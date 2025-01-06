package main

import (
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/metadata/internal/core"
	"github.com/karaMuha/go-movie/metadata/internal/repository/memory"
	"github.com/karaMuha/go-movie/metadata/internal/rest/v1"
)

func main() {
	log.Println("Starting movie metadata service")
	metadataRepo := memory.New()
	app := core.New(metadataRepo)
	metadataHandlerV1 := rest.NewMetadataHandlerV1(&app)

	mux := http.NewServeMux()
	setupEndpoints(mux, metadataHandlerV1)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupEndpoints(mux *http.ServeMux, metadataHandlerV1 rest.MetadataHandlerV1) {
	metadataV1 := http.NewServeMux()
	metadataV1.HandleFunc("GET /get-metadata", metadataHandlerV1.GetMetadata)

	mux.Handle("/v1/", http.StripPrefix("/v1", metadataV1))
}
