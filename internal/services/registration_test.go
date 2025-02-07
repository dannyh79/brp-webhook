package services_test

import (
	"errors"
	"testing"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	shouldFail bool
	error
}

func (r *mockRepo) Save(g *g.Group) (*g.Group, error) {
	if !r.shouldFail {
		return g, nil
	}
	if r.error == nil {
		r.error = errors.New("failed to save group")
	}
	return nil, r.error
}

func Test_RegistrationService(t *testing.T) {
	tcs := []struct {
		name              string
		dto               s.GroupDto
		expectRepoError   bool
		expectedRepoError error
		expectError       bool
		expectedError     error
	}{
		{
			name:        "Does not return error",
			dto:         s.GroupDto{Group: &g.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"}},
			expectError: false,
		},
		{
			name:              "Returns error when group already registered",
			dto:               s.GroupDto{Group: &g.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"}},
			expectRepoError:   true,
			expectedRepoError: r.ErrorAlreadyExists,
			expectError:       true,
			expectedError:     s.ErrorGroupAlreadyRegistered,
		},
		{
			name:            "Returns error",
			dto:             s.GroupDto{Group: &g.Group{Id: "C56781862c40c77487fc60baf98fa7a6a"}},
			expectRepoError: true,
			expectError:     true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := &mockRepo{shouldFail: tc.expectRepoError, error: tc.expectedRepoError}
			s := s.NewRegistrationService(r)
			err := s.Execute(&tc.dto)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err, "Expected error to be %v but got: %v", tc.expectedError, err)
				}
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
			}
		})
	}
}
