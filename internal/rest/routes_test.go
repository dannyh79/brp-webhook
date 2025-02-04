package routes_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/gin-gonic/gin"
)

const stubSecret = "some-line-channel-secret"

func Test_POSTCallback(t *testing.T) {
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
			// TODO
			reqBody:    []byte(`{"foo":"bar"}`),
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
			// TODO
			reqBody:    []byte(`{"foo":"bar"}`),
			statusCode: 401,
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

func newTestSuite(s string) *testSuite {
	r := gin.New()
	routes.AddRoutes(r, s)

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
		if got := rr.Result().StatusCode; got != want {
			t.Errorf("got HTTP status %v, want %v", got, want)
		}
	}
}
