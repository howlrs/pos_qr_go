package models

import (
	"golang.org/x/crypto/bcrypt"
)

// Manager は管理者ユーザーを表す構造体です。
type Manager struct {
	Email    string
	Password string
}

func NewManager(email, password string) *Manager {
	return &Manager{
		Email:    email,
		Password: password,
	}
}

// パスワードのハッシュ化
func (u *Manager) ToEncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// パスワードの検証
func (u *Manager) IsVerifyPassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword))
}
