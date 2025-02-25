package domain

import (
	"errors"
	"strings"
)

func CreateMetadata(title, director, recordType string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(director) == "" {
		return errors.New("director cannot be empty")
	}

	if strings.TrimSpace(recordType) == "" {
		return errors.New("record type cannot be empty")
	}

	return nil
}
