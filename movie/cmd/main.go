package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

const serviceName = "movie"

func main() {
	log.Println("Starting movie service")
	config := config.NewConfig()

	registry, err := consul.NewConsulRegistry(config.ConsulAddress)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s%s", config.Domain, config.Port))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := registry.ReportHealthyState(instanceID, serviceName)
			if err != nil {
				log.Printf("Failed to report healty state: %s\n", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	//metadataGateway := restgateway.NewMetadataGateway(registry)
	//ratingGateway := restgateway.NewRatginGateway(registry)
	metadataGateway := grpcgateway.NewMetadataGateway(registry)
	ratingGateway := grpcgateway.NewRatingGateway(registry)
	producer := producer.NewMessageProducer(config.KafkaAddress, "ratings")
	defer producer.Writer.Close()

	app := core.New(&metadataGateway, &ratingGateway, producer)

	startRest(&app, config.Port)

}

func setupRestEndpoints(mux *http.ServeMux, movieHandlerV1 rest.MovieHandlerV1) {
	movieV1 := http.NewServeMux()
	movieV1.HandleFunc("GET /get-movie-details", movieHandlerV1.HandleGetMovieDetails)
	movieV1.HandleFunc("POST /submit-rating", movieHandlerV1.HandleSubmitRating)

	mux.Handle("/v1/", http.StripPrefix("/v1", movieV1))
}

func startRest(app driving.IApplication, port string) {
	movieHandlerV1 := rest.NewMovieHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, movieHandlerV1)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
