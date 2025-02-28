package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
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

const serviceName = "metadata-service"

func main() {
	log.Println("Starting metadata service")
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

	// Message producer
	producer := producer.NewMessageProducer(config.KafkaAddress, "metadata")
	defer producer.Writer.Close()

	// Application
	metadataRepo := postgres_repo.NewMetadataRepository(db)
	metadataEventRepo := postgres_repo.NewMetadataEventRepository(db)
	app := core.New(&metadataRepo, producer, &metadataEventRepo)

	// Cronjob
	cronjob := cronjob.NewCronjob(&metadataEventRepo, producer)
	defer cronjob.GracefulStop()
	go cronjob.Run()

	// Server
	server := initGrpcServer(&app)
	go runServer(server, config.PortExposed)

	// Graceful stop
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()
	log.Println("Shutting down server gracefully")
	server.GracefulStop()
}

func initGrpcServer(app driving.IApplication) *grpc.Server {
	metdataDataHandlerGrpc := grpchandler.NewMetadataHandler(app)

	server := grpc.NewServer()
	pb.RegisterMetadataServiceServer(server, &metdataDataHandlerGrpc)
	reflection.Register(server)
	return server
}

func runServer(srv *grpc.Server, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Cannot start grpc server: %v\n", err)
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
