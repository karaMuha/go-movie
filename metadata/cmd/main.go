package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/karaMuha/go-movie/metadata/config"
	"github.com/karaMuha/go-movie/metadata/internal/core"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/metadata/internal/cronjob"
	grpchandler "github.com/karaMuha/go-movie/metadata/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/metadata/internal/queue/producer"
	postgres_repo "github.com/karaMuha/go-movie/metadata/internal/repository/postgres"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/database/postgres"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

const serviceName = "metadata"

func main() {
	log.Println("Starting metadata service")
	config := config.NewConfig()

	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("%s%s", "metadata-service", config.PortExposed)
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
	fmt.Println("Connected to database")

	metadataRepo := postgres_repo.NewMetadataRepository(db)
	metadataEventRepo := postgres_repo.NewMetadataEventRepository(db)
	producer := producer.NewMessageProducer(config.KafkaAddress, "metadata")
	defer producer.Writer.Close()
	app := core.New(&metadataRepo, producer, &metadataEventRepo)
	cronjob := cronjob.NewCronjob(&metadataEventRepo, producer)

	// if events failed to be published they are saved in the database
	// cronjob loops through the table each minute and tries to publish them
	go cronjob.Run()

	startGrpc(&app, config.Port)
}

func startGrpc(app driving.IApplication, port string) {
	metdataDataHandlerGrpc := grpchandler.NewMetadataHandler(app)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetadataServiceServer(server, &metdataDataHandlerGrpc)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Cannot start grpc server: %v\n", err)
	}
}
