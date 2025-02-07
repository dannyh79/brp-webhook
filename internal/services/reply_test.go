package services_test

import (
	"net/http"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
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

			mockClient := u.NewMockHttpClient(tc.mockRespStatusCode)
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
