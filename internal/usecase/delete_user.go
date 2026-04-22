package usecase

import (
	"context"

	"github.com/google/uuid"
)

type deleteUserRepository interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleteUser struct {
	repo deleteUserRepository
}

func NewDeleteUser(repo deleteUserRepository) *DeleteUser {
	return &DeleteUser{repo: repo}
}

func (uc *DeleteUser) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
