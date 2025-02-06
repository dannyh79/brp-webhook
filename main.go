package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dannyh79/brp-webhook/internal/repositories"
	routes "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	secret := os.Getenv("LINE_CHANNEL_SECRET")
	endpoint := os.Getenv("D1_GROUP_QUERY_ENDPOINT")

	repo := repositories.NewD1GroupRepository(endpoint, &http.Client{})
	service := services.NewRegistrationService(repo)

	routes.AddRoutes(router, secret, service)
	err := router.Run()
	if err != nil {
		panic(fmt.Sprintf("Error in starting the app: %v", err))
	}
}
