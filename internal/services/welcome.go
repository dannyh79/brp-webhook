package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dannyh79/brp-webhook/internal/sentry"
)

var _ Service[GroupDto] = (*WelcomeService)(nil)

type WelcomeService struct {
	token  string
	client HttpDoer
}

func NewWelcomeService(token string, client *http.Client) Service[GroupDto] {
	return &WelcomeService{token: token, client: sentry.NewSentryHttpClient(client)}
}

func (s *WelcomeService) Execute(g *GroupDto) error {
	p := ReplyMessageRequest{
		ReplyToken: g.ReplyToken,
		Messages: []message{
			{Type: "text", Text: msgWelcome},
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal reply JSON: %w", err)
	}

	req, err := http.NewRequest("POST", lineReplyApiEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create welcome request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send welcome request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from welcoming: %d", resp.StatusCode)
	}

	return nil
}
