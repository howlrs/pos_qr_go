package models

import (
	"fmt"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword は、指定されたパスワードをbcryptでハッシュ化し、そのハッシュ値を返します。
//
// # 引数
//
// password: パスワードが8文字未満の場合はエラーを返します。
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters long")
	} else if len(password) > 72 {
		return "", fmt.Errorf("password must be at most 72 characters long")
	}

	// bcryptのコストパラメータを設定します。
	// bcrypt.DefaultCost (10) は一般的なウェブアプリケーションに適しています。
	cost := bcrypt.DefaultCost

	// パスワードをハッシュ化します。
	// bcrypt.GenerateFromPasswordはソルトを自動生成し、ハッシュに含めます。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateUniqueID(prefix string) string {
	// Sortrable ID generation using the xid package
	// 例: prefix_9m4e2mr0ui3e8a215n4g
	return prefix + xid.New().String()
}
