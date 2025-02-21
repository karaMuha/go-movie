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
	grpchandler "github.com/karaMuha/go-movie/metadata/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/metadata/internal/queue/producer"
	"github.com/karaMuha/go-movie/metadata/internal/repository/memory"
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
	//metadataPostgresRepo := postgres_repo.NewMetadataRepository(db)

	metadataRepo := memory.New()
	producer := producer.NewMessageProducer(config.KafkaAddress, "metadata")
	defer producer.Writer.Close()
	app := core.New(metadataRepo, producer)

	//startRest(&app, port)
	startGrpc(&app, config.Port)
}

func startGrpc(app driving.IApplication, port string) {
	metdataDataHandlerGrpc := grpchandler.NewMetadataHandler(app)
	// address := fmt.Sprintf("%s%s", domain, port)

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

/* func setupRestEndpoints(mux *http.ServeMux, metadataHandlerV1 rest.MetadataHandlerV1) {
	metadataV1 := http.NewServeMux()
	metadataV1.HandleFunc("GET /get-metadata", metadataHandlerV1.GetMetadata)

	mux.Handle("/v1/", http.StripPrefix("/v1", metadataV1))
}

func startRest(app driving.IApplication, port string) {
	metadataHandlerV1 := rest.NewMetadataHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, metadataHandlerV1)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
} */
