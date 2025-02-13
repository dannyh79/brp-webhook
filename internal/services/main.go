package services

import (
	g "github.com/dannyh79/brp-webhook/internal/groups"
)

func NewServiceContext(unlistService, registrationService, replyService, welcomeService Service[GroupDto]) *ServiceContext {
	return &ServiceContext{
		UnlistService:       unlistService,
		RegistrationService: registrationService,
		ReplyService:        replyService,
		WelcomeService:      welcomeService,
	}
}

func NewGroupDto(g *g.Group, t ReplyToken) *GroupDto {
	return &GroupDto{Group: g, ReplyToken: t}
}
