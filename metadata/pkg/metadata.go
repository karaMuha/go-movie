package metadataModel

// Metadata defines the movie metadata
type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}

type MetadataRecordType string

const RecordTypeMovie MetadataRecordType = "movie"

type MetadataEventType string

const MetadataEventTypeSubmitted MetadataEventType = "submit"

type MetadataEvent struct {
	ID         string             `json:"id"`
	RecordType MetadataRecordType `json:"record_type"`
	EventType  MetadataEventType  `json:"event_type"`
}
