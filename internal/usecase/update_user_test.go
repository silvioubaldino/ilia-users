package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(s string) *string { return &s }

func TestUpdateUser_Execute(t *testing.T) {
	var (
		userID  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		updated = domain.User{
			ID:        userID,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane@example.com",
		}
	)

	type input struct {
		id    uuid.UUID
		patch usecase.UpdateUserInput
	}
	type mocks struct {
		repoUpdate    *domain.User
		repoUpdateErr error
		repoUpdates   *domain.User
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
			input: input{id: userID, patch: usecase.UpdateUserInput{FirstName: strPtr("Jane")}},
			mocks: mocks{
				repoUpdate:    &domain.User{},
				repoUpdateErr: apperrors.ErrNotFound,
				repoUpdates:   &domain.User{FirstName: "Jane"},
			},
			expected: expected{output: domain.User{}, err: apperrors.ErrNotFound},
		},
		"should return error when repo fails": {
			input: input{id: userID, patch: usecase.UpdateUserInput{FirstName: strPtr("Jane")}},
			mocks: mocks{
				repoUpdate:    &domain.User{},
				repoUpdateErr: assert.AnError,
				repoUpdates:   &domain.User{FirstName: "Jane"},
			},
			expected: expected{output: domain.User{}, err: assert.AnError},
		},
		"should update user when input is valid": {
			input: input{id: userID, patch: usecase.UpdateUserInput{FirstName: strPtr("Jane")}},
			mocks: mocks{
				repoUpdate:    &updated,
				repoUpdateErr: nil,
				repoUpdates:   &domain.User{FirstName: "Jane"},
			},
			expected: expected{output: updated, err: nil},
		},
		"should hash password when password is provided": {
			input: input{id: userID, patch: usecase.UpdateUserInput{Password: strPtr("newpassword")}},
			mocks: mocks{
				repoUpdate:    &updated,
				repoUpdateErr: nil,
				repoUpdates:   nil, // matched via MatchedBy
			},
			expected: expected{output: updated, err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockUpdateUserRepository{}

			if tt.input.patch.Password != nil {
				repo.On("Update", tt.input.id, mock.MatchedBy(func(u domain.User) bool {
					return u.Password != *tt.input.patch.Password && u.Password != ""
				})).Return(*tt.mocks.repoUpdate, tt.mocks.repoUpdateErr)
			} else {
				repo.On("Update", tt.input.id, *tt.mocks.repoUpdates).
					Return(*tt.mocks.repoUpdate, tt.mocks.repoUpdateErr)
			}

			uc := usecase.NewUpdateUser(repo)

			// Act
			got, err := uc.Execute(context.Background(), tt.input.id, tt.input.patch)

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
