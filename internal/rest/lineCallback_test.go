package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/rest"
	s "github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRegistrationService struct {
	shouldFail bool
	error
}

func (m *mockRegistrationService) Execute(g *s.GroupDto) error {
	if !m.shouldFail {
		return nil
	}
	if m.error == nil {
		m.error = fmt.Errorf("failed to register group")
	}
	return m.error
}

type mockReplyService struct {
	shouldFail  bool
	calledTimes int
}

func (m *mockReplyService) Execute(g *s.GroupDto) error {
	m.calledTimes++
	if m.shouldFail {
		return fmt.Errorf("failed to send reply")
	}
	return nil
}

func (m *mockReplyService) CalledTimes() int {
	return m.calledTimes
}

func setupRouter(sCtx *s.ServiceContext) *gin.Engine {
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
		registerFailError  error
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
							"text": "請好好靈修每日推播靈修內容到這",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
						"replyToken": "test-reply-token",
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
							"text": "請好好靈修每日推播靈修內容到這",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
						"replyToken": "test-reply-token",
					},
				},
			},
			expectStatus:       http.StatusOK,
			shouldRegisterFail: true,
			shouldReplyFail:    false,
			expectedReplies:    0,
		},
		{
			name: "Group registration fails from record already exists",
			requestBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "message",
						"message": map[string]interface{}{
							"text": "請好好靈修每日推播靈修內容到這",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
						"replyToken": "test-reply-token",
					},
				},
			},
			expectStatus:       http.StatusOK,
			shouldRegisterFail: true,
			registerFailError:  s.ErrorGroupAlreadyRegistered,
			shouldReplyFail:    false,
			expectedReplies:    1,
		},
		{
			name: "Reply service fails",
			requestBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "message",
						"message": map[string]interface{}{
							"text": "請好好靈修每日推播靈修內容到這",
						},
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
						"replyToken": "test-reply-token",
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
			sCtx := &s.ServiceContext{
				RegistrationService: &mockRegistrationService{shouldFail: tc.shouldRegisterFail, error: tc.registerFailError},
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
