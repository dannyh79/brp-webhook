package services

import "github.com/dannyh79/brp-webhook/internal/groups"

type Service[T any] interface {
	Execute(*T) error
}

type ServiceContext struct {
	RegistrationService Service[groups.Group]
	ReplyService        Service[string]
}

func NewServiceContext(regS Service[groups.Group], replyS Service[string]) *ServiceContext {
	return &ServiceContext{
		RegistrationService: regS,
		ReplyService:        replyS,
	}
}
