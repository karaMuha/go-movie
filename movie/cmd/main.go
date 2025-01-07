package main

import (
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core"
	"github.com/karaMuha/go-movie/movie/internal/gateway/metadataGateway"
	"github.com/karaMuha/go-movie/movie/internal/gateway/ratingGateway"
	"github.com/karaMuha/go-movie/movie/internal/rest/v1"
)

func main() {
	log.Println("Starting movie service")
	metadataGateway := metadataGateway.NewMetadataRestGateway("http://localhost:8080")
	ratingGateway := ratingGateway.NewRatginRestGateway("http://localhost:8081")
	app := core.New(&metadataGateway, &ratingGateway)
	movieHandlerV1 := rest.NewMovieHandlerV1(&app)

	mux := http.NewServeMux()
	setupEndpoints(mux, movieHandlerV1)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupEndpoints(mux *http.ServeMux, movieHandlerV1 rest.MovieHandlerV1) {
	movieV1 := http.NewServeMux()
	movieV1.HandleFunc("GET /get-movie-details", movieHandlerV1.HandleGetMovieDetails)

	mux.Handle("/v1/", http.StripPrefix("/v1", movieV1))
}
