package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
)

type getUserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type GetUser struct {
	repo getUserRepository
}

func NewGetUser(repo getUserRepository) *GetUser {
	return &GetUser{repo: repo}
}

func (uc *GetUser) Execute(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return uc.repo.GetByID(ctx, id)
}
