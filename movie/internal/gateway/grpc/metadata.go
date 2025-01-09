package grpcgateway

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	"github.com/karaMuha/go-movie/pkg/grpcutil"
)

type MetadataGateway struct {
	registry discovery.Registry
}

var _ driven.IMetadataGateway = (*MetadataGateway)(nil)

func NewMetadataGateway(registry discovery.Registry) MetadataGateway {
	return MetadataGateway{registry: registry}
}

func (g *MetadataGateway) GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &pb.GetMetadataRequest{MovieId: movieID})
	if err != nil {
		return nil, err
	}

	return metadataModel.MetadataFromProto(resp.Metadata), nil
}
