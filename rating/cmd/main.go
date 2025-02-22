package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/database/postgres"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
	"github.com/karaMuha/go-movie/rating/config"
	"github.com/karaMuha/go-movie/rating/internal/core"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/rating/internal/cronjob"
	grpchandler "github.com/karaMuha/go-movie/rating/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/rating/internal/queue/consumer"
	postgres_repo "github.com/karaMuha/go-movie/rating/internal/repository/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

const serviceName = "rating"

func main() {
	log.Println("Starting movie rating service")
	config := config.NewConfig()

	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("%s%s", "ratings-service", config.PortExposed)
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	db, err := postgres.ConnectToDb(config.DbDriver, config.DbConnection)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ratingRepo := postgres_repo.NewRatingRepository(db)
	aggregatedRatingRepo := postgres_repo.NewAggregatedMetadataRepository(db)
	app := core.New(&ratingRepo, &aggregatedRatingRepo)

	ratingConsumer := consumer.NewRatingEventConsumer(&app, config.KafkaAddress, "ratings", "rating")
	defer ratingConsumer.RatingReader.Close()
	go ratingConsumer.StartReadingRatingEvents()

	metadataEventRepository := postgres_repo.NewMetadataEventRepository(db)
	metadataConsumer := consumer.NewMetadataEventConsumer(&app, config.KafkaAddress, "metadata", "metadata", &metadataEventRepository)
	defer metadataConsumer.Reader.Close()
	go metadataConsumer.StartReadingMetadataEvents()

	cronjob := cronjob.NewCronjob(&metadataEventRepository, app)
	go cronjob.RunMetadata()

	startGrpc(&app, config.PortExposed)

}

func startGrpc(app driving.IApplication, port string) {
	ratingHandler := grpchandler.NewRatingHandler(app)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRatingServiceServer(server, &ratingHandler)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}
