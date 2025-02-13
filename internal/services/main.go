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

type ServiceContext struct {
	UnlistService       Service[GroupDto]
	RegistrationService Service[GroupDto]
	ReplyService        Service[GroupDto]
	WelcomeService      Service[GroupDto]
}

func NewServiceContext(unlistService, registrationService, replyService, welcomeService Service[GroupDto]) *ServiceContext {
	return &ServiceContext{
		UnlistService:       unlistService,
		RegistrationService: registrationService,
		ReplyService:        replyService,
		WelcomeService:      welcomeService,
	}
}

type ReplyToken = string

type GroupDto struct {
	*g.Group
	ReplyToken
	WasRegistered bool
}

func NewGroupDto(g *g.Group, t ReplyToken) *GroupDto {
	return &GroupDto{Group: g, ReplyToken: t}
}
