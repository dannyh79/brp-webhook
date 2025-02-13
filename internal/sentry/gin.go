package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func SentryMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{})
}

func TagBy(ctx *gin.Context, name string) {
	tx := sentry.TransactionFromContext(ctx.Request.Context())
	if tx != nil {
		tx.SetTag("handler", name)
	}
}
