package ratingModel

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across al types.
type RecordType string

// Existing record types.
const RecordTypeMovie = RecordType("movie")

type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for
// some record.
type Rating struct {
	RecordID   string `json:"record_id"`
	RecordType string `json:"record_type"`
	UserID     string `json:"user_id"`
	Value      int    `json:"value"`
}
