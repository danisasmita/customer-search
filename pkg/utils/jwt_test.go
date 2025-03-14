package utils_test

import (
	"testing"
	"time"

	"github.com/danisasmita/customer-search/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := uint(123)
	token, err := utils.GenerateJWT(userID)
	assert.NoError(t, err, "GenerateJWT should not return an error")
	assert.NotEmpty(t, token, "Generated token should not be empty")
}

func TestValidateJWT(t *testing.T) {
	userID := uint(123)
	token, err := utils.GenerateJWT(userID)
	assert.NoError(t, err, "GenerateJWT should not return an error")
	assert.NotEmpty(t, token, "Generated token should not be empty")

	claims, err := utils.ValidateJWT(token)
	assert.NoError(t, err, "ValidateJWT should not return an error")
	assert.NotNil(t, claims, "Claims should not be nil")
	assert.Equal(t, userID, claims.UserID, "UserID should match")
}

func TestValidateJWTInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	claims, err := utils.ValidateJWT(invalidToken)
	assert.Error(t, err, "ValidateJWT should return an error for invalid token")
	assert.Nil(t, claims, "Claims should be nil for invalid token")
}

func TestValidateJWTExpiredToken(t *testing.T) {
	expiredClaims := &utils.Claims{
		UserID: 123,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(), // Expired token
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	signedToken, _ := token.SignedString([]byte("your-secret-key"))

	claims, err := utils.ValidateJWT(signedToken)
	assert.Error(t, err, "ValidateJWT should return an error for expired token")
	assert.Nil(t, claims, "Claims should be nil for expired token")
}
