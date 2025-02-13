package services

import (
	"net/http"
)

var _ Service[GroupDto] = (*WelcomeService)(nil)

type WelcomeService struct {
	token  string
	client HttpDoer
}

func NewWelcomeService(token string, client *http.Client) Service[GroupDto] {
	return &WelcomeService{token: token, client: client}
}

func (s *WelcomeService) Execute(g *GroupDto) error {
	return nil
}
