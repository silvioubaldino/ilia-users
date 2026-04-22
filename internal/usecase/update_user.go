package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type updateUserRepository interface {
	Update(ctx context.Context, id uuid.UUID, updates domain.User) (domain.User, error)
}

type UpdateUser struct {
	repo updateUserRepository
}

func NewUpdateUser(repo updateUserRepository) *UpdateUser {
	return &UpdateUser{repo: repo}
}

type UpdateUserInput struct {
	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
}

func (uc *UpdateUser) Execute(ctx context.Context, id uuid.UUID, input UpdateUserInput) (domain.User, error) {
	updates := domain.User{}

	if input.FirstName != nil {
		updates.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		updates.LastName = *input.LastName
	}
	if input.Email != nil {
		updates.Email = *input.Email
	}
	if input.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 12)
		if err != nil {
			return domain.User{}, err
		}
		updates.Password = string(hashed)
	}

	return uc.repo.Update(ctx, id, updates)
}
