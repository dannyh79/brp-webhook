package services

import g "github.com/dannyh79/brp-webhook/internal/groups"

type Service[T any] interface {
	Execute(*T) error
}

type ServiceContext struct {
	RegistrationService Service[g.Group]
	ReplyService        Service[string]
}

func NewServiceContext(regS Service[g.Group], replyS Service[string]) *ServiceContext {
	return &ServiceContext{
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
	return &GroupDto{Group: g}
}
