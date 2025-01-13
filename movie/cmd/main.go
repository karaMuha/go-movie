package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/movie/internal/core"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/movie/internal/endpoint/rest/v1"
	grpcgateway "github.com/karaMuha/go-movie/movie/internal/gateway/grpc"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
)

const serviceName = "movie"
const domain = "localhost"

func main() {
	log.Println("Starting movie service")
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	registry, err := consul.NewConsulRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port))
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
	app := core.New(&metadataGateway, &ratingGateway)
	startRest(&app, port)

}

func setupRestEndpoints(mux *http.ServeMux, movieHandlerV1 rest.MovieHandlerV1) {
	movieV1 := http.NewServeMux()
	movieV1.HandleFunc("GET /get-movie-details", movieHandlerV1.HandleGetMovieDetails)

	mux.Handle("/v1/", http.StripPrefix("/v1", movieV1))
}

func startRest(app driving.IApplication, port int) {
	movieHandlerV1 := rest.NewMovieHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, movieHandlerV1)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
