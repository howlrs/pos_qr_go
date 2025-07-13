package models

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name  string
	Email string
	Admin bool
	jwt.RegisteredClaims
}

func NewClaims(manager *Manager, isAdmin bool, exp time.Time) *Claims {
	return &Claims{
		Email: manager.Email,
		Admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
}

func (p *Claims) ToJwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

// SessionClaims はセッションのクレームを表す構造体です。
// 注文を行うセッションの有効期限と基礎情報を格納するためのものです。
type SessionClaims struct {
	StoreID string
	SeatID  string
	Name    string
	jwt.RegisteredClaims
}

func NewSessionClaims(store *Store, seat *Seat, exp time.Time) *SessionClaims {
	return &SessionClaims{
		StoreID: store.ID,
		SeatID:  seat.ID,
		Name:    seat.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
}

func (p *SessionClaims) ToJwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
