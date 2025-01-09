package metadataModel

import "github.com/karaMuha/go-movie/pb"

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
