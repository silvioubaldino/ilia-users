package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
)

type authenticateUserUseCase interface {
	Execute(ctx context.Context, input usecase.AuthenticateInput) (usecase.AuthenticateOutput, error)
}

type AuthHandler struct {
	authenticateUC authenticateUserUseCase
}

func NewAuthHandler(authenticateUC authenticateUserUseCase) *AuthHandler {
	return &AuthHandler{authenticateUC: authenticateUC}
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	User        usersResponse `json:"user"`
	AccessToken string        `json:"access_token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.authenticateUC.Execute(c.Request.Context(), usecase.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		User:        toUsersResponse(out.User),
		AccessToken: out.AccessToken,
	})
}
