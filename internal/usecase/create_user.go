package usecase

import (
	"context"

	"github.com/silvioubaldino/ilia-users/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type createUserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type CreateUser struct {
	repo createUserRepository
}

func NewCreateUser(repo createUserRepository) *CreateUser {
	return &CreateUser{repo: repo}
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (uc *CreateUser) Execute(ctx context.Context, input CreateUserInput) (domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashed),
	}

	return uc.repo.Create(ctx, user)
}
