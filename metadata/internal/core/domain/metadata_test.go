package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMetadata(t *testing.T) {
	tests := []struct {
		testName   string
		title      string
		director   string
		recordType string
		err        error
	}{
		{"TestNoTitle", "", "Christopher Nolan", "movie", errors.New("title cannot be empty")},
		{"TestNoDirector", "Batman", "", "movie", errors.New("director cannot be empty")},
		{"TestSuccess", "Batman", "Christopher Nolan", "movie", nil},
	}

	for _, test := range tests {
		err := CreateMetadata(test.title, test.director, test.recordType)
		require.Equal(t, test.err, err)
	}
}
