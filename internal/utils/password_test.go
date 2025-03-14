package utils_test

import (
	"testing"

	"github.com/danisasmita/customer-search/internal/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	password := "securepassword"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		hashedPwd    string
		inputPwd     string
		expectsError bool
	}{
		{
			name:         "Valid Password",
			hashedPwd:    string(hashedPassword),
			inputPwd:     password,
			expectsError: false,
		},
		{
			name:         "Invalid Password",
			hashedPwd:    string(hashedPassword),
			inputPwd:     "wrongpassword",
			expectsError: true,
		},
		{
			name:         "Empty Password",
			hashedPwd:    string(hashedPassword),
			inputPwd:     "",
			expectsError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.CheckPassword(tt.hashedPwd, tt.inputPwd)
			if tt.expectsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
