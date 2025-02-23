package metadataModel

import (
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

// MetadataToProto converts a Metadata struct into a
// generated proto counterpart
func MetadataToProto(metadata *Metadata) *pb.Metadata {
	return &pb.Metadata{
		Id:          metadata.ID,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}

func MetadataFromProto(metadata *pb.Metadata) *Metadata {
	return &Metadata{
		ID:          metadata.Id,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}

func RespErrToProto(respErr *dtos.RespErr) *pb.ResponseStatus {
	return &pb.ResponseStatus{
		StatusCode: int32(respErr.StatusCode),
		Message:    respErr.StatusMessage,
	}
}
