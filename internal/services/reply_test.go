package services_test

import (
	"testing"

	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/stretchr/testify/assert"
)

func Test_ReplyService(t *testing.T) {
	tcs := []struct {
		name        string
		replyToken  string
		expectError bool
	}{
		{
			name:        "Does not return error",
			replyToken:  "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
			expectError: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := services.NewReplyService()
			err := s.Execute(tc.replyToken)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
			}
		})
	}
}
