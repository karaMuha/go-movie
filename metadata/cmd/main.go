package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/metadata/config"
	"github.com/karaMuha/go-movie/metadata/internal/core"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	grpchandler "github.com/karaMuha/go-movie/metadata/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/metadata/internal/endpoint/rest/v1"
	"github.com/karaMuha/go-movie/metadata/internal/repository/memory"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

const serviceName = "metadata"

func main() {
	log.Println("Starting movie metadata service")
	config := config.NewConfig()

	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost%s", config.Port)); err != nil {
		panic(err)
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

	/* connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName, config.DbSslMode)
	db, err := postgres.ConnectToDb(config.DbDriver, connectionString)
	if err != nil {
		panic(err)
	}
	metadataPostgresRepo := postgres_repo.NewMetadataRepository(db) */

	metadataRepo := memory.New()
	app := core.New(metadataRepo)

	//startRest(&app, port)
	startGrpc(&app, config.Domain, config.Port)
}

func setupRestEndpoints(mux *http.ServeMux, metadataHandlerV1 rest.MetadataHandlerV1) {
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
}

func startGrpc(app driving.IApplication, domain string, port string) {
	metdataDataHandlerGrpc := grpchandler.NewMetadataHandler(app)
	address := fmt.Sprintf("%s%s", domain, port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetadataServiceServer(server, &metdataDataHandlerGrpc)
	server.Serve(listener)
}
