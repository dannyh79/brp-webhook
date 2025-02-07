package routes

import (
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine, cs channelSecret, sCtx *s.ServiceContext) {
	r.Use(LineAuthMiddleware(cs))
	r.POST("/api/v1/callback",
		LineMsgEventsHandler,
		LineGroupRegistrationHandler(sCtx),
		LineReplyHandler(sCtx),
		successHandler,
	)
}

func successHandler(ctx *gin.Context) { ctx.Status(200) }
