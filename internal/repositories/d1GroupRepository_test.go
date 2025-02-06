package repositories_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/repositories"
)

func TestSaveGroupParams_MarshalJSON(t *testing.T) {
	tcs := []struct {
		name           string
		input          repositories.SaveGroupParams
		expectedOutput string
	}{
		{
			name:           "Marshals to JSON",
			input:          repositories.SaveGroupParams{Id: "C1234f49365c6b492b337189e3343a9d9"},
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
		name        string
		statusCode  int
		expectError bool
	}{
		{"Returns 204 (Success)", http.StatusNoContent, false},
		{"Returns 304 (Success)", http.StatusNotModified, false},
		{"Returns 400 (Failure)", http.StatusBadRequest, true},
		{"Returns 500 (Failure)", http.StatusInternalServerError, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := &http.Client{
				Transport: &mockHTTPTransport{statusCode: tc.statusCode},
			}

			repo := repositories.NewD1GroupRepository("https://example.com/api/v1/groups", mockClient)

			group := &groups.Group{Id: "C1234f49365c6b492b337189e3343a9d9"}
			result, err := repo.Save(group)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none for status code %d", tc.statusCode)
				}
				if result != nil {
					t.Errorf("Expected result to be nil on error, but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got one: %v", err)
				}
				if result == nil {
					t.Fatal("Expected result to be non-nil but got nil")
				}
				if result.Id != group.Id {
					t.Errorf("Returned group does not match input: got %v, want %v", result, group)
				}
			}
		})
	}
}

type mockHTTPTransport struct {
	statusCode int
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
		Header:     make(http.Header),
	}, nil
}
