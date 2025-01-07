package metadataGateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	metadataModel "github.com/karaMuha/go-movie/metadata/metadataModel"
	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/discovery"
)

type MetadataRestGateway struct {
	registry discovery.Registry
}

var _ driven.IMetadataGateway = (*MetadataRestGateway)(nil)

func NewMetadataRestGateway(registry discovery.Registry) MetadataRestGateway {
	return MetadataRestGateway{
		registry: registry,
	}
}

func (g *MetadataRestGateway) GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, error) {
	addresses, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/v1/get-metadata", addresses[rand.Intn(len(addresses))])
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", movieID)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var v *metadataModel.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
