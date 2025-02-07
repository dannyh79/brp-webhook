package routes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

type channelSecret = string

func AddRoutes(r *gin.Engine, cs channelSecret, sCtx *services.ServiceContext) {
	r.Use(lineAuthMiddleware(cs))
	r.POST("/api/v1/callback",
		msgEventsHandler,
		groupRegistrationHandler(sCtx),
		replyHandler(sCtx),
		successHandler,
	)
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
		// Writes the request body back after inspection
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

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

type groupDto struct {
	*groups.Group
	ReplyToken string
}

func msgEventsHandler(ctx *gin.Context) {
	defer ctx.Request.Body.Close()

	var b LineCallbackBody
	if err := ctx.ShouldBindJSON(&b); err != nil {
		fmt.Printf("Error in unmarshalling request body: %v", err)
		ctx.Next()
		return
	}

	var gs []*groupDto
	for _, e := range b.Events {
		if e.Type == "message" && e.Message.Text == RegisterMyGroupMsg && len(e.Message.ReplyToken) > 0 {
			gs = append(gs, &groupDto{Group: groups.NewGroup(e.Source.GroupId), ReplyToken: e.Message.ReplyToken})
		}
	}

	ctx.Set("groups", gs)

	ctx.Next()
}

func groupRegistrationHandler(sCtx *services.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gsIf, exists := ctx.Get("groups")
		if !exists {
			ctx.Next()
			return
		}

		gs, ok := gsIf.([]*groupDto)
		if !ok {
			ctx.Next()
			return
		}

		var registered []*groupDto
		for _, g := range gs {
			if err := sCtx.RegistrationService.Execute(g.Group); err != nil {
				fmt.Printf("Error in registering group: %v", err)
			} else {
				registered = append(registered, g)
			}
		}

		ctx.Set("group", registered)

		ctx.Next()
	}
}

func replyHandler(sCtx *services.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gsIf, exists := ctx.Get("groups")
		if !exists {
			ctx.Next()
			return
		}

		gs, ok := gsIf.([]*groupDto)
		if !ok {
			ctx.Next()
			return
		}

		for _, g := range gs {
			if err := sCtx.ReplyService.Execute(g.ReplyToken); err != nil {
				fmt.Printf("Error in replying to completed registration for group %v via LINE: %v", g.Id, err)
			}
		}

		ctx.Next()
	}
}

func successHandler(ctx *gin.Context) { ctx.Status(200) }
