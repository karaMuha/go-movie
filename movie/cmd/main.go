package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/karaMuha/go-movie/movie/config"
	"github.com/karaMuha/go-movie/movie/internal/core"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/movie/internal/endpoint/rest/v1"
	grpcgateway "github.com/karaMuha/go-movie/movie/internal/gateway/grpc"
	"github.com/karaMuha/go-movie/movie/internal/queue/producer"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"

	_ "github.com/lib/pq"
)

const serviceName = "movies-service"

func main() {
	log.Println("Starting movie service")
	config := config.NewConfig()

	// Service discovery
	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("%s%s", serviceName, config.PortExposed)
	err = registry.Register(ctx, instanceID, serviceName, hostPort)
	if err != nil {
		log.Fatalln(err)
	}
	go reportHealtyhState(registry, instanceID)
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Message producer
	producer := producer.NewMessageProducer(config.KafkaAddress, "ratings")
	defer producer.Writer.Close()

	// Application
	metadataGateway := grpcgateway.NewMetadataGateway(registry)
	ratingGateway := grpcgateway.NewRatingGateway(registry)
	app := core.New(&metadataGateway, &ratingGateway, producer)

	// Server
	server := initRestServer(&app, config.Port)
	go startServer(server)

	// Graceful shutdown
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	log.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown with error: %v", err)
	}
	log.Println("Shutdown complete")
}

func setupRestEndpoints(mux *http.ServeMux, movieHandlerV1 rest.MovieHandlerV1) {
	movieV1 := http.NewServeMux()
	movieV1.HandleFunc("GET /get-movie-details", movieHandlerV1.HandleGetMovieDetails)
	movieV1.HandleFunc("POST /submit-rating", movieHandlerV1.HandleSubmitRating)
	movieV1.HandleFunc("GET /get-metadata", movieHandlerV1.HandleGetMetadata)
	movieV1.HandleFunc("POST /submit-metadata", movieHandlerV1.HandleSubmitMetadata)

	mux.Handle("/v1/", http.StripPrefix("/v1", movieV1))
}

func initRestServer(app driving.IApplication, port string) *http.Server {
	movieHandlerV1 := rest.NewMovieHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, movieHandlerV1)

	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server
}

func startServer(src *http.Server) {
	if err := src.ListenAndServe(); err != nil {
		log.Printf("Stopped listening: %v\n", err)
	}
}

func reportHealtyhState(registry *consul.ConsulRegistry, instanceID string) {
	for {
		err := registry.ReportHealthyState(instanceID, serviceName)
		if err != nil {
			log.Printf("Failed to report healty state: %s\n", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
