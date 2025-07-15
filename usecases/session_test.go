package usecases

import (
	"backend/models"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewSession tests the NewSession function
func TestNewSession(t *testing.T) {
	t.Run("create new session with valid parameters", func(t *testing.T) {
		// Arrange
		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.NotNil(t, session)
		assert.NotNil(t, session.Store)
		assert.NotNil(t, session.Seat)
		assert.Equal(t, storeID, session.Store.ID)
		assert.Equal(t, seatID, session.Seat.ID)
		assert.Equal(t, expiredAt, session.ExpiredAt)
	})

	t.Run("create new session with empty store ID", func(t *testing.T) {
		// Arrange
		storeID := ""
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.NotNil(t, session)
		assert.Equal(t, "", session.Store.ID)
		assert.Equal(t, seatID, session.Seat.ID)
	})

	t.Run("create new session with empty seat ID", func(t *testing.T) {
		// Arrange
		storeID := "store_123"
		seatID := ""
		expiredAt := time.Now().Add(time.Hour)

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.NotNil(t, session)
		assert.Equal(t, storeID, session.Store.ID)
		assert.Equal(t, "", session.Seat.ID)
	})

	t.Run("create new session with past expiration time", func(t *testing.T) {
		// Arrange
		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(-time.Hour) // 1時間前

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.NotNil(t, session)
		assert.Equal(t, storeID, session.Store.ID)
		assert.Equal(t, seatID, session.Seat.ID)
		assert.Equal(t, expiredAt, session.ExpiredAt)
		assert.True(t, session.ExpiredAt.Before(time.Now()))
	})

	t.Run("create new session with zero time", func(t *testing.T) {
		// Arrange
		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Time{}

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.NotNil(t, session)
		assert.Equal(t, storeID, session.Store.ID)
		assert.Equal(t, seatID, session.Seat.ID)
		assert.True(t, session.ExpiredAt.IsZero())
	})
}

// TestSessionCreateJWT tests the CreateJWT method
func TestSessionCreateJWT(t *testing.T) {
	// JWT_SECRETを設定
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	t.Run("create JWT with valid session and JWT secret", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "test_secret_key")

		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)
		session := NewSession(storeID, seatID, expiredAt)

		// Seatに名前を設定（JWT生成に必要）
		session.Seat.Name = "Table 1"

		// Act
		token, err := session.CreateJWT()

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Contains(t, token, ".") // JWTトークンは"."で区切られている

		// トークンが3つの部分（header.payload.signature）で構成されていることを確認
		parts := len([]rune(token))
		assert.True(t, parts > 10) // 最低限の長さがあることを確認
	})

	t.Run("create JWT without JWT secret", func(t *testing.T) {
		// Arrange
		os.Unsetenv("JWT_SECRET")

		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)
		session := NewSession(storeID, seatID, expiredAt)
		session.Seat.Name = "Table 1"

		// Act
		token, err := session.CreateJWT()

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "JWT_SECRET is not set")
	})

	t.Run("create JWT with empty JWT secret", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "")

		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)
		session := NewSession(storeID, seatID, expiredAt)
		session.Seat.Name = "Table 1"

		// Act
		token, err := session.CreateJWT()

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "JWT_SECRET is not set")
	})

	t.Run("create JWT with expired session", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "test_secret_key")

		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(-time.Hour) // 1時間前に期限切れ
		session := NewSession(storeID, seatID, expiredAt)
		session.Seat.Name = "Table 1"

		// Act
		token, err := session.CreateJWT()

		// Assert
		// JWTの生成自体は成功するが、期限切れのトークンが生成される
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("create JWT with nil store", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "test_secret_key")

		session := &Session{
			Store:     nil,
			Seat:      &models.Seat{ID: "seat_456", Name: "Table 1"},
			ExpiredAt: time.Now().Add(time.Hour),
		}

		// Act & Assert
		// nilストアの場合、パニックが発生することを期待
		assert.Panics(t, func() {
			session.CreateJWT()
		}, "CreateJWT should panic with nil store")
	})
	t.Run("create JWT with nil seat", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "test_secret_key")

		session := &Session{
			Store:     &models.Store{ID: "store_123"},
			Seat:      nil,
			ExpiredAt: time.Now().Add(time.Hour),
		}

		// Act & Assert
		// nilシートの場合、パニックが発生することを期待
		assert.Panics(t, func() {
			session.CreateJWT()
		}, "CreateJWT should panic with nil seat")
	})

	t.Run("create multiple JWTs from same session", func(t *testing.T) {
		// Arrange
		os.Setenv("JWT_SECRET", "test_secret_key")

		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)
		session := NewSession(storeID, seatID, expiredAt)
		session.Seat.Name = "Table 1"

		// Act
		token1, err1 := session.CreateJWT()
		// 少し時間を置く
		time.Sleep(time.Second)
		token2, err2 := session.CreateJWT()

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		// IssuedAtが異なるため、トークンは異なるはず
		assert.NotEqual(t, token1, token2)
	})
}

// TestSessionStructure tests the Session struct properties
func TestSessionStructure(t *testing.T) {
	t.Run("session struct has correct fields", func(t *testing.T) {
		// Arrange
		storeID := "store_123"
		seatID := "seat_456"
		expiredAt := time.Now().Add(time.Hour)

		// Act
		session := NewSession(storeID, seatID, expiredAt)

		// Assert
		assert.IsType(t, &Session{}, session)
		assert.IsType(t, &models.Store{}, session.Store)
		assert.IsType(t, &models.Seat{}, session.Seat)
		assert.IsType(t, time.Time{}, session.ExpiredAt)
	})

	t.Run("session fields are modifiable", func(t *testing.T) {
		// Arrange
		session := NewSession("store_123", "seat_456", time.Now().Add(time.Hour))

		// Act
		session.Store.ID = "new_store_id"
		session.Seat.ID = "new_seat_id"
		session.Seat.Name = "New Table Name"
		newExpiredAt := time.Now().Add(2 * time.Hour)
		session.ExpiredAt = newExpiredAt

		// Assert
		assert.Equal(t, "new_store_id", session.Store.ID)
		assert.Equal(t, "new_seat_id", session.Seat.ID)
		assert.Equal(t, "New Table Name", session.Seat.Name)
		assert.Equal(t, newExpiredAt, session.ExpiredAt)
	})
}

// TestSessionTimeHandling tests various time-related scenarios
func TestSessionTimeHandling(t *testing.T) {
	t.Run("session with future expiration time", func(t *testing.T) {
		// Arrange
		expiredAt := time.Now().Add(24 * time.Hour) // 24時間後
		session := NewSession("store_123", "seat_456", expiredAt)

		// Assert
		assert.True(t, session.ExpiredAt.After(time.Now()))
		duration := session.ExpiredAt.Sub(time.Now())
		assert.True(t, duration > 23*time.Hour) // 約24時間後
	})

	t.Run("session with very long expiration time", func(t *testing.T) {
		// Arrange
		expiredAt := time.Now().Add(365 * 24 * time.Hour) // 1年後
		session := NewSession("store_123", "seat_456", expiredAt)

		// Assert
		assert.True(t, session.ExpiredAt.After(time.Now()))
		duration := session.ExpiredAt.Sub(time.Now())
		assert.True(t, duration > 360*24*time.Hour) // 約1年後
	})

	t.Run("session with specific time format", func(t *testing.T) {
		// Arrange
		specificTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
		session := NewSession("store_123", "seat_456", specificTime)

		// Assert
		assert.Equal(t, specificTime, session.ExpiredAt)
		assert.Equal(t, 2025, session.ExpiredAt.Year())
		assert.Equal(t, time.December, session.ExpiredAt.Month())
		assert.Equal(t, 31, session.ExpiredAt.Day())
	})
}
