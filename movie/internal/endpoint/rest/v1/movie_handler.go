package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
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
	movieDetails, err := h.app.GetMovieDetails(r.Context(), movieID)

	if errors.Is(err, domain.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("GetMovieDetails error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(movieDetails)
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

	err = h.app.SubmitRating(ctx, &rating)
	if err != nil {
		log.Printf("SaveRating error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
