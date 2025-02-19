package movieModel

import metadataModel "github.com/karaMuha/go-movie/metadata/pkg"

// MovieDetails includes movie metadata and its
// aggregated rating.
type MovieDetails struct {
	Rating        float64                `json:"rating"`
	AmountRatings int                    `json:"amount_ratings"`
	Metadata      metadataModel.Metadata `json:"metadata"`
}
