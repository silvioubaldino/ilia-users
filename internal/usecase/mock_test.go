package usecase_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/stretchr/testify/mock"
)

type mockCreateUserRepository struct {
	mock.Mock
}

func (m *mockCreateUserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockGetUserRepository struct {
	mock.Mock
}

func (m *mockGetUserRepository) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockListUsersRepository struct {
	mock.Mock
}

func (m *mockListUsersRepository) List(_ context.Context) ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

type mockUpdateUserRepository struct {
	mock.Mock
}

func (m *mockUpdateUserRepository) Update(_ context.Context, id uuid.UUID, updates domain.User) (domain.User, error) {
	args := m.Called(id, updates)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockDeleteUserRepository struct {
	mock.Mock
}

func (m *mockDeleteUserRepository) Delete(_ context.Context, id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockAuthenticateUserRepository struct {
	mock.Mock
}

func (m *mockAuthenticateUserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(domain.User), args.Error(1)
}
