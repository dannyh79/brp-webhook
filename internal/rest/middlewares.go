package routes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"

	"github.com/dannyh79/brp-webhook/internal/sentry"
	"github.com/gin-gonic/gin"
)

type channelSecret = string

func LineAuthMiddleware(s channelSecret) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sentry.TagBy(ctx, "LineAuthMiddleware")

		if ctx.Request.ContentLength == 0 {
			ctx.AbortWithStatus(400)
			return
		}

		defer func() {
			err := ctx.Request.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil || len(body) == 0 {
			ctx.AbortWithStatus(400)
			return
		}
		// Writes the request body back after inspection
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		decoded, err := base64.StdEncoding.DecodeString(ctx.Request.Header.Get("x-line-signature"))
		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		hash := hmac.New(sha256.New, []byte(s))
		hash.Write(body)
		if !hmac.Equal(decoded, hash.Sum(nil)) {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Next()
	}
}
