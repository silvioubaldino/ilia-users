package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_Execute(t *testing.T) {
	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	)

	type input struct {
		id uuid.UUID
	}
	type mocks struct {
		repoDeleteErr error
	}
	type expected struct {
		err error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return error when user is not found": {
			input:    input{id: userID},
			mocks:    mocks{repoDeleteErr: apperrors.ErrNotFound},
			expected: expected{err: apperrors.ErrNotFound},
		},
		"should return error when repo fails": {
			input:    input{id: userID},
			mocks:    mocks{repoDeleteErr: assert.AnError},
			expected: expected{err: assert.AnError},
		},
		"should delete user when found": {
			input:    input{id: userID},
			mocks:    mocks{repoDeleteErr: nil},
			expected: expected{err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockDeleteUserRepository{}
			repo.On("Delete", tt.input.id).Return(tt.mocks.repoDeleteErr)

			uc := usecase.NewDeleteUser(repo)

			// Act
			err := uc.Execute(context.Background(), tt.input.id)

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			repo.AssertExpectations(t)
		})
	}
}
