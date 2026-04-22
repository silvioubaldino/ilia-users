package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/adapter/http/handler"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = domain.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}
		authOutput = usecase.AuthenticateOutput{
			User:        stored,
			AccessToken: "token-abc",
		}
	)

	type inputBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type mocks struct {
		ucInput  *usecase.AuthenticateInput
		ucOutput *usecase.AuthenticateOutput
		ucErr    error
		ucCalled bool
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		inputBody inputBody
		mocks     mocks
		expected  expected
	}{
		"should return 400 when body is missing required fields": {
			inputBody: inputBody{},
			mocks:     mocks{ucCalled: false},
			expected:  expected{statusCode: http.StatusBadRequest},
		},
		"should return 401 when credentials are invalid": {
			inputBody: inputBody{Email: "john@example.com", Password: "wrong"},
			mocks: mocks{
				ucInput:  &usecase.AuthenticateInput{Email: "john@example.com", Password: "wrong"},
				ucOutput: &usecase.AuthenticateOutput{},
				ucErr:    apperrors.ErrUnauthorized,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusUnauthorized},
		},
		"should return 500 when usecase fails": {
			inputBody: inputBody{Email: "john@example.com", Password: "secret"},
			mocks: mocks{
				ucInput:  &usecase.AuthenticateInput{Email: "john@example.com", Password: "secret"},
				ucOutput: &usecase.AuthenticateOutput{},
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 200 with token when credentials are valid": {
			inputBody: inputBody{Email: "john@example.com", Password: "secret"},
			mocks: mocks{
				ucInput:  &usecase.AuthenticateInput{Email: "john@example.com", Password: "secret"},
				ucOutput: &authOutput,
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusOK},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockAuthenticateUserUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", *tt.mocks.ucInput).Return(*tt.mocks.ucOutput, tt.mocks.ucErr)
			}

			h := handler.NewAuthHandler(uc)

			router := gin.New()
			router.POST("/auth", h.Login)

			body, _ := json.Marshal(tt.inputBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}
