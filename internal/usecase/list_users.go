package usecase

import (
	"context"

	"github.com/silvioubaldino/ilia-users/internal/domain"
)

type listUsersRepository interface {
	List(ctx context.Context) ([]domain.User, error)
}

type ListUsers struct {
	repo listUsersRepository
}

func NewListUsers(repo listUsersRepository) *ListUsers {
	return &ListUsers{repo: repo}
}

func (uc *ListUsers) Execute(ctx context.Context) ([]domain.User, error) {
	users, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if users == nil {
		return []domain.User{}, nil
	}
	return users, nil
}
