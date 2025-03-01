package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
)

type MetadataHandlerV1 struct {
	app driving.IApplication
}

func NewMetadataHandlerV1(app driving.IApplication) MetadataHandlerV1 {
	return MetadataHandlerV1{
		app: app,
	}
}

func (h MetadataHandlerV1) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "no id specified", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	metadata, respErr := h.app.GetMetadata(ctx, id)

	if respErr != nil {
		http.Error(w, respErr.StatusMessage, respErr.StatusCode)
		return
	}

	err := json.NewEncoder(w).Encode(metadata)
	if err != nil {
		log.Printf("GetMetadata error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
