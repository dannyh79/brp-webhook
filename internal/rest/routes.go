package routes

import (
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine, cs channelSecret, sCtx *s.ServiceContext) {
	r.Use(lineAuthMiddleware(cs))
	r.POST("/api/v1/callback",
		msgEventsHandler,
		groupRegistrationHandler(sCtx),
		replyHandler(sCtx),
		successHandler,
	)
}

func successHandler(ctx *gin.Context) { ctx.Status(200) }
