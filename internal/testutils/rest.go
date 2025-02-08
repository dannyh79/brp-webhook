package testutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

func GenerateSignature(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type RoutesTestSuite struct {
	Router         *gin.Engine
	ServiceContext *s.ServiceContext
}

const StubSecret = "some-line-channel-secret"

// Set test mode for gin
func InitRoutesTest() {
	gin.SetMode(gin.TestMode)
}

func NewRoutesTestSuite() *RoutesTestSuite {
	r := gin.New()

	return &RoutesTestSuite{
		Router: r,
	}
}
