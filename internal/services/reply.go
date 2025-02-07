package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var _ Service[GroupDto] = (*ReplyService)(nil)

const lineReplyApiEndpoint = "https://api.line.me/v2/bot/message/reply"

type replyMessageRequest struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []message `json:"messages"`
}

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ReplyService struct {
	token  string
	client *http.Client
}

func NewReplyService(token string, client *http.Client) Service[GroupDto] {
	return &ReplyService{token: token, client: client}
}

func (s *ReplyService) Execute(g *GroupDto) error {
	p := replyMessageRequest{
		ReplyToken: g.ReplyToken,
		Messages: []message{
			{Type: "text", Text: "好的"},
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from replying: %d", resp.StatusCode)
	}

	return nil
}
