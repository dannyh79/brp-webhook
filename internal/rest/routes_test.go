package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	s "github.com/dannyh79/brp-webhook/internal/services"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
)

var textMessageEvent = routes.MessageEvent{
	Event: routes.Event{
		Type: "message",
		Source: routes.Source{
			Type:    "user",
			GroupId: "C1234f49365c6b492b337189e3343a9d9",
			UserId:  "U123425e31582f9bdc77b386c1d02477e",
		},
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
	},
	Message: routes.MessageEventBody{
		Type: "text",
		Text: routes.RegisterMyGroupMsg,
	},
}

func Test_POSTCallback(t *testing.T) {
	u.InitRoutesTest()

	textMessageEventString, _ := json.Marshal(textMessageEvent)
	fmt.Printf(`{"events":[%s]}`, textMessageEventString)

	tcs := []struct {
		name                string
		hasReqHead          bool
		hasInvalidSignature bool
		reqBody             []byte
		statusCode          int
	}{
		{
			name:       "Returns 200",
			hasReqHead: true,
			reqBody:    []byte(`{"events":[]}`),
			statusCode: 200,
		},
		{
			name:       "Returns 200 when there is a text message event",
			hasReqHead: true,
			reqBody:    []byte(fmt.Sprintf(`{"events":[%s]}`, textMessageEventString)),
			statusCode: 200,
		},
		{
			name:       "Returns 400 when the request is missing body",
			hasReqHead: true,
			statusCode: 400,
		},
		{
			name:       "Returns 401 when the request is missing signature header",
			reqBody:    []byte(`{"foo":"bar"}`),
			statusCode: 401,
		},
		{
			name:                "Returns 401 when the request is not authorized",
			hasReqHead:          true,
			hasInvalidSignature: true,
			reqBody:             []byte(`{"events":[]}`),
			statusCode:          401,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			regS := u.NewMockService[s.GroupDto](false)
			repS := u.NewMockService[s.GroupDto](false)
			sCtx := &s.ServiceContext{
				RegistrationService: regS,
				ReplyService:        repS,
			}

			suite := u.NewRoutesTestSuite()
			routes.AddRoutes(suite.Router, u.StubSecret, sCtx)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/callback", bytes.NewBuffer(tc.reqBody))
			if tc.hasReqHead && len(tc.reqBody) > 0 {
				s := u.GenerateSignature(u.StubSecret, string(tc.reqBody))
				if tc.hasInvalidSignature {
					s = u.GenerateSignature("some-invalid-line-channel-secret", string(tc.reqBody))
				}

				req.Header.Add("x-line-signature", s)
			}

			rr := httptest.NewRecorder()
			suite.Router.ServeHTTP(rr, req)

			u.AssertHttpStatus(t)(rr, tc.statusCode)
		})
	}
}
