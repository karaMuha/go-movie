package domain

import "errors"

func CreateMetadata(title, director string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	if director == "" {
		return errors.New("director cannot be empty")
	}

	return nil
}
