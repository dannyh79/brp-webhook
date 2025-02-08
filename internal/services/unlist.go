package services

import (
	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
)

var _ Service[GroupDto] = (*UnlistService)(nil)

// TODO: impl
type UnlistService struct {
	repo r.Repository[g.Group]
}

func (s *UnlistService) Execute(g *GroupDto) error {
	return nil
}

func NewUnlistService(r r.Repository[g.Group]) Service[GroupDto] {
	return &UnlistService{r}
}
