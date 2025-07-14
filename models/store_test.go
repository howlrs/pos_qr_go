package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewStore は NewStore 関数のテストです
func TestNewStore(t *testing.T) {
	name := "テストストア"
	email := "test@example.com"
	password := "password123"
	address := "東京都渋谷区1-1-1"
	phone := "03-1234-5678"

	store := NewStore(name, email, password, address, phone)

	assert.NotNil(t, store, "NewStore should return a non-nil store")
	assert.NotEmpty(t, store.ID, "Store ID should not be empty")
	assert.True(t, len(store.ID) > len(StorePrefix), "Store ID should be longer than prefix")
	assert.Contains(t, store.ID, StorePrefix, "Store ID should contain the prefix")
	assert.Equal(t, name, store.Name, "Store name should match")
	assert.Equal(t, email, store.Email, "Store email should match")
	assert.Equal(t, password, store.Password, "Store password should match")
	assert.Equal(t, address, store.Address, "Store address should match")
	assert.Equal(t, phone, store.Phone, "Store phone should match")
	assert.False(t, store.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, store.UpdatedAt.IsZero(), "UpdatedAt should be set")
	assert.True(t, store.CreatedAt.Equal(store.UpdatedAt), "CreatedAt and UpdatedAt should be equal for new store")
}

// TestStore_ResetMetaFields は ResetMetaFields メソッドのテストです
func TestStore_ResetMetaFields(t *testing.T) {
	store := &Store{
		ID:        "old_id",
		Name:      "テストストア",
		Email:     "test@example.com",
		Password:  "password123",
		Address:   "東京都渋谷区1-1-1",
		Phone:     "03-1234-5678",
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	oldID := store.ID
	oldCreatedAt := store.CreatedAt
	oldUpdatedAt := store.UpdatedAt

	store.ResetMetaFields()

	assert.NotEqual(t, oldID, store.ID, "ID should be reset")
	assert.Contains(t, store.ID, StorePrefix, "New ID should contain the prefix")
	assert.NotEqual(t, oldCreatedAt, store.CreatedAt, "CreatedAt should be reset")
	assert.NotEqual(t, oldUpdatedAt, store.UpdatedAt, "UpdatedAt should be reset")
	assert.True(t, store.CreatedAt.Equal(store.UpdatedAt), "CreatedAt and UpdatedAt should be equal after reset")

	// ビジネスデータは変更されないことを確認
	assert.Equal(t, "テストストア", store.Name, "Name should remain unchanged")
	assert.Equal(t, "test@example.com", store.Email, "Email should remain unchanged")
}

// TestStore_FromValue は FromValue メソッドのテストです
func TestStore_FromValue(t *testing.T) {
	t.Run("有効なJSONデータ", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"name":     "JSONストア",
			"email":    "json@example.com",
			"password": "jsonpass123",
			"address":  "大阪府大阪市1-1-1",
			"phone":    "06-1234-5678",
		}

		jsonBytes, err := json.Marshal(jsonData)
		require.NoError(t, err, "JSON marshaling should not fail")

		store := &Store{}
		err = store.FromValue(jsonBytes)

		assert.NoError(t, err, "FromValue should not return an error")
		assert.NotEmpty(t, store.ID, "ID should be generated")
		assert.Equal(t, "JSONストア", store.Name, "Name should match JSON data")
		assert.Equal(t, "json@example.com", store.Email, "Email should match JSON data")
		assert.Equal(t, "jsonpass123", store.Password, "Password should match JSON data")
		assert.Equal(t, "大阪府大阪市1-1-1", store.Address, "Address should match JSON data")
		assert.Equal(t, "06-1234-5678", store.Phone, "Phone should match JSON data")
		assert.False(t, store.CreatedAt.IsZero(), "CreatedAt should be set")
		assert.False(t, store.UpdatedAt.IsZero(), "UpdatedAt should be set")
	})

	t.Run("部分的なJSONデータ", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"name":  "部分ストア",
			"email": "partial@example.com",
			// password, address, phone は省略
		}

		jsonBytes, err := json.Marshal(jsonData)
		require.NoError(t, err, "JSON marshaling should not fail")

		store := &Store{}
		err = store.FromValue(jsonBytes)

		assert.NoError(t, err, "FromValue should not return an error")
		assert.Equal(t, "部分ストア", store.Name, "Name should match JSON data")
		assert.Equal(t, "partial@example.com", store.Email, "Email should match JSON data")
		assert.Empty(t, store.Password, "Password should be empty")
		assert.Empty(t, store.Address, "Address should be empty")
		assert.Empty(t, store.Phone, "Phone should be empty")
	})

	t.Run("不正なJSON形式", func(t *testing.T) {
		invalidJSON := []byte(`{"name": "test", "email":}`)

		store := &Store{}
		err := store.FromValue(invalidJSON)

		assert.Error(t, err, "FromValue should return an error for invalid JSON")
	})

	t.Run("型が合わないJSON", func(t *testing.T) {
		jsonData := map[string]interface{}{
			"name":     123, // 数値（文字列ではない）
			"email":    "test@example.com",
			"password": true, // bool（文字列ではない）
			"address":  "住所",
			"phone":    []string{"123", "456"}, // 配列（文字列ではない）
		}

		jsonBytes, err := json.Marshal(jsonData)
		require.NoError(t, err, "JSON marshaling should not fail")

		store := &Store{}
		err = store.FromValue(jsonBytes)

		assert.NoError(t, err, "FromValue should not return an error")
		assert.Empty(t, store.Name, "Name should be empty for non-string value")
		assert.Equal(t, "test@example.com", store.Email, "Email should match JSON data")
		assert.Empty(t, store.Password, "Password should be empty for non-string value")
		assert.Equal(t, "住所", store.Address, "Address should match JSON data")
		assert.Empty(t, store.Phone, "Phone should be empty for non-string value")
	})
}

// TestStore_ValidateRequiredFields は ValidateRequiredFields メソッドのテストです
func TestStore_ValidateRequiredFields(t *testing.T) {
	t.Run("すべてのフィールドが有効", func(t *testing.T) {
		store := &Store{
			Name:     "有効ストア",
			Email:    "valid@example.com",
			Password: "validpass123",
			Address:  "有効住所",
			Phone:    "03-1234-5678",
		}

		err := store.ValidateRequiredFields()
		assert.NoError(t, err, "ValidateRequiredFields should not return an error for valid store")
	})

	t.Run("単一フィールドが空", func(t *testing.T) {
		store := &Store{
			Name:     "", // 空
			Email:    "test@example.com",
			Password: "password123",
			Address:  "住所",
			Phone:    "03-1234-5678",
		}

		err := store.ValidateRequiredFields()
		assert.Error(t, err, "ValidateRequiredFields should return an error for empty name")
	})

	t.Run("複数フィールドが空", func(t *testing.T) {
		store := &Store{
			Name:     "", // 空
			Email:    "test@example.com",
			Password: "", // 空
			Address:  "住所",
			Phone:    "", // 空
		}

		err := store.ValidateRequiredFields()
		assert.Error(t, err, "ValidateRequiredFields should return an error for multiple empty fields")
	})

	t.Run("空白文字のみのフィールド", func(t *testing.T) {
		store := &Store{
			Name:     "   ", // 空白のみ
			Email:    "test@example.com",
			Password: "\t\n", // タブと改行
			Address:  "住所",
			Phone:    "  \t  ", // 空白とタブ
		}

		err := store.ValidateRequiredFields()
		assert.Error(t, err, "ValidateRequiredFields should return an error for whitespace-only fields")
	})

	t.Run("すべてのフィールドが空", func(t *testing.T) {
		store := &Store{}

		err := store.ValidateRequiredFields()
		assert.Error(t, err, "ValidateRequiredFields should return an error for all empty fields")
	})
}

// TestStore_PasswordToHash は PasswordToHash メソッドのテストです
func TestStore_PasswordToHash(t *testing.T) {
	t.Run("有効なパスワード", func(t *testing.T) {
		store := &Store{
			Password: "plainpassword123",
		}

		originalPassword := store.Password
		err := store.PasswordToHash()

		assert.NoError(t, err, "PasswordToHash should not return an error")
		assert.NotEqual(t, originalPassword, store.Password, "Password should be hashed")
		assert.NotEmpty(t, store.Password, "Hashed password should not be empty")
		assert.True(t, len(store.Password) > len(originalPassword), "Hashed password should be longer than original")
	})

	t.Run("空のパスワード", func(t *testing.T) {
		store := &Store{
			Password: "",
		}

		// 7文字以下と73文字以上はエラーになる
		err := store.PasswordToHash()
		// 空のパスワードでもハッシュ化は可能（bcryptの動作による）
		assert.Error(t, err, "PasswordToHash should handle empty password")
	})
}

// TestStore_NormalizeTimestamps は NormalizeTimestamps メソッドのテストです
func TestStore_NormalizeTimestamps(t *testing.T) {
	t.Run("ゼロ値のタイムスタンプ", func(t *testing.T) {
		store := &Store{}

		store.NormalizeTimestamps()

		assert.False(t, store.CreatedAt.IsZero(), "CreatedAt should be set")
		assert.False(t, store.UpdatedAt.IsZero(), "UpdatedAt should be set")
		assert.True(t, time.Since(store.CreatedAt) < time.Second, "CreatedAt should be recent")
		assert.True(t, time.Since(store.UpdatedAt) < time.Second, "UpdatedAt should be recent")
	})

	t.Run("既存のタイムスタンプ", func(t *testing.T) {
		fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		store := &Store{
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}

		store.NormalizeTimestamps()

		assert.Equal(t, fixedTime, store.CreatedAt, "CreatedAt should remain unchanged")
		assert.Equal(t, fixedTime, store.UpdatedAt, "UpdatedAt should remain unchanged")
	})

	t.Run("部分的にゼロ値のタイムスタンプ", func(t *testing.T) {
		fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		store := &Store{
			CreatedAt: fixedTime,
			// UpdatedAt はゼロ値
		}

		store.NormalizeTimestamps()

		assert.Equal(t, fixedTime, store.CreatedAt, "CreatedAt should remain unchanged")
		assert.False(t, store.UpdatedAt.IsZero(), "UpdatedAt should be set")
		assert.True(t, time.Since(store.UpdatedAt) < time.Second, "UpdatedAt should be recent")
	})
}

// Test_getStringFromJSON は getStringFromJSON ヘルパー関数のテストです
func Test_getStringFromJSON(t *testing.T) {
	t.Run("有効な文字列値", func(t *testing.T) {
		data := map[string]interface{}{
			"key": "value",
		}

		result := getStringFromJSON(data, "key")
		assert.Equal(t, "value", result, "Should return the string value")
	})

	t.Run("存在しないキー", func(t *testing.T) {
		data := map[string]interface{}{
			"key": "value",
		}

		result := getStringFromJSON(data, "nonexistent")
		assert.Equal(t, "", result, "Should return empty string for nonexistent key")
	})

	t.Run("非文字列値", func(t *testing.T) {
		data := map[string]interface{}{
			"number": 123,
			"bool":   true,
			"array":  []string{"a", "b"},
			"object": map[string]string{"nested": "value"},
		}

		assert.Equal(t, "", getStringFromJSON(data, "number"), "Should return empty string for number")
		assert.Equal(t, "", getStringFromJSON(data, "bool"), "Should return empty string for bool")
		assert.Equal(t, "", getStringFromJSON(data, "array"), "Should return empty string for array")
		assert.Equal(t, "", getStringFromJSON(data, "object"), "Should return empty string for object")
	})

	t.Run("nil値", func(t *testing.T) {
		data := map[string]interface{}{
			"nil": nil,
		}

		result := getStringFromJSON(data, "nil")
		assert.Equal(t, "", result, "Should return empty string for nil value")
	})

	t.Run("空の文字列", func(t *testing.T) {
		data := map[string]interface{}{
			"empty": "",
		}

		result := getStringFromJSON(data, "empty")
		assert.Equal(t, "", result, "Should return empty string for empty string value")
	})
}

// TestStore_Integration は Store 構造体の統合テストです
func TestStore_Integration(t *testing.T) {
	t.Run("完全なライフサイクル", func(t *testing.T) {
		// 1. NewStore でストアを作成
		store := NewStore("統合テストストア", "integration@example.com", "password123", "統合住所", "03-9999-9999")

		// 2. 必須フィールドの検証
		err := store.ValidateRequiredFields()
		assert.NoError(t, err, "New store should pass validation")

		// 3. パスワードのハッシュ化
		originalPassword := store.Password
		err = store.PasswordToHash()
		assert.NoError(t, err, "Password hashing should succeed")
		assert.NotEqual(t, originalPassword, store.Password, "Password should be hashed")

		// 4. タイムスタンプの正規化（既に設定されているので変更されない）
		originalCreatedAt := store.CreatedAt
		originalUpdatedAt := store.UpdatedAt
		store.NormalizeTimestamps()
		assert.Equal(t, originalCreatedAt, store.CreatedAt, "CreatedAt should remain the same")
		assert.Equal(t, originalUpdatedAt, store.UpdatedAt, "UpdatedAt should remain the same")

		// 5. メタフィールドのリセット
		store.ResetMetaFields()
		assert.NotEqual(t, originalCreatedAt, store.CreatedAt, "CreatedAt should be reset")
		assert.NotEqual(t, originalUpdatedAt, store.UpdatedAt, "UpdatedAt should be reset")

		// 6. JSONからの復元
		jsonData := map[string]interface{}{
			"name":     "JSONから復元されたストア",
			"email":    "restored@example.com",
			"password": "restoredpass",
			"address":  "復元住所",
			"phone":    "03-8888-8888",
		}
		jsonBytes, _ := json.Marshal(jsonData)

		err = store.FromValue(jsonBytes)
		assert.NoError(t, err, "FromValue should succeed")
		assert.Equal(t, "JSONから復元されたストア", store.Name, "Name should be restored from JSON")
		assert.Equal(t, "restored@example.com", store.Email, "Email should be restored from JSON")
		assert.Equal(t, "restoredpass", store.Password, "Password should be restored from JSON")

		// 7. 最終検証
		err = store.ValidateRequiredFields()
		assert.NoError(t, err, "Restored store should pass validation")
	})
}
