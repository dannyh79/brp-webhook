package services_test

import (
	"errors"
	"testing"

	"github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/services"
)

type mockRepo struct {
	shouldFail bool
}

func (r *mockRepo) Save(g *groups.Group) (*groups.Group, error) {
	if r.shouldFail {
		return nil, errors.New("failed to save group")
	}
	return g, nil
}

func Test_RegistrationService(t *testing.T) {
	tcs := []struct {
		name        string
		group       groups.Group
		expectError bool
	}{
		{
			name:        "Does not return error",
			group:       groups.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"},
			expectError: false,
		},
		{
			name:        "Returns error",
			group:       groups.Group{Id: "C56781862c40c77487fc60baf98fa7a6a"},
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := &mockRepo{shouldFail: tc.expectError}
			s := services.NewRegistrationService(r)
			err := s.Execute(&tc.group)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got one: %v", err)
				}
			}
		})
	}
}
