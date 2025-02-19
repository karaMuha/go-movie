package ratingmodel

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across al types.
type RecordType string

// Existing record types.
const RecordTypeMovie RecordType = "movie"

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

type RatingEventType string

const RatingEventTypeSubmit RatingEventType = "submit"
const RatingEventTypeDelete RatingEventType = "delete"

type RatingEvent struct {
	RecordID   string          `json:"record_id"`
	RecordType string          `json:"record_type"`
	UserID     string          `json:"user_id"`
	Value      int             `json:"value"`
	EventType  RatingEventType `json:"event_type"`
}

type AggregatedRating struct {
	ID         string  `json:"id"`
	RecordType string  `json:"record_type"`
	Rating     float64 `json:"rating"`
}
