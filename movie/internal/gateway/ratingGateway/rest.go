package ratingGateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	ratingModel "github.com/karaMuha/go-movie/rating/ratingModel"
)

type RatingRestGateway struct {
	address string
}

var _ driven.IRatingGateway = (*RatingRestGateway)(nil)

func NewRatginRestGateway(address string) RatingRestGateway {
	return RatingRestGateway{
		address: address,
	}
}

func (g *RatingRestGateway) GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error) {
	url := fmt.Sprintf("%s/v1/get-rating", g.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("record_id", string(recordID))
	values.Add("record_type", string(recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, domain.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (g *RatingRestGateway) SubmitRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error {
	url := fmt.Sprintf("%s/v1/submit-rating", g.address)
	body := ratingModel.Rating{
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
