package repositories_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestSaveGroupParams_MarshalJSON(t *testing.T) {
	tcs := []struct {
		name           string
		input          r.SaveGroupParams
		expectedOutput string
	}{
		{
			name:           "Marshals to JSON",
			input:          r.SaveGroupParams{Id: "C1234f49365c6b492b337189e3343a9d9"},
			expectedOutput: `{"id":"C1234f49365c6b492b337189e3343a9d9"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			marshaledData, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			jsonString := string(marshaledData)

			if jsonString != tc.expectedOutput {
				t.Errorf("Expected JSON %s, but got %s", tc.expectedOutput, jsonString)
			}
		})
	}
}

func TestD1GroupRepository_Save(t *testing.T) {
	testCases := []struct {
		name          string
		statusCode    int
		expectError   bool
		expectedError error
	}{
		{"Returns 204 (Success)", http.StatusNoContent, false, nil},
		{"Returns 304 (Success)", http.StatusNotModified, true, r.ErrorAlreadyExists},
		{"Returns 400 (Failure)", http.StatusBadRequest, true, fmt.Errorf("unexpected response status: %d", http.StatusBadRequest)},
		{"Returns 500 (Failure)", http.StatusInternalServerError, true, fmt.Errorf("unexpected response status: %d", http.StatusInternalServerError)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := u.NewMockHttpClient(tc.statusCode)
			repo := r.NewD1GroupRepository("https://example.com/api/v1/groups", mockClient)

			group := &g.Group{Id: "C1234f49365c6b492b337189e3343a9d9"}
			result, err := repo.Save(group)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none for status code %d", tc.statusCode)
				assert.Nil(t, result, "Expected result to be nil on error, but got: %v", result)
				if tc.statusCode == http.StatusNotModified {
					assert.Equal(t, tc.expectedError, err, "Expected error to be ErrorAlreadyExists but got: %v", err)
				}
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
				assert.NotNil(t, result, "Expected result to be non-nil but got nil")
				assert.Equal(t, group.Id, result.Id, "Returned group ID, %s, does not match input group ID %s", result.Id, group.Id)
			}
		})
	}
}

func TestD1GroupRepository_Destroy(t *testing.T) {
	tcs := []struct {
		name          string
		statusCode    int
		shouldFail    bool
		expectedError error
	}{
		{
			name:          "Returns nil on successful delete",
			statusCode:    http.StatusNoContent,
			shouldFail:    false,
			expectedError: nil,
		},
		{
			name:          "Returns ErrorNotFound when group does not exist",
			statusCode:    http.StatusNotFound,
			shouldFail:    false,
			expectedError: r.ErrorNotFound,
		},
		{
			name:          "Returns an error on unexpected status code",
			statusCode:    http.StatusInternalServerError,
			shouldFail:    false,
			expectedError: fmt.Errorf("unexpected response status: %d", http.StatusInternalServerError),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := u.NewMockHttpClient(tc.statusCode)
			repo := r.NewD1GroupRepository("https://example.com/api/v1/groups", mockClient)

			group := &g.Group{Id: "C1234f49365c6b492b337189e3343a9d9"}

			err := repo.Destroy(group)

			if tc.expectedError == nil {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
			} else {
				assert.Error(t, err, "Expected an error but got none for status code %d", tc.statusCode)
			}
		})
	}
}
