package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/gin-gonic/gin"
)

func Test_POSTCallback(t *testing.T) {
	tcs := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Returns 200",
			statusCode: 200,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := newTestSuite()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/callback", nil)

			suite.Router.ServeHTTP(rr, req)

			assertHttpStatus(t)(rr, tc.statusCode)
		})
	}
}

type testSuite struct {
	Router *gin.Engine
}

func newTestSuite() *testSuite {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.AddRoutes(r)

	return &testSuite{
		Router: r,
	}
}

func assertHttpStatus(t *testing.T) func(rr *httptest.ResponseRecorder, want int) {
	return func(rr *httptest.ResponseRecorder, want int) {
		t.Helper()
		if got := rr.Result().StatusCode; got != want {
			t.Errorf("got HTTP status %v, want %v", got, want)
		}
	}
}
