package services

import (
	"net/http"

	g "github.com/dannyh79/brp-webhook/internal/groups"
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service[T any] interface {
	Execute(*T) error
}

type ReplyToken = string

type GroupDto struct {
	*g.Group
	ReplyToken
	WasRegistered bool
}

type ServiceContext struct {
	UnlistService       Service[GroupDto]
	RegistrationService Service[GroupDto]
	ReplyService        Service[GroupDto]
	WelcomeService      Service[GroupDto]
}

type ReplyMessageRequest struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []message `json:"messages"`
}

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
