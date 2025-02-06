package services

import (
	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
)

type RegistrationService struct {
	repo r.Repository[g.Group]
}

func (s *RegistrationService) Execute(g *g.Group) error {
	_, err := s.repo.Save(g)
	if err != nil {
		return err
	}

	return nil
}

func NewRegistrationService(r r.Repository[g.Group]) *RegistrationService {
	return &RegistrationService{r}
}
