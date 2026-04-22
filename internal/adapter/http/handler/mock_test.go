package handler_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type mockCreateUserUseCase struct {
	mock.Mock
}

func (m *mockCreateUserUseCase) Execute(_ context.Context, input usecase.CreateUserInput) (domain.User, error) {
	args := m.Called(input)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockListUsersUseCase struct {
	mock.Mock
}

func (m *mockListUsersUseCase) Execute(_ context.Context) ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

type mockGetUserUseCase struct {
	mock.Mock
}

func (m *mockGetUserUseCase) Execute(_ context.Context, id uuid.UUID) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockUpdateUserUseCase struct {
	mock.Mock
}

func (m *mockUpdateUserUseCase) Execute(_ context.Context, id uuid.UUID, input usecase.UpdateUserInput) (domain.User, error) {
	args := m.Called(id, input)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockDeleteUserUseCase struct {
	mock.Mock
}

func (m *mockDeleteUserUseCase) Execute(_ context.Context, id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockAuthenticateUserUseCase struct {
	mock.Mock
}

func (m *mockAuthenticateUserUseCase) Execute(_ context.Context, input usecase.AuthenticateInput) (usecase.AuthenticateOutput, error) {
	args := m.Called(input)
	return args.Get(0).(usecase.AuthenticateOutput), args.Error(1)
}
