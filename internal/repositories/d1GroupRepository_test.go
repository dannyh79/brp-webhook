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
		name               string
		expectedReqHeaders map[string]string
		expectedRespStatus int
		expectError        bool
		expectedError      error
	}{
		{
			name:               "Returns 204 (Success)",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusNoContent,
			expectError:        false,
		},
		{
			name:               "Returns 304 (Success)",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusNotModified,
			expectError:        true,
			expectedError:      r.ErrorAlreadyExists,
		},
		{
			name:               "Returns 400 (Failure)",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusBadRequest,
			expectError:        true,
			expectedError:      fmt.Errorf("unexpected response status: %d", http.StatusBadRequest),
		},
		{
			name:               "Returns 500 (Failure)",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusInternalServerError,
			expectError:        true,
			expectedError:      fmt.Errorf("unexpected response status: %d", http.StatusInternalServerError),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var hs http.Header
			mockClient := u.NewMockHttpClient(tc.expectedRespStatus, func(req *http.Request) {
				hs = req.Header
			})
			repo := r.NewD1GroupRepository("https://example.com/api/v1/groups", mockClient)

			group := &g.Group{Id: "C1234f49365c6b492b337189e3343a9d9"}
			result, err := repo.Save(group)

			for k, v := range tc.expectedReqHeaders {
				assert.Equal(t, v, hs.Get(k), "Expected header %s to be %s, but got %s", k, v, hs.Get(k))
			}

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none for status code %d", tc.expectedRespStatus)
				assert.Nil(t, result, "Expected result to be nil on error, but got: %v", result)
				if tc.expectedRespStatus == http.StatusNotModified {
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
		name               string
		expectedReqHeaders map[string]string
		expectedRespStatus int
		shouldFail         bool
		expectedError      error
	}{
		{
			name:               "Returns nil on successful delete",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusNoContent,
			shouldFail:         false,
			expectedError:      nil,
		},
		{
			name:               "Returns ErrorNotFound when group does not exist",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusNotFound,
			shouldFail:         false,
			expectedError:      r.ErrorNotFound,
		},
		{
			name:               "Returns an error on unexpected status code",
			expectedReqHeaders: map[string]string{"Content-Type": "application/json"},
			expectedRespStatus: http.StatusInternalServerError,
			shouldFail:         false,
			expectedError:      fmt.Errorf("unexpected response status: %d", http.StatusInternalServerError),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var hs http.Header
			mockClient := u.NewMockHttpClient(tc.expectedRespStatus, func(req *http.Request) {
				hs = req.Header
			})
			repo := r.NewD1GroupRepository("https://example.com/api/v1/groups", mockClient)

			group := &g.Group{Id: "C1234f49365c6b492b337189e3343a9d9"}

			err := repo.Destroy(group)

			for k, v := range tc.expectedReqHeaders {
				assert.Equal(t, v, hs.Get(k), "Expected header %s to be %s, but got %s", k, v, hs.Get(k))
			}

			if tc.expectedError == nil {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
			} else {
				assert.Error(t, err, "Expected an error but got none for status code %d", tc.expectedRespStatus)
			}
		})
	}
}
