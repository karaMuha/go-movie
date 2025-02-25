package restgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/discovery"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type MetadataGateway struct {
	registry discovery.Registry
}

var _ driven.IMetadataGateway = (*MetadataGateway)(nil)

func NewMetadataGateway(registry discovery.Registry) MetadataGateway {
	return MetadataGateway{
		registry: registry,
	}
}

func (g *MetadataGateway) GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, *dtos.RespErr) {
	addresses, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	url := fmt.Sprintf("http://%s/v1/get-metadata", addresses[rand.Intn(len(addresses))]) // #nosec G404
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", movieID)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return nil, &dtos.RespErr{
			StatusCode:    resp.StatusCode,
			StatusMessage: resp.Status,
		}
	}

	var v *metadataModel.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return v, nil
}

func (g *MetadataGateway) SubmitMetadata(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	return nil, nil
}
