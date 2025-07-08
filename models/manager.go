package models

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	COLLECTION_USER = "users"
)

// Manager は管理者ユーザーを表す構造体です。
type Manager struct {
	ID       string `json:"id" db:"id" firestore:"id"`
	Email    string `json:"email" db:"email" firestore:"email"`
	Password string `json:"password" db:"password" firestore:"password"`
}

// ToCollection はユーザーマネージャーのコレクション名を返します。
func (u *Manager) ToCollection(isTest bool) string {
	if isTest {
		return "test_" + COLLECTION_USER
	}
	return COLLECTION_USER
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
