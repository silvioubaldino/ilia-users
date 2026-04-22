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

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = domain.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "hashed",
		}
	)

	type inputBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	type mocks struct {
		ucInput  *usecase.CreateUserInput
		ucOutput *domain.User
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
		"should return 409 when email already exists": {
			inputBody: inputBody{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
			mocks: mocks{
				ucInput:  &usecase.CreateUserInput{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
				ucOutput: &domain.User{},
				ucErr:    apperrors.ErrConflict,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusConflict},
		},
		"should return 500 when usecase fails": {
			inputBody: inputBody{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
			mocks: mocks{
				ucInput:  &usecase.CreateUserInput{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
				ucOutput: &domain.User{},
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 201 when user is created": {
			inputBody: inputBody{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
			mocks: mocks{
				ucInput:  &usecase.CreateUserInput{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"},
				ucOutput: &stored,
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusCreated},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockCreateUserUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", *tt.mocks.ucInput).Return(*tt.mocks.ucOutput, tt.mocks.ucErr)
			}

			h := handler.NewUserHandler(uc, &mockListUsersUseCase{}, &mockGetUserUseCase{}, &mockUpdateUserUseCase{}, &mockDeleteUserUseCase{})

			router := gin.New()
			router.POST("/users", h.Create)

			body, _ := json.Marshal(tt.inputBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		users  = []domain.User{
			{ID: userID, FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		}
	)

	type mocks struct {
		ucOutput []domain.User
		ucErr    error
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		mocks    mocks
		expected expected
	}{
		"should return 500 when usecase fails": {
			mocks:    mocks{ucOutput: nil, ucErr: assert.AnError},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 200 with users list": {
			mocks:    mocks{ucOutput: users, ucErr: nil},
			expected: expected{statusCode: http.StatusOK},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockListUsersUseCase{}
			uc.On("Execute").Return(tt.mocks.ucOutput, tt.mocks.ucErr)

			h := handler.NewUserHandler(&mockCreateUserUseCase{}, uc, &mockGetUserUseCase{}, &mockUpdateUserUseCase{}, &mockDeleteUserUseCase{})

			router := gin.New()
			router.GET("/users", h.List)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = domain.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	)

	type mocks struct {
		ucOutput *domain.User
		ucErr    error
		ucCalled bool
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		paramID  string
		mocks    mocks
		expected expected
	}{
		"should return 400 when id is invalid": {
			paramID:  "invalid",
			mocks:    mocks{ucCalled: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		"should return 404 when user is not found": {
			paramID: userID.String(),
			mocks: mocks{
				ucOutput: &domain.User{},
				ucErr:    apperrors.ErrNotFound,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusNotFound},
		},
		"should return 500 when usecase fails": {
			paramID: userID.String(),
			mocks: mocks{
				ucOutput: &domain.User{},
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 200 when user is found": {
			paramID: userID.String(),
			mocks: mocks{
				ucOutput: &stored,
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusOK},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockGetUserUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", userID).Return(*tt.mocks.ucOutput, tt.mocks.ucErr)
			}

			h := handler.NewUserHandler(&mockCreateUserUseCase{}, &mockListUsersUseCase{}, uc, &mockUpdateUserUseCase{}, &mockDeleteUserUseCase{})

			router := gin.New()
			router.GET("/users/:id", h.Get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.paramID, nil)

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID    = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		firstName = "Jane"
		stored    = domain.User{ID: userID, FirstName: "Jane", LastName: "Doe", Email: "john@example.com"}
	)

	type inputBody struct {
		FirstName *string `json:"first_name,omitempty"`
	}
	type mocks struct {
		ucInput  *usecase.UpdateUserInput
		ucOutput *domain.User
		ucErr    error
		ucCalled bool
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		paramID   string
		inputBody inputBody
		mocks     mocks
		expected  expected
	}{
		"should return 400 when id is invalid": {
			paramID:   "invalid",
			inputBody: inputBody{},
			mocks:     mocks{ucCalled: false},
			expected:  expected{statusCode: http.StatusBadRequest},
		},
		"should return 404 when user is not found": {
			paramID:   userID.String(),
			inputBody: inputBody{FirstName: &firstName},
			mocks: mocks{
				ucInput:  &usecase.UpdateUserInput{FirstName: &firstName},
				ucOutput: &domain.User{},
				ucErr:    apperrors.ErrNotFound,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusNotFound},
		},
		"should return 500 when usecase fails": {
			paramID:   userID.String(),
			inputBody: inputBody{FirstName: &firstName},
			mocks: mocks{
				ucInput:  &usecase.UpdateUserInput{FirstName: &firstName},
				ucOutput: &domain.User{},
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 200 when user is updated": {
			paramID:   userID.String(),
			inputBody: inputBody{FirstName: &firstName},
			mocks: mocks{
				ucInput:  &usecase.UpdateUserInput{FirstName: &firstName},
				ucOutput: &stored,
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusOK},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockUpdateUserUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", userID, *tt.mocks.ucInput).Return(*tt.mocks.ucOutput, tt.mocks.ucErr)
			}

			h := handler.NewUserHandler(&mockCreateUserUseCase{}, &mockListUsersUseCase{}, &mockGetUserUseCase{}, uc, &mockDeleteUserUseCase{})

			router := gin.New()
			router.PATCH("/users/:id", h.Update)

			body, _ := json.Marshal(tt.inputBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/users/"+tt.paramID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	)

	type mocks struct {
		ucErr    error
		ucCalled bool
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		paramID  string
		mocks    mocks
		expected expected
	}{
		"should return 400 when id is invalid": {
			paramID:  "invalid",
			mocks:    mocks{ucCalled: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		"should return 404 when user is not found": {
			paramID: userID.String(),
			mocks: mocks{
				ucErr:    apperrors.ErrNotFound,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusNotFound},
		},
		"should return 500 when usecase fails": {
			paramID: userID.String(),
			mocks: mocks{
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 204 when user is deleted": {
			paramID: userID.String(),
			mocks: mocks{
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusNoContent},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockDeleteUserUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", userID).Return(tt.mocks.ucErr)
			}

			h := handler.NewUserHandler(&mockCreateUserUseCase{}, &mockListUsersUseCase{}, &mockGetUserUseCase{}, &mockUpdateUserUseCase{}, uc)

			router := gin.New()
			router.DELETE("/users/:id", h.Delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.paramID, nil)

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}
