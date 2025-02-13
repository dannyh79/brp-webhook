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

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			unlS := u.NewMockService[s.GroupDto](tc.shouldUnlistFail)
			regS := u.NewMockService[s.GroupDto](tc.shouldRegisterFail, tc.registerFailError)
			repS := u.NewMockService[s.GroupDto](tc.shouldReplyFail)
			welS := u.NewMockService[s.GroupDto](tc.shouldWelcomeFail)
			sCtx := &s.ServiceContext{UnlistService: unlS, RegistrationService: regS, ReplyService: repS, WelcomeService: welS}

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
			assert.Equal(t, tc.expectedWelcomes, welS.CalledTimes(), "Unexpected number of welcomes triggered")
			assert.Equal(t, tc.expectedUnlistings, unlS.CalledTimes(), "Unexpected number of unlistings triggered")
			assert.Equal(t, tc.expectedRegistrations, regS.CalledTimes(), "Unexpected number of registrations triggered")
			assert.Equal(t, tc.expectedReplies, repS.CalledTimes(), "Unexpected number of replies triggered")
		})
	}
}
