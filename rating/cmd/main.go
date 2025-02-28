package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
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

const serviceName = "ratings-service"

func main() {
	log.Println("Starting rating service")
	config := config.NewConfig()

	// Service discovery
	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("%s%s", serviceName, config.PortExposed)
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		log.Fatalln(err)
	}
	defer registry.Deregister(ctx, instanceID, serviceName)
	go reportHealthyState(registry, instanceID)

	// Database
	db, err := postgres.ConnectToDb(config.DbDriver, config.DbConnection)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Application
	ratingRepo := postgres_repo.NewRatingRepository(db)
	aggregatedRatingRepo := postgres_repo.NewAggregatedRatingRepository(db)
	metadataEventRepository := postgres_repo.NewMetadataEventRepository(db)
	app := core.New(&ratingRepo, &aggregatedRatingRepo)

	// Message consumers
	ratingConsumer := consumer.NewRatingEventConsumer(&app, config.KafkaAddress, "ratings", "rating")
	defer ratingConsumer.RatingReader.Close()
	go ratingConsumer.StartReadingRatingEvents()

	metadataConsumer := consumer.NewMetadataEventConsumer(&app, config.KafkaAddress, "metadata", "metadata", &metadataEventRepository)
	defer metadataConsumer.StopReadingEvents()
	defer metadataConsumer.Reader.Close()
	go metadataConsumer.StartReadingMetadataEvents()

	// Cronjob
	cronjob := cronjob.NewCronjob(&metadataEventRepository, &app)
	defer cronjob.GracefulStop()
	go cronjob.RunMetadata()

	// Server
	server := initGrpcServer(&app)
	go startServer(server, config.PortExposed)

	// Graceful stop
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()
	log.Println("Shutting down server gracefully")
	server.GracefulStop()
}

func initGrpcServer(app driving.IApplication) *grpc.Server {
	ratingHandler := grpchandler.NewRatingHandler(app)
	server := grpc.NewServer()
	pb.RegisterRatingServiceServer(server, &ratingHandler)
	reflection.Register(server)
	return server
}

func startServer(srv *grpc.Server, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Start listening")
	if err := srv.Serve(listener); err != nil {
		log.Printf("Stopped listening: %v\n", err)
	}
}

func reportHealthyState(registry *consul.ConsulRegistry, instanceID string) {
	for {
		if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
			log.Println("Failed to report healthy state: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
