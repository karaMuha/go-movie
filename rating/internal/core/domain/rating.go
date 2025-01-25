package domain

import "errors"

func CreateRating(recordID, recordType, userID string, value int) error {
	if recordID == "" {
		return errors.New("record ID is empty")
	}
	if recordType != "movie" {
		return errors.New("unsupported record type")
	}
	if userID == "" {
		return errors.New("user ID is empty")
	}
	if value > 10 || value < 0 {
		return errors.New("value must be between 0 and 10")
	}
	return nil
}
