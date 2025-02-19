package restgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/discovery"
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

func (g *RatingGateway) GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, error) {
	addresses, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, 0, err
	}
	url := fmt.Sprintf("http://%s/v1/get-rating", addresses[rand.Intn(len(addresses))])
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("record_id", string(recordID))
	values.Add("record_type", string(recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, 0, domain.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, 0, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, 0, err
	}
	return v, 0, nil
}

func (g *RatingGateway) SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error {
	addresses, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v1/submit-rating", addresses[rand.Intn(len(addresses))])
	body := ratingmodel.Rating{
		RecordID:   string(recordID),
		RecordType: string(recordType),
		UserID:     rating.UserID,
		Value:      rating.Value,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}
	return nil
}
