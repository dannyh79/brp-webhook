package main

import (
	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.AddRoutes(router)
	router.Run()
}
