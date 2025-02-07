package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRegistrationService struct {
	shouldFail bool
}

func (m *mockRegistrationService) Execute(g *groups.Group) error {
	if m.shouldFail {
		return fmt.Errorf("failed to register group")
	}
	return nil
}

type mockReplyService struct {
	shouldFail  bool
	calledTimes int
}

func (m *mockReplyService) Execute(replyToken *string) error {
	m.calledTimes++
	if m.shouldFail {
		return fmt.Errorf("failed to send reply")
	}
	return nil
}

func (m *mockReplyService) CalledTimes() int {
	return m.calledTimes
}

func setupRouter(sCtx *services.ServiceContext) *gin.Engine {
	router := gin.New()

	router.POST("/callback",
		routes.LineMsgEventsHandler,
		routes.LineGroupRegistrationHandler(sCtx),
		routes.LineReplyHandler(sCtx),
	)

	return router
}

func TestLineHandlers(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        map[string]interface{}
		expectStatus       int
		shouldRegisterFail bool
		shouldReplyFail    bool
		expectedReplies    int
	}{
		{
			name: "Successful message processing",
			requestBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "message",
						"message": map[string]interface{}{
							"text":       "請好好靈修每日推播靈修內容到這",
							"replyToken": "test-reply-token",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
					},
				},
			},
			expectStatus:       http.StatusOK,
			shouldRegisterFail: false,
			shouldReplyFail:    false,
			expectedReplies:    1,
		},
		{
			name: "Group registration fails",
			requestBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "message",
						"message": map[string]interface{}{
							"text":       "請好好靈修每日推播靈修內容到這",
							"replyToken": "test-reply-token",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
					},
				},
			},
			expectStatus:       http.StatusOK,
			shouldRegisterFail: true,
			shouldReplyFail:    false,
			expectedReplies:    0,
		},
		{
			name: "Reply service fails",
			requestBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "message",
						"message": map[string]interface{}{
							"text":       "請好好靈修每日推播靈修內容到這",
							"replyToken": "test-reply-token",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
					},
				},
			},
			expectStatus:       http.StatusOK,
			shouldRegisterFail: false,
			shouldReplyFail:    true,
			expectedReplies:    1,
		},
		{
			name:            "Invalid event type",
			requestBody:     map[string]interface{}{"events": []map[string]interface{}{}},
			expectStatus:    http.StatusOK,
			expectedReplies: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			replyService := &mockReplyService{shouldFail: tc.shouldReplyFail}
			sCtx := &services.ServiceContext{
				RegistrationService: &mockRegistrationService{shouldFail: tc.shouldRegisterFail},
				ReplyService:        replyService,
			}

			router := setupRouter(sCtx)
			body, _ := json.Marshal(tc.requestBody)

			req := httptest.NewRequest("POST", "/callback", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			u.AssertHttpStatus(t)(rr, tc.expectStatus)
			assert.Equal(t, tc.expectedReplies, replyService.CalledTimes(), "Unexpected number of replies triggered")
		})
	}
}
