package sentry

import (
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func SentryMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{})
}
