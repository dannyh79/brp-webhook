package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	s "github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_POSTCallback(t *testing.T) {
	u.InitRoutesTest()

	tcs := []struct {
		name         string
		expectStatus int

		reqBody             map[string]interface{}
		noSignatureHead     bool
		hasInvalidSignature bool

		expectedUnlistings    int
		shouldUnlistFail      bool
		shouldRegisterFail    bool
		expectedRegistrations int
		registerFailError     error
		expectedReplies       int
		shouldReplyFail       bool
	}{
		//#region authentication for LINE requests
		{
			name:         "Returns 400 when the request is missing body",
			expectStatus: http.StatusBadRequest,
		},
		{
			name:            "Returns 401 when the request is missing signature header",
			expectStatus:    http.StatusUnauthorized,
			reqBody:         map[string]interface{}{},
			noSignatureHead: true,
		},
		{
			name:                "Returns 401 when the request is not authorized",
			expectStatus:        http.StatusUnauthorized,
			hasInvalidSignature: true,
			reqBody:             map[string]interface{}{},
		},
		//#endregion

		{
			name:         `Returns 200 when there is no events in request body`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
				"events": []map[string]interface{}{},
			},
			expectedUnlistings:    0,
			expectedRegistrations: 0,
			expectedReplies:       0,
		},
		{
			name:         `Returns 200 when receiving "請好好靈修每日推播靈修內容到這" text message event from a group`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
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
			expectedUnlistings:    0,
			expectedRegistrations: 1,
			expectedReplies:       1,
		},
		{
			name:         `Returns 200 when receiving "請好好靈修每日推播靈修內容到這" text message event & leave event from groups`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
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
					{
						"type": "leave",
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
					},
				},
			},
			expectedUnlistings:    1,
			expectedRegistrations: 1,
			expectedReplies:       1,
		},
		{
			name:         `Returns 200 when receiving "請好好靈修每日推播靈修內容到這" text message event from a group - registration failed from having existing record`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
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
			expectedUnlistings:    0,
			expectedRegistrations: 1,
			shouldRegisterFail:    true,
			registerFailError:     s.ErrorGroupAlreadyRegistered,
			expectedReplies:       1,
		},
		{
			name:         `Returns 200 when receiving "請好好靈修每日推播靈修內容到這" text message event from a group - registration failed for any other reason`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
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
			expectedUnlistings:    0,
			expectedRegistrations: 1,
			shouldRegisterFail:    true,
			expectedReplies:       0,
		},
		{
			name:         `Returns 200 when receiving "請好好靈修每日推播靈修內容到這" text message event from a group - reply failed`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
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
			expectedUnlistings:    0,
			expectedRegistrations: 1,
			shouldRegisterFail:    false,
			expectedReplies:       1,
			shouldReplyFail:       true,
		},
		{
			name:         `Returns 200 when receiving leave event`,
			expectStatus: http.StatusOK,
			reqBody: map[string]interface{}{
				"events": []map[string]interface{}{
					{
						"type": "leave",
						"source": map[string]interface{}{
							"groupId": "C1234",
						},
					},
				},
			},
			expectedUnlistings:    1,
			shouldUnlistFail:      false,
			expectedRegistrations: 0,
			shouldRegisterFail:    false,
			expectedReplies:       0,
			shouldReplyFail:       false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			unlS := u.NewMockService[s.GroupDto](tc.shouldUnlistFail)
			regS := u.NewMockService[s.GroupDto](tc.shouldRegisterFail, tc.registerFailError)
			repS := u.NewMockService[s.GroupDto](tc.shouldReplyFail)
			sCtx := &s.ServiceContext{UnlistService: unlS, RegistrationService: regS, ReplyService: repS}

			suite := u.NewRoutesTestSuite()
			routes.AddRoutes(suite.Router, u.StubSecret, sCtx)

			body, _ := json.Marshal(tc.reqBody)
			var req *http.Request
			if tc.reqBody == nil {
				req, _ = http.NewRequest(http.MethodPost, "/api/v1/callback", nil)
			} else {
				req, _ = http.NewRequest(http.MethodPost, "/api/v1/callback", bytes.NewBuffer(body))
			}

			if !tc.noSignatureHead {
				s := u.GenerateSignature(u.StubSecret, string(body))
				if tc.hasInvalidSignature {
					s = u.GenerateSignature("some-invalid-line-channel-secret", string(body))
				}
				req.Header.Add("x-line-signature", s)
			}

			rr := httptest.NewRecorder()
			suite.Router.ServeHTTP(rr, req)

			u.AssertHttpStatus(t)(rr, tc.expectStatus)
			assert.Equal(t, tc.expectedUnlistings, unlS.CalledTimes(), "Unexpected number of unlistings triggered")
			assert.Equal(t, tc.expectedRegistrations, regS.CalledTimes(), "Unexpected number of registrations triggered")
			assert.Equal(t, tc.expectedReplies, repS.CalledTimes(), "Unexpected number of replies triggered")
		})
	}
}
