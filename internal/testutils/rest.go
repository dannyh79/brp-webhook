package testutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSignature(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
