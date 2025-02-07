package services_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/stretchr/testify/assert"
)

const stubChannelToken = "some-line-channel-token"

func Test_ReplyService(t *testing.T) {
	tcs := []struct {
		name               string
		replyToken         string
		mockRespStatusCode int
		expectError        bool
	}{
		{
			name:               "Does not return error",
			replyToken:         "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
			mockRespStatusCode: http.StatusOK,
			expectError:        false,
		},
		{
			name:               "Returns error",
			replyToken:         "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
			mockRespStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockClient := &http.Client{
				Transport: &mockHTTPTransport{statusCode: tc.mockRespStatusCode},
			}

			s := services.NewReplyService(stubChannelToken, mockClient)
			err := s.Execute(tc.replyToken)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
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
