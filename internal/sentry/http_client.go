package sentry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryHttpClient struct {
	client *http.Client
}

func NewSentryHttpClient(client *http.Client) *SentryHttpClient {
	return &SentryHttpClient{client: client}
}

func (s *SentryHttpClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	tx := sentry.TransactionFromContext(ctx)

	if tx == nil {
		tx = sentry.StartTransaction(ctx, "Default Transaction")
		defer tx.Finish()
	}

	span := tx.StartChild("http.client", sentry.WithTransactionName(req.Method+" "+req.URL.String()))
	defer span.Finish()

	span.SetTag("http.method", req.Method)
	span.SetTag("http.url", req.URL.String())

	startTime := time.Now()
	resp, err := s.client.Do(req)
	duration := time.Since(startTime)

	span.SetData("duration_ms", duration.Milliseconds())

	if err != nil {
		sentry.CaptureException(err)
		span.Status = sentry.SpanStatusInternalError
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	span.SetTag("http.status_code", fmt.Sprintf("%d", resp.StatusCode))

	if resp.StatusCode >= 400 {
		sentry.CaptureMessage(fmt.Sprintf("HTTP request failed: %s %s (%d)", req.Method, req.URL.String(), resp.StatusCode))
		span.Status = sentry.SpanStatusUnknown
	}

	return resp, nil
}
