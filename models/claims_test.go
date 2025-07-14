package models

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のJWTシークレット
const testJWTSecret = "test_secret_key"

func TestMain(m *testing.M) {
	// テスト実行前にJWT_SECRETを設定
	os.Setenv("JWT_SECRET", testJWTSecret)
	code := m.Run()
	// テスト実行後にクリーンアップ
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func TestNewClaims(t *testing.T) {
	manager := &Manager{
		Email: "test@example.com",
	}
	isAdmin := true
	exp := time.Now().Add(time.Hour)

	claims := NewClaims(manager, isAdmin, exp)

	assert.Equal(t, manager.Email, claims.Email)
	assert.Equal(t, isAdmin, claims.Admin)
	assert.Equal(t, exp.Unix(), claims.ExpiresAt.Unix())
}

func TestClaims_ToJwtToken(t *testing.T) {
	t.Run("正常なケース", func(t *testing.T) {
		manager := &Manager{
			Email: "test@example.com",
		}
		claims := NewClaims(manager, true, time.Now().Add(time.Hour))

		token, err := claims.ToJwtToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// トークンをパースして検証
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(testJWTSecret), nil
		})

		require.NoError(t, err)
		require.True(t, parsedToken.Valid)

		parsedClaims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)
		assert.Equal(t, claims.Email, parsedClaims.Email)
		assert.Equal(t, claims.Admin, parsedClaims.Admin)
	})

	t.Run("JWT_SECRETが設定されていない場合", func(t *testing.T) {
		// 一時的にJWT_SECRETを削除
		originalSecret := os.Getenv("JWT_SECRET")
		os.Unsetenv("JWT_SECRET")
		defer os.Setenv("JWT_SECRET", originalSecret)

		manager := &Manager{
			Email: "test@example.com",
		}
		claims := NewClaims(manager, true, time.Now().Add(time.Hour))

		token, err := claims.ToJwtToken()

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "JWT_SECRET is not set")
	})
}

func TestNewSessionClaims(t *testing.T) {
	store := &Store{
		ID:   "store_123",
		Name: "Test Store",
	}
	seat := &Seat{
		ID:   "seat_456",
		Name: "Table 1",
	}
	exp := time.Now().Add(time.Hour * 2)

	sessionClaims := NewSessionClaims(store, seat, exp)

	assert.Equal(t, store.ID, sessionClaims.StoreID)
	assert.Equal(t, seat.ID, sessionClaims.SeatID)
	assert.Equal(t, seat.Name, sessionClaims.Name)
	assert.Equal(t, exp.Unix(), sessionClaims.ExpiresAt.Unix())
	assert.NotNil(t, sessionClaims.IssuedAt)
}

func TestSessionClaims_ToJwtToken(t *testing.T) {
	t.Run("正常なケース", func(t *testing.T) {
		store := &Store{
			ID:   "store_123",
			Name: "Test Store",
		}
		seat := &Seat{
			ID:   "seat_456",
			Name: "Table 1",
		}
		sessionClaims := NewSessionClaims(store, seat, time.Now().Add(time.Hour*2))

		token, err := sessionClaims.ToJwtToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// トークンをパースして検証
		parsedToken, err := jwt.ParseWithClaims(token, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(testJWTSecret), nil
		})

		require.NoError(t, err)
		require.True(t, parsedToken.Valid)

		parsedClaims, ok := parsedToken.Claims.(*SessionClaims)
		require.True(t, ok)
		assert.Equal(t, sessionClaims.StoreID, parsedClaims.StoreID)
		assert.Equal(t, sessionClaims.SeatID, parsedClaims.SeatID)
		assert.Equal(t, sessionClaims.Name, parsedClaims.Name)
	})

	t.Run("JWT_SECRETが設定されていない場合", func(t *testing.T) {
		// 一時的にJWT_SECRETを削除
		originalSecret := os.Getenv("JWT_SECRET")
		os.Unsetenv("JWT_SECRET")
		defer os.Setenv("JWT_SECRET", originalSecret)

		store := &Store{
			ID:   "store_123",
			Name: "Test Store",
		}
		seat := &Seat{
			ID:   "seat_456",
			Name: "Table 1",
		}
		sessionClaims := NewSessionClaims(store, seat, time.Now().Add(time.Hour*2))

		token, err := sessionClaims.ToJwtToken()

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "JWT_SECRET is not set")
	})
}

func TestClaimsTokenExpiration(t *testing.T) {
	t.Run("期限切れのトークン", func(t *testing.T) {
		manager := &Manager{
			Email: "test@example.com",
		}
		// 過去の時刻を設定
		pastTime := time.Now().Add(-time.Hour)
		claims := NewClaims(manager, true, pastTime)

		token, err := claims.ToJwtToken()
		require.NoError(t, err)

		// 期限切れトークンのパース
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(testJWTSecret), nil
		})

		// パースはできるが、有効性チェックで false になる
		assert.Error(t, err)
		assert.False(t, parsedToken.Valid)
	})
}

func TestSessionClaimsTokenExpiration(t *testing.T) {
	t.Run("期限切れのセッショントークン", func(t *testing.T) {
		store := &Store{
			ID:   "store_123",
			Name: "Test Store",
		}
		seat := &Seat{
			ID:   "seat_456",
			Name: "Table 1",
		}
		// 過去の時刻を設定
		pastTime := time.Now().Add(-time.Hour)
		sessionClaims := NewSessionClaims(store, seat, pastTime)

		token, err := sessionClaims.ToJwtToken()
		require.NoError(t, err)

		// 期限切れトークンのパース
		parsedToken, err := jwt.ParseWithClaims(token, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(testJWTSecret), nil
		})

		// パースはできるが、有効性チェックで false になる
		assert.Error(t, err)
		assert.False(t, parsedToken.Valid)
	})
}
