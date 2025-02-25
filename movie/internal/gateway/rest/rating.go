package restgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/discovery"
	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type RatingGateway struct {
	registry discovery.Registry
}

var _ driven.IRatingGateway = (*RatingGateway)(nil)

func NewRatginGateway(registry discovery.Registry) RatingGateway {
	return RatingGateway{
		registry: registry,
	}
}

func (g *RatingGateway) GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, *dtos.RespErr) {
	addresses, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	url := fmt.Sprintf("http://%s/v1/get-rating", addresses[rand.Intn(len(addresses))]) // #nosec G404
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("record_id", string(recordID))
	values.Add("record_type", string(recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return 0.0, 0, &dtos.RespErr{
			StatusCode:    resp.StatusCode,
			StatusMessage: resp.Status,
		}
	}
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return v, 0, nil
}

func (g *RatingGateway) SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr {
	addresses, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	url := fmt.Sprintf("%s/v1/submit-rating", addresses[rand.Intn(len(addresses))]) // #nosec G404
	body := ratingmodel.Rating{
		RecordID:   string(recordID),
		RecordType: string(recordType),
		UserID:     rating.UserID,
		Value:      rating.Value,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return &dtos.RespErr{
			StatusCode:    resp.StatusCode,
			StatusMessage: resp.Status,
		}
	}
	return nil
}
