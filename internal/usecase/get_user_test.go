package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_Execute(t *testing.T) {
	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = domain.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}
	)

	type input struct {
		id uuid.UUID
	}
	type mocks struct {
		repoGetByID    *domain.User
		repoGetByIDErr error
	}
	type expected struct {
		output domain.User
		err    error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return error when user is not found": {
			input: input{id: userID},
			mocks: mocks{
				repoGetByID:    &domain.User{},
				repoGetByIDErr: apperrors.ErrNotFound,
			},
			expected: expected{output: domain.User{}, err: apperrors.ErrNotFound},
		},
		"should return error when repo fails": {
			input: input{id: userID},
			mocks: mocks{
				repoGetByID:    &domain.User{},
				repoGetByIDErr: assert.AnError,
			},
			expected: expected{output: domain.User{}, err: assert.AnError},
		},
		"should return user when found": {
			input: input{id: userID},
			mocks: mocks{
				repoGetByID:    &stored,
				repoGetByIDErr: nil,
			},
			expected: expected{output: stored, err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockGetUserRepository{}
			repo.On("GetByID", tt.input.id).Return(*tt.mocks.repoGetByID, tt.mocks.repoGetByIDErr)

			uc := usecase.NewGetUser(repo)

			// Act
			got, err := uc.Execute(context.Background(), tt.input.id)

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
