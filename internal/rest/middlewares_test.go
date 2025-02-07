package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	u "github.com/dannyh79/brp-webhook/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLineAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name            string
		body            string
		secret          string
		modifySignature bool
		expectStatus    int
	}{
		{
			name:            "Valid request - Body remains available",
			body:            `{"message":"hello"}`,
			secret:          "test-secret",
			modifySignature: false,
			expectStatus:    http.StatusOK,
		},
		{
			name:            "Invalid signature - Unauthorized",
			body:            `{"message":"hello"}`,
			secret:          "test-secret",
			modifySignature: true,
			expectStatus:    http.StatusUnauthorized,
		},
		{
			name:            "Missing body - Bad Request",
			body:            ``,
			secret:          "test-secret",
			modifySignature: false,
			expectStatus:    http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			router := gin.New()
			router.Use(routes.LineAuthMiddleware(tc.secret))
			router.POST("/test", func(ctx *gin.Context) {
				body, err := io.ReadAll(ctx.Request.Body)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read body"})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"body": string(body)})
			})

			signature := u.GenerateSignature(tc.secret, tc.body)
			if tc.modifySignature {
				signature = "some-invalid-signature"
			}

			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-line-signature", signature)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			u.AssertHttpStatus(t)(rr, tc.expectStatus)

			if tc.expectStatus == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err, "Failed to parse JSON response")
				assert.Equal(t, tc.body, resp["body"], "Request body was lost or modified")
			}
		})
	}
}
