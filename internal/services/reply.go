package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dannyh79/brp-webhook/internal/sentry"
)

var _ Service[GroupDto] = (*ReplyService)(nil)

type ReplyService struct {
	token  string
	client HttpDoer
}

func NewReplyService(token string, client *http.Client) Service[GroupDto] {
	return &ReplyService{token: token, client: sentry.NewSentryHttpClient(client)}
}

func (s *ReplyService) Execute(g *GroupDto) error {
	m := msgOk
	if g.WasRegistered {
		m = msgAlreadyRegistered
	}
	p := ReplyMessageRequest{
		ReplyToken: g.ReplyToken,
		Messages: []message{
			{Type: "text", Text: m},
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal reply JSON: %w", err)
	}

	req, err := http.NewRequest("POST", lineReplyApiEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create reply request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send reply request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from replying: %d", resp.StatusCode)
	}

	return nil
}
