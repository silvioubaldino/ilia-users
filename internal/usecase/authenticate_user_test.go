package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticateUser_Execute(t *testing.T) {
	correctPassword := "correct-password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.MinCost)

	var (
		userID    = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		jwtSecret = "test-secret"
		stored    = domain.User{
			ID:       userID,
			Email:    "john@example.com",
			Password: string(hashedPassword),
		}
	)

	type input struct {
		email    string
		password string
	}
	type mocks struct {
		repoGetByEmail    *domain.User
		repoGetByEmailErr error
	}
	type expected struct {
		hasToken bool
		err      error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return unauthorized when user is not found": {
			input: input{email: "john@example.com", password: correctPassword},
			mocks: mocks{
				repoGetByEmail:    &domain.User{},
				repoGetByEmailErr: apperrors.ErrNotFound,
			},
			expected: expected{hasToken: false, err: apperrors.ErrUnauthorized},
		},
		"should return unauthorized when password is wrong": {
			input: input{email: "john@example.com", password: "wrong-password"},
			mocks: mocks{
				repoGetByEmail:    &stored,
				repoGetByEmailErr: nil,
			},
			expected: expected{hasToken: false, err: apperrors.ErrUnauthorized},
		},
		"should return token when credentials are valid": {
			input: input{email: "john@example.com", password: correctPassword},
			mocks: mocks{
				repoGetByEmail:    &stored,
				repoGetByEmailErr: nil,
			},
			expected: expected{hasToken: true, err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockAuthenticateUserRepository{}
			repo.On("GetByEmail", tt.input.email).Return(*tt.mocks.repoGetByEmail, tt.mocks.repoGetByEmailErr)

			uc := usecase.NewAuthenticateUser(repo, jwtSecret)

			// Act
			got, err := uc.Execute(context.Background(), usecase.AuthenticateInput{
				Email:    tt.input.email,
				Password: tt.input.password,
			})

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			if tt.expected.hasToken {
				assert.NotEmpty(t, got.AccessToken)
			} else {
				assert.Empty(t, got.AccessToken)
			}
			repo.AssertExpectations(t)
		})
	}
}
