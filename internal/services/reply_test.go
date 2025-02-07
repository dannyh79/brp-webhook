package services_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	s "github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/stretchr/testify/assert"
)

const stubChannelToken = "some-line-channel-token"

func Test_ReplyService(t *testing.T) {
	tcs := []struct {
		name               string
		dto                s.GroupDto
		expectedMsg        string
		replyToken         string
		mockRespStatusCode int
		expectError        bool
	}{
		{
			name:               "Replies with MsgOk for new registration",
			dto:                s.GroupDto{ReplyToken: "test-reply-token", WasRegistered: false},
			expectedMsg:        "已加入推播清單。",
			replyToken:         "test-reply-token",
			mockRespStatusCode: http.StatusOK,
			expectError:        false,
		},
		{
			name:               "Replies with MsgAlreadyRegistered for existing registration",
			dto:                s.GroupDto{ReplyToken: "test-reply-token", WasRegistered: true},
			expectedMsg:        "已在加入推播清單中。 :D",
			replyToken:         "test-reply-token",
			mockRespStatusCode: http.StatusOK,
			expectError:        false,
		},
		{
			name:               "Returns error on bad request",
			dto:                s.GroupDto{ReplyToken: "test-reply-token"},
			replyToken:         "test-reply-token",
			mockRespStatusCode: http.StatusBadRequest,
			expectError:        true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var reqBody []byte
			mockClient := u.NewMockHttpClient(tc.mockRespStatusCode, func(req *http.Request) {
				b, _ := req.GetBody()
				reqBody, _ = io.ReadAll(b)
			})

			svc := s.NewReplyService(stubChannelToken, mockClient)
			err := svc.Execute(&tc.dto)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)

				var sent s.ReplyMessageRequest
				err := json.Unmarshal(reqBody, &sent)
				assert.NoError(t, err, "Failed to unmarshal sent request body")

				assert.Equal(t, tc.replyToken, sent.ReplyToken, "Reply token mismatch")

				assert.Len(t, sent.Messages, 1, "Expected one message object to be sent")
				assert.Equal(t, "text", sent.Messages[0].Type, `Expected message type to be "text"`)
				assert.Equal(t, tc.expectedMsg, sent.Messages[0].Text, `Expect "%s", got "%s"`, sent.Messages[0].Text, tc.expectedMsg)
			}
		})
	}
}
