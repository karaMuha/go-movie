package grpcgateway

import (
	"context"
	"net/http"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/pkg/grpcutil"
)

type MetadataGateway struct {
	registry discovery.Registry
}

var _ driven.IMetadataGateway = (*MetadataGateway)(nil)

func NewMetadataGateway(registry discovery.Registry) MetadataGateway {
	return MetadataGateway{registry: registry}
}

func (g *MetadataGateway) GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, *dtos.RespErr) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata-service", g.registry)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	defer conn.Close()

	client := pb.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &pb.GetMetadataRequest{MovieId: movieID})
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	if resp.ResponseStatus.StatusCode > 299 {
		return nil, &dtos.RespErr{
			StatusCode:    int(resp.ResponseStatus.StatusCode),
			StatusMessage: resp.ResponseStatus.Message,
		}
	}

	return metadataModel.MetadataFromProto(resp.Metadata), nil
}

func (g *MetadataGateway) SubmitMetadata(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata-service", g.registry)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	defer conn.Close()

	client := pb.NewMetadataServiceClient(conn)
	params := pb.SubmitMetadataRequest{
		Metadata: metadataModel.MetadataToProto(metadata),
	}
	resp, err := client.SubmitMetadata(ctx, &params)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	if resp.ResponseStatus.StatusCode > 299 {
		return nil, &dtos.RespErr{
			StatusCode:    int(resp.ResponseStatus.StatusCode),
			StatusMessage: resp.ResponseStatus.Message,
		}
	}

	return metadataModel.MetadataFromProto(resp.Metadata), nil
}
