package main

import (
	"fmt"

	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.AddRoutes(router)
	err := router.Run()
	if err != nil {
		panic(fmt.Sprintf("Error in starting the app: %v", err))
	}
}
