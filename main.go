package main

import (
	"fmt"
	"os"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	secret := os.Getenv("LINE_CHANNEL_SECRET")
	routes.AddRoutes(router, secret)
	err := router.Run()
	if err != nil {
		panic(fmt.Sprintf("Error in starting the app: %v", err))
	}
}
