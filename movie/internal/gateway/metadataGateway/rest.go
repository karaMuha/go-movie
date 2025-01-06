package metadataGateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	metadataModel "github.com/karaMuha/go-movie/metadata/metadataModel"
	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
)

type MetadataRestGateway struct {
	address string
}

var _ driven.IMetadataGateway = (*MetadataRestGateway)(nil)

func NewMetadataRestGateway(address string) MetadataRestGateway {
	return MetadataRestGateway{
		address: address,
	}
}

func (g *MetadataRestGateway) GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, error) {
	url := fmt.Sprintf("%s/v1/get-metadata", g.address)
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
