package movieModel

import metadataModel "github.com/karaMuha/go-movie/metadata/metadataModel"

// MovieDetails includes movie metadata and its
// aggregated rating.
type MovieDetails struct {
	Rating   float64
	Metadata metadataModel.Metadata
}
