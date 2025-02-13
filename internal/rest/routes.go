package routes

import (
	"github.com/dannyh79/brp-webhook/internal/sentry"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine, cs channelSecret, sCtx *s.ServiceContext) {
	r.Use(LineAuthMiddleware(cs))
	r.POST("/api/v1/callback",
		lineEventsHandler(sCtx),
		successHandler,
	)
}

func successHandler(ctx *gin.Context) {
	sentry.TagBy(ctx, "successHandler")
	ctx.Status(200)
}
