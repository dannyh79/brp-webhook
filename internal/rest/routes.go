package routes

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine) {
	r.POST("/api/v1/callback", func(ctx *gin.Context) {
		ctx.Status(200)
	})
}
