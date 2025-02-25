package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/pkg/http/response"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type MovieHandlerV1 struct {
	app driving.IApplication
}

func NewMovieHandlerV1(app driving.IApplication) MovieHandlerV1 {
	return MovieHandlerV1{
		app: app,
	}
}

func (h *MovieHandlerV1) HandleGetMovieDetails(w http.ResponseWriter, r *http.Request) {
	movieID := r.URL.Query().Get("id")
	movieDetails, respErr := h.app.GetMovieDetails(r.Context(), movieID)

	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movieDetails)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MovieHandlerV1) HandleSubmitRating(w http.ResponseWriter, r *http.Request) {
	var rating ratingmodel.Rating
	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	respErr := h.app.SubmitRating(ctx, &rating)
	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}
}

func (h *MovieHandlerV1) HandleGetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	metadata, respErr := h.app.GetMetadata(r.Context(), id)

	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}

	response.WithJson(w, http.StatusOK, metadata)
}

func (h *MovieHandlerV1) HandleSubmitMetadata(w http.ResponseWriter, r *http.Request) {
	var metadata metadataModel.Metadata
	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, respErr := h.app.SubmitMetadata(r.Context(), &metadata)
	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}

	response.WithJson(w, http.StatusCreated, resp)
}
