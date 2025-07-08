package models

//
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// HashPasswordのテスト
func TestHashPassword(t *testing.T) {
	password := "passwird"
	hashed, err := HashPassword(password)
	assert.NoError(t, err, "HashPassword should not return an error")

	// ハッシュが空でないことを確認
	assert.NotEmpty(t, hashed, "Hashed password should not be empty")

	// ハッシュが正しいパスワードで検証できることを確認
	err = CheckPasswordHash(password, hashed)
	t.Log("Hashed password:", hashed, "Password:", password)
	assert.NoError(t, err, "CheckPasswordHash should not return an error for the correct password")
}
