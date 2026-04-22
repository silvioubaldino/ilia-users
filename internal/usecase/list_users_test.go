package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestListUsers_Execute(t *testing.T) {
	var (
		users = []domain.User{
			{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
		}
	)

	type mocks struct {
		repoList    []domain.User
		repoListErr error
	}
	type expected struct {
		output []domain.User
		err    error
	}

	tests := map[string]struct {
		mocks    mocks
		expected expected
	}{
		"should return error when repo fails": {
			mocks:    mocks{repoList: nil, repoListErr: assert.AnError},
			expected: expected{output: nil, err: assert.AnError},
		},
		"should return empty slice when no users exist": {
			mocks:    mocks{repoList: []domain.User{}, repoListErr: nil},
			expected: expected{output: []domain.User{}, err: nil},
		},
		"should return users when they exist": {
			mocks:    mocks{repoList: users, repoListErr: nil},
			expected: expected{output: users, err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockListUsersRepository{}
			repo.On("List").Return(tt.mocks.repoList, tt.mocks.repoListErr)

			uc := usecase.NewListUsers(repo)

			// Act
			got, err := uc.Execute(context.Background())

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
