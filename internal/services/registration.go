package services

import (
	"errors"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
)

var _ Service[g.Group] = (*RegistrationService)(nil)

type RegistrationService struct {
	repo r.Repository[g.Group]
}

func (s *RegistrationService) Execute(g *g.Group) error {
	_, err := s.repo.Save(g)
	if err == r.ErrorAlreadyExists {
		return ErrorGroupAlreadyRegistered
	}
	return err
}

func NewRegistrationService(r r.Repository[g.Group]) Service[g.Group] {
	return &RegistrationService{r}
}

// Group already registered.
var ErrorGroupAlreadyRegistered = errors.New("Group already registered")
