package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

type channelSecret = string

func AddRoutes(r *gin.Engine, s channelSecret) {
	r.Use(lineAuthMiddleware(s))
	r.POST("/api/v1/callback", lineCallbackHandler)
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

const RegisterMyGroupMsg = "請好好靈修每日推播靈修內容到這"

type stubGroupRepo struct{}

func (r *stubGroupRepo) Save(g *groups.Group) (*groups.Group, error) {
	return g, nil
}

func lineCallbackHandler(ctx *gin.Context) {
	defer ctx.Request.Body.Close()

	body, err := io.ReadAll(ctx.Request.Body)
	if err == nil {
		var b LineCallbackBody
		if err := json.Unmarshal(body, &b); err != nil {
			fmt.Printf("Error in unmarshalling request body: %v", err)
		}

		for _, e := range b.Events {
			if e.Type == "message" && e.Message.Text == RegisterMyGroupMsg {
				g := groups.NewGroup(e.Source.GroupId)
				s := services.NewRegistrationService(&stubGroupRepo{})
				if err := s.Execute(g); err != nil {
					fmt.Printf("Error in registering group: %v", err)
				}
			}
		}
	}

	ctx.Status(200)
}
