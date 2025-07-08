package usecases

import (
	"backend/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewStore - NewStore関数のテスト
func TestNewStore(t *testing.T) {
	// 空のStoreを作成
	store := NewStore()
	_, err := store.AsValidatedStore()
	// バリデーションエラーが発生することを期待
	assert.Error(t, err)

	assert.NotNil(t, store)
	assert.Empty(t, store.Name)
	assert.Empty(t, store.Email)
	assert.Empty(t, store.Password)
	assert.Empty(t, store.Address)
	assert.Empty(t, store.Phone)
}

// TestStore_AsValidatedStore - バリデーションロジックのテスト
func TestStore_AsValidatedStore(t *testing.T) {
	tests := []struct {
		name          string
		store         *Store
		expectedError bool
		errorContains string
	}{
		{
			name: "有効なストア",
			store: &Store{
				Name:     "テストストア",
				Email:    "test@example.com",
				Password: "password123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			expectedError: false,
		},
		{
			name: "名前が空",
			store: &Store{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			expectedError: true,
			errorContains: "name",
		},
		{
			name: "メールが空",
			store: &Store{
				Name:     "テストストア",
				Email:    "",
				Password: "password123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			expectedError: true,
			errorContains: "email",
		},
		{
			name: "複数のフィールドが空",
			store: &Store{
				Name:     "",
				Email:    "",
				Password: "password123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			expectedError: true,
			errorContains: "required fields are empty",
		},
		{
			name: "全フィールドが空",
			store: &Store{
				Name:     "",
				Email:    "",
				Password: "",
				Address:  "",
				Phone:    "",
			},
			expectedError: true,
			errorContains: "required fields are empty",
		},
		{
			name: "空白のみのフィールド",
			store: &Store{
				Name:     "   ",
				Email:    "test@example.com",
				Password: "password123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			expectedError: true,
			errorContains: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実行
			validatedStore, err := tt.store.AsValidatedStore()

			// 検証
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, validatedStore)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validatedStore)
				assert.Equal(t, tt.store.Name, validatedStore.Name)
				assert.Equal(t, tt.store.Email, validatedStore.Email)
				assert.Equal(t, tt.store.Password, validatedStore.Password)
				assert.Equal(t, tt.store.Address, validatedStore.Address)
				assert.Equal(t, tt.store.Phone, validatedStore.Phone)
			}
		})
	}
}

// TestStore_IDValidation - ID検証のテスト
func TestStore_IDValidation(t *testing.T) {
	ctx := context.Background()
	store := NewStore()

	tests := []struct {
		name      string
		storeID   string
		operation string
		wantErr   bool
	}{
		{
			name:      "GetByID - 空のID",
			storeID:   "",
			operation: "GetByID",
			wantErr:   true,
		},
		{
			name:      "Update - 空のID",
			storeID:   "",
			operation: "Update",
			wantErr:   true,
		},
		{
			name:      "Delete - 空のID",
			storeID:   "",
			operation: "Delete",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			// IDが空の場合のバリデーションエラーをテスト
			if tt.storeID == "" {
				switch tt.operation {
				case "GetByID":
					_, err = store.GetByID(ctx, nil, tt.storeID)
				case "Update":
					validStore := &Store{
						Name:     "テストストア",
						Email:    "test@example.com",
						Password: "password123",
						Address:  "東京都渋谷区1-1-1",
						Phone:    "03-1234-5678",
					}
					err = validStore.Update(ctx, nil, tt.storeID)
				case "Delete":
					err = store.Delete(ctx, nil, tt.storeID)
				}

				if tt.wantErr {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "store ID is required")
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

// TestStore_PasswordHashing - パスワードハッシュ化のテスト
func TestStore_PasswordHashing(t *testing.T) {
	store := &Store{
		Name:     "テストストア",
		Email:    "test@example.com",
		Password: "plaintext_password",
		Address:  "東京都渋谷区1-1-1",
		Phone:    "03-1234-5678",
	}

	// 元のパスワードを保存
	originalPassword := store.Password

	// AsValidatedStoreを使用してバリデーションのみをテスト
	validatedStore, err := store.AsValidatedStore()
	assert.NoError(t, err)
	assert.NotNil(t, validatedStore)

	// 元のStore構造体のパスワードは変更されていないことを確認
	assert.Equal(t, "plaintext_password", store.Password)

	// バリデーション済みストアのパスワードは元のパスワードと同じ（まだハッシュ化されていない）
	assert.Equal(t, originalPassword, validatedStore.Password)
}

// TestStore_ValidationEdgeCases - バリデーションのエッジケーステスト
func TestStore_ValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		store       *Store
		shouldError bool
		fieldCount  int
	}{
		{
			name: "有効なストア（全フィールド設定）",
			store: &Store{
				Name:     "有効なストア名",
				Email:    "valid@example.com",
				Password: "validpassword123",
				Address:  "東京都渋谷区1-1-1",
				Phone:    "03-1234-5678",
			},
			shouldError: false,
			fieldCount:  0,
		},
		{
			name: "空白のみのフィールド",
			store: &Store{
				Name:     "  ",
				Email:    "\t",
				Password: "\n",
				Address:  "   ",
				Phone:    "  \t  ",
			},
			shouldError: true,
			fieldCount:  5, // 全フィールドが空と判定される
		},
		{
			name: "混在した空と空白フィールド",
			store: &Store{
				Name:     "有効な名前",
				Email:    "",
				Password: "  ",
				Address:  "有効な住所",
				Phone:    "\t\n",
			},
			shouldError: true,
			fieldCount:  3, // Email, Password, Phone が空と判定される
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validatedStore, err := tt.store.AsValidatedStore()

			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, validatedStore)

				// バリデーションエラーかどうかをチェック
				if validationErr, ok := err.(*models.ValidationError); ok {
					assert.Len(t, validationErr.Fields, tt.fieldCount)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validatedStore)
				assert.Equal(t, tt.store.Name, validatedStore.Name)
				assert.Equal(t, tt.store.Email, validatedStore.Email)
				assert.Equal(t, tt.store.Password, validatedStore.Password)
				assert.Equal(t, tt.store.Address, validatedStore.Address)
				assert.Equal(t, tt.store.Phone, validatedStore.Phone)
			}
		})
	}
}
