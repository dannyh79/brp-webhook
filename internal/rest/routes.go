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
	r.POST("/api/v1/callback", func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.Status(400)
			return
		}

		defer ctx.Request.Body.Close()

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil || len(body) == 0 {
			ctx.Status(400)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(ctx.Request.Header.Get("x-line-signature"))
		if err != nil {
			ctx.Status(401)
			return
		}

		hash := hmac.New(sha256.New, []byte(s))
		hash.Write(body)
		if !hmac.Equal(decoded, hash.Sum(nil)) {
			ctx.Status(401)
			return
		}

		ctx.Status(200)
	})
}
