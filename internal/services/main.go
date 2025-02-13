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
}

func NewServiceContext(unlS, regS, replyS Service[GroupDto]) *ServiceContext {
	return &ServiceContext{
		UnlistService:       unlS,
		RegistrationService: regS,
		ReplyService:        replyS,
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
