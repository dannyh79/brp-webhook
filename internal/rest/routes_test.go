package routes_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/groups"
	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const stubSecret = "some-line-channel-secret"

var textMessageEvent = routes.MessageEvent{
	Event: routes.Event{
		Type: "message",
		Source: routes.Source{
			Type:    "user",
			GroupId: "C1234f49365c6b492b337189e3343a9d9",
			UserId:  "U123425e31582f9bdc77b386c1d02477e",
		},
	},
	Message: routes.MessageEventBody{
		Type:       "text",
		Text:       routes.RegisterMyGroupMsg,
		ReplyToken: "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
	},
}

func Test_POSTCallback(t *testing.T) {
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

			suite := newTestSuite(stubSecret)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/callback", bytes.NewBuffer(tc.reqBody))
			if tc.hasReqHead && len(tc.reqBody) > 0 {
				s := generateSignature(stubSecret, tc.reqBody)
				if tc.hasInvalidSignature {
					s = generateSignature("some-invalid-line-channel-secret", tc.reqBody)
				}

				req.Header.Add("x-line-signature", s)
			}

			suite.Router.ServeHTTP(rr, req)

			assertHttpStatus(t)(rr, tc.statusCode)
		})
	}
}

func init() {
	gin.SetMode(gin.TestMode)
}

type testSuite struct {
	Router *gin.Engine
}

type stubGroupRepo struct{}

func (r *stubGroupRepo) Save(g *groups.Group) (*groups.Group, error) {
	return g, nil
}

func newTestSuite(cs string) *testSuite {
	r := gin.New()

	s := services.NewRegistrationService(&stubGroupRepo{})
	routes.AddRoutes(r, cs, s)

	return &testSuite{
		Router: r,
	}
}

func generateSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func assertHttpStatus(t *testing.T) func(rr *httptest.ResponseRecorder, want int) {
	return func(rr *httptest.ResponseRecorder, want int) {
		t.Helper()
		got := rr.Result().StatusCode
		assert.Equal(t, got, want, "got HTTP status %v, want %v", got, want)
	}
}
