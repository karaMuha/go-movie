package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/metadata/internal/core"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	grpchandler "github.com/karaMuha/go-movie/metadata/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/metadata/internal/endpoint/rest/v1"
	"github.com/karaMuha/go-movie/metadata/internal/repository/memory"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
	"google.golang.org/grpc"
)

const serviceName = "metadata"
const domain = "localhost"

func main() {
	log.Println("Starting movie metadata service")
	var port int
	flag.IntVar(&port, "port", 8080, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)

	registry, err := consul.NewConsulRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
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

	metadataRepo := memory.New()
	app := core.New(metadataRepo)

	//startRest(&app, port)
	startGrpc(&app, port)
}

func setupRestEndpoints(mux *http.ServeMux, metadataHandlerV1 rest.MetadataHandlerV1) {
	metadataV1 := http.NewServeMux()
	metadataV1.HandleFunc("GET /get-metadata", metadataHandlerV1.GetMetadata)

	mux.Handle("/v1/", http.StripPrefix("/v1", metadataV1))
}

func startRest(app driving.IApplication, port int) {
	metadataHandlerV1 := rest.NewMetadataHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, metadataHandlerV1)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func startGrpc(app driving.IApplication, port int) {
	metdataDataHandlerGrpc := grpchandler.NewMetadataHandler(app)
	address := fmt.Sprintf("%s:%d", domain, port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetadataServiceServer(server, &metdataDataHandlerGrpc)
	server.Serve(listener)
}
