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

func TestCreateUser_Execute(t *testing.T) {
	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = domain.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "$2a$12$hashedpassword",
		}
	)

	type input struct {
		firstName string
		lastName  string
		email     string
		password  string
	}
	type mocks struct {
		repoCreate       *domain.User
		repoCreateErr    error
		repoCreateCalled bool
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
		"should return error when repo fails": {
			input: input{firstName: "John", lastName: "Doe", email: "john@example.com", password: "secret"},
			mocks: mocks{
				repoCreate:       &domain.User{},
				repoCreateErr:    assert.AnError,
				repoCreateCalled: true,
			},
			expected: expected{output: domain.User{}, err: assert.AnError},
		},
		"should return conflict when email already exists": {
			input: input{firstName: "John", lastName: "Doe", email: "john@example.com", password: "secret"},
			mocks: mocks{
				repoCreate:       &domain.User{},
				repoCreateErr:    apperrors.ErrConflict,
				repoCreateCalled: true,
			},
			expected: expected{output: domain.User{}, err: apperrors.ErrConflict},
		},
		"should create user when input is valid": {
			input: input{firstName: "John", lastName: "Doe", email: "john@example.com", password: "secret"},
			mocks: mocks{
				repoCreate:       &stored,
				repoCreateErr:    nil,
				repoCreateCalled: true,
			},
			expected: expected{output: stored, err: nil},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockCreateUserRepository{}
			if tt.mocks.repoCreateCalled {
				repo.On("Create", mock.MatchedBy(func(u domain.User) bool {
					return u.FirstName == tt.input.firstName &&
						u.LastName == tt.input.lastName &&
						u.Email == tt.input.email &&
						u.Password != tt.input.password &&
						u.Password != ""
				})).Return(*tt.mocks.repoCreate, tt.mocks.repoCreateErr)
			}

			uc := usecase.NewCreateUser(repo)

			// Act
			got, err := uc.Execute(context.Background(), usecase.CreateUserInput{
				FirstName: tt.input.firstName,
				LastName:  tt.input.lastName,
				Email:     tt.input.email,
				Password:  tt.input.password,
			})

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
