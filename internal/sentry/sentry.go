package sentry

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type Monitor struct {
	env string
	dsn string
}

func NewMonitor(env, dsn string) *Monitor {
	return &Monitor{env: env, dsn: dsn}
}

func (s *Monitor) isProd() bool {
	return s.env == "production"
}

func (s *Monitor) Init() {
	if !s.isProd() {
		log.Printf("Sentry disabled in %s environment", s.env)
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              s.dsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0, // Capture all requests
	})
	if err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
	}
}

func (s *Monitor) Flush() {
	if !s.isProd() {
		return
	}

	sentry.Flush(2 * time.Second)
}
