package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/movie/internal/core"
	"github.com/karaMuha/go-movie/movie/internal/gateway/metadataGateway"
	"github.com/karaMuha/go-movie/movie/internal/gateway/ratingGateway"
	"github.com/karaMuha/go-movie/movie/internal/rest/v1"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
)

const serviceName = "movie"

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

	metadataGateway := metadataGateway.NewMetadataRestGateway(registry)
	ratingGateway := ratingGateway.NewRatginRestGateway(registry)
	app := core.New(&metadataGateway, &ratingGateway)
	movieHandlerV1 := rest.NewMovieHandlerV1(&app)

	mux := http.NewServeMux()
	setupEndpoints(mux, movieHandlerV1)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupEndpoints(mux *http.ServeMux, movieHandlerV1 rest.MovieHandlerV1) {
	movieV1 := http.NewServeMux()
	movieV1.HandleFunc("GET /get-movie-details", movieHandlerV1.HandleGetMovieDetails)

	mux.Handle("/v1/", http.StripPrefix("/v1", movieV1))
}
