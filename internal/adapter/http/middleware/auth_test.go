package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/silvioubaldino/ilia-users/internal/adapter/http/middleware"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret"

func generateToken(userID, secret string, expiry time.Time) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		inputAuthHeader    string
		expectedStatusCode int
	}{
		"should return 401 when authorization header is missing": {
			inputAuthHeader:    "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		"should return 401 when token is expired": {
			inputAuthHeader:    "Bearer " + generateToken("user-1", testSecret, time.Now().Add(-1*time.Hour)),
			expectedStatusCode: http.StatusUnauthorized,
		},
		"should return 401 when token has invalid signature": {
			inputAuthHeader:    "Bearer " + generateToken("user-1", "wrong-secret", time.Now().Add(time.Hour)),
			expectedStatusCode: http.StatusUnauthorized,
		},
		"should return 200 when token is valid": {
			inputAuthHeader:    "Bearer " + generateToken("user-1", testSecret, time.Now().Add(time.Hour)),
			expectedStatusCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			var (
				w      = httptest.NewRecorder()
				req    = httptest.NewRequest(http.MethodGet, "/test", nil)
				router = gin.New()
			)

			router.Use(middleware.Auth(testSecret))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			if tt.inputAuthHeader != "" {
				req.Header.Set("Authorization", tt.inputAuthHeader)
			}

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
