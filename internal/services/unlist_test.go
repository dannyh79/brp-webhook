package services_test

import (
	"errors"
	"testing"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	r "github.com/dannyh79/brp-webhook/internal/repositories"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/stretchr/testify/assert"
)

var _ r.Repository[g.Group] = (*mockUnlistRepo)(nil)

type mockUnlistRepo struct {
	shouldFail bool
	returnErr  error
}

func (r *mockUnlistRepo) Save(g *g.Group) (*g.Group, error) {
	return g, nil
}

func (r *mockUnlistRepo) Destroy(g *g.Group) error {
	if r.shouldFail {
		return r.returnErr
	}
	return nil
}

func Test_UnlistService(t *testing.T) {
	tcs := []struct {
		name              string
		dto               s.GroupDto
		expectRepoError   bool
		expectedRepoError error
		expectError       bool
		expectedError     error
	}{
		{
			name:        "Successfully unlists a group",
			dto:         s.GroupDto{Group: &g.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"}},
			expectError: false,
		},
		{
			name:              "Returns ErrorGroupNotFound when group is not found",
			dto:               s.GroupDto{Group: &g.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"}},
			expectRepoError:   true,
			expectedRepoError: r.ErrorNotFound,
			expectError:       true,
			expectedError:     s.ErrorGroupNotFound,
		},
		{
			name:              "Returns generic error when Destroy fails",
			dto:               s.GroupDto{Group: &g.Group{Id: "C12343d7945aa7d4a1f0ab43bc6cfa351"}},
			expectRepoError:   true,
			expectedRepoError: errors.New("repo error"),
			expectError:       true,
			expectedError:     errors.New("repo error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &mockUnlistRepo{shouldFail: tc.expectRepoError, returnErr: tc.expectedRepoError}
			service := s.NewUnlistService(repo)

			err := service.Execute(&tc.dto)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				assert.Equal(t, tc.expectedError, err, "Expected error %v but got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}
