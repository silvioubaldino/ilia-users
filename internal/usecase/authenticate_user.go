package usecase

import (
	"context"

	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/silvioubaldino/ilia-users/pkg/jwtutil"
	"golang.org/x/crypto/bcrypt"
)

type authenticateUserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type AuthenticateUser struct {
	repo      authenticateUserRepository
	jwtSecret string
}

func NewAuthenticateUser(repo authenticateUserRepository, jwtSecret string) *AuthenticateUser {
	return &AuthenticateUser{repo: repo, jwtSecret: jwtSecret}
}

type AuthenticateInput struct {
	Email    string
	Password string
}

type AuthenticateOutput struct {
	User        domain.User
	AccessToken string
}

func (uc *AuthenticateUser) Execute(ctx context.Context, input AuthenticateInput) (AuthenticateOutput, error) {
	user, err := uc.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return AuthenticateOutput{}, apperrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return AuthenticateOutput{}, apperrors.ErrUnauthorized
	}

	token, err := jwtutil.GenerateToken(user.ID, user.Email, uc.jwtSecret)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	return AuthenticateOutput{User: user, AccessToken: token}, nil
}
