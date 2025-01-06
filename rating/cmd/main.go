package main

import (
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/rating/internal/core"
	"github.com/karaMuha/go-movie/rating/internal/repository/memory"
	"github.com/karaMuha/go-movie/rating/internal/rest/v1"
)

func main() {
	log.Println("Starting movie rating service")
	ratingRepo := memory.New()
	app := core.New(&ratingRepo)
	ratingHandlerV1 := rest.NewRatingHandlerV1(&app)

	mux := http.NewServeMux()
	setupEndpoints(mux, ratingHandlerV1)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupEndpoints(mux *http.ServeMux, ratingHandlerV1 rest.RatingHandlerV1) {
	ratingV1 := http.NewServeMux()
	ratingV1.HandleFunc("GET /get-rating", ratingHandlerV1.HandleGetRating)
	ratingV1.HandleFunc("POST /submit-rating", ratingHandlerV1.HandleSubmitRating)

	mux.Handle("/v1/", http.StripPrefix("/v1", ratingV1))
}
