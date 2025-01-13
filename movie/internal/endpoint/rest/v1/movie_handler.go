package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
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
