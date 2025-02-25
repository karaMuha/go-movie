package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestSubmitRating(t *testing.T) {
	tests := []struct {
		testName   string
		recordID   string
		recordType string
		userID     string
		value      int
		wantErr    bool
	}{
		{"TestNoRecordID", "", "movie", "", 5, true},
		{"TestNoRecordType", uuid.NewString(), "", "", 5, true},
		{"TestWrongRecordType", uuid.NewString(), "music", "", 5, true},
		{"TestValueToHigh", uuid.NewString(), "movie", "", 11, true},
		{"TestValueNoLow", uuid.NewString(), "movie", "", -1, true},
		{"TestSuccessfulSubmit", uuid.NewString(), "movie", "", 5, false},
	}

	for _, test := range tests {
		err := SubmitRating(test.recordID, test.recordType, test.userID, test.value)
		if err == nil && test.wantErr {
			t.Errorf("Submit rating test error: want error but got none for test case: %s", test.testName)
		}
		if err != nil && !test.wantErr {
			t.Errorf("Submit rating test error: want no error but got error for test case: %s", test.testName)
		}
	}
}
