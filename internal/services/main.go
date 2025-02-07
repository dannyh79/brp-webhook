package services

import g "github.com/dannyh79/brp-webhook/internal/groups"

type Service[T any] interface {
	Execute(*T) error
}

type ServiceContext struct {
	RegistrationService Service[GroupDto]
	ReplyService        Service[GroupDto]
}

func NewServiceContext(regS, replyS Service[GroupDto]) *ServiceContext {
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
	return &GroupDto{Group: g, ReplyToken: t}
}
