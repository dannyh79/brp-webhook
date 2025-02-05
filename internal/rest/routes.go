package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/gin-gonic/gin"
)

type channelSecret = string

func AddRoutes(r *gin.Engine, s channelSecret) {
	r.Use(lineAuthMiddleware(s))
	r.POST("/api/v1/callback", func(ctx *gin.Context) {
		ctx.Status(200)
	})
}

func lineAuthMiddleware(s channelSecret) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.AbortWithStatus(400)
		}

		defer ctx.Request.Body.Close()

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil || len(body) == 0 {
			ctx.AbortWithStatus(400)
		}

		decoded, err := base64.StdEncoding.DecodeString(ctx.Request.Header.Get("x-line-signature"))
		if err != nil {
			ctx.AbortWithStatus(401)
		}

		hash := hmac.New(sha256.New, []byte(s))
		hash.Write(body)
		if !hmac.Equal(decoded, hash.Sum(nil)) {
			ctx.AbortWithStatus(401)
		}

		ctx.Next()
	}
}
