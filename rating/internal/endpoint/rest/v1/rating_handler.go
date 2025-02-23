package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	model "github.com/karaMuha/go-movie/rating/pkg"
)

type RatingHandlerV1 struct {
	app driving.IApplication
}

func NewRatingHandlerV1(app driving.IApplication) RatingHandlerV1 {
	return RatingHandlerV1{
		app: app,
	}
}

func (h *RatingHandlerV1) HandleGetRating(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.URL.Query().Get("record_id"))
	if recordID == "" {
		http.Error(w, "no record ID specified", http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.URL.Query().Get("record_type"))
	if recordType == "" {
		http.Error(w, "no record type specified", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	aggregatedRating, _, respErr := h.app.GetAggregatedRating(ctx, recordID, recordType)
	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}

	err := json.NewEncoder(w).Encode(aggregatedRating)
	if err != nil {
		log.Printf("GetAggregatedRating error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RatingHandlerV1) HandleSubmitRating(w http.ResponseWriter, r *http.Request) {
	var rating model.Rating
	err := json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	respErr := h.app.SubmitRating(ctx, model.RecordID(rating.RecordID), model.RecordType(rating.RecordType), &rating)
	if respErr != nil {
		log.Printf("SaveRating error: %v\n", err)
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}
}
