package utils_test

import (
	"testing"

	"github.com/danisasmita/customer-search/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mysecurepassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecurepassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Test correct password
	assert.True(t, utils.CheckPasswordHash(password, hashedPassword))

	// Test incorrect password
	assert.False(t, utils.CheckPasswordHash("wrongpassword", hashedPassword))
}
