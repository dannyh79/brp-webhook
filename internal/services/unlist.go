package services

import (
	"errors"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
)

var _ Service[GroupDto] = (*UnlistService)(nil)

type UnlistService struct {
	repo r.Repository[g.Group]
}

func (s *UnlistService) Execute(g *GroupDto) error {
	err := s.repo.Destroy(g.Group)
	if err == r.ErrorNotFound {
		return ErrorGroupNotFound
	}
	return err
}

func NewUnlistService(r r.Repository[g.Group]) Service[GroupDto] {
	return &UnlistService{r}
}

// Group not found.
var ErrorGroupNotFound = errors.New("group not found")
