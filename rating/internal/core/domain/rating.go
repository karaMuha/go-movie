package domain

import (
	"errors"

	"github.com/google/uuid"
)

func SubmitRating(recordID, recordType, userID string, value int) error {
	err := uuid.Validate(recordID)
	if err != nil {
		return err
	}

	if recordType != "movie" {
		return errors.New("unsupported record type")
	}

	// register and login not implemented yet
	/* if userID == "" {
		return errors.New("user ID is empty")
	} */

	if value > 10 || value < 0 {
		return errors.New("value must be between 0 and 10")
	}

	return nil
}
