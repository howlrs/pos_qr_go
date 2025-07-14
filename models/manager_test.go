package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestNewManager はシンプルで明確なため、変更ありません。
func TestNewManager(t *testing.T) {
	email := "test@example.com"
	password := "password123"

	manager := NewManager(email, password)

	assert.NotNil(t, manager)
	assert.Equal(t, email, manager.Email)
	assert.Equal(t, password, manager.Password)
}

// TestManager_ToEncryptPassword はテーブル駆動テストにリファクタリングしました。
// 様々な種類のパスワードに対する暗号化処理をまとめてテストします。
func TestManager_ToEncryptPassword(t *testing.T) {
	// bcrypt.GenerateFromPasswordは72バイトを超えるパスワードでエラーを返します。
	longPassword := "this_is_a_very_long_password_that_is_definitely_longer_than_72_bytes_which_is_the_limit_for_bcrypt"

	testCases := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "正常なパスワード",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "空のパスワード",
			password:    "", // 空文字列でもbcryptはハッシュを生成します
			expectError: false,
		},
		{
			name:        "特殊文字を含むパスワード",
			password:    "p@ssw0rd!@#$%^&*()",
			expectError: false,
		},
		{
			name:        "Unicode文字を含むパスワード",
			password:    "パスワード🔒",
			expectError: false,
		},
		{
			name:        "非常に短いパスワード",
			password:    "a",
			expectError: false,
		},
		{
			name:        "長すぎるパスワード (bcryptの上限超え)",
			password:    longPassword,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewManager("test@example.com", tc.password)

			err := manager.ToEncryptPassword()

			if tc.expectError {
				assert.Error(t, err, "エラーが期待されるケースです")
				// エラー発生時は、元のパスワードが変更されていないことを確認します
				assert.Equal(t, tc.password, manager.Password, "エラー発生時、パスワードは変更されないべきです")
				return
			}

			// エラーが期待されないケース
			require.NoError(t, err, "エラーは期待されないケースです")
			assert.NotEmpty(t, manager.Password, "暗号化されたパスワードは空であってはなりません")
			assert.NotEqual(t, tc.password, manager.Password, "パスワードは暗号化されている必要があります")

			// 暗号化されたパスワードが、元のパスワードで検証できることを確認します
			err = manager.IsVerifyPassword(tc.password)
			assert.NoError(t, err, "暗号化されたパスワードは元のパスワードで検証できる必要があります")
		})
	}
}

// TestManager_IsVerifyPassword もテーブル駆動テストにリファクタリングしました。
// 検証ロジックの様々なパターンに焦点を当ててテストします。
func TestManager_IsVerifyPassword(t *testing.T) {
	originalPassword := "password123"
	wrongPassword := "wrongpassword"

	// 事前に暗号化したManagerを準備
	encryptedManager := NewManager("test@example.com", originalPassword)
	err := encryptedManager.ToEncryptPassword()
	require.NoError(t, err, "テストセットアップ: パスワードの暗号化に失敗しました")

	// 暗号化していない生のパスワードを持つManagerを準備
	plainManager := NewManager("test@example.com", originalPassword)

	testCases := []struct {
		name             string
		manager          *Manager // テスト対象のManagerインスタンス
		passwordToVerify string
		expectedErr      error // 期待される特定のエラー（nilの場合はエラーなし）
	}{
		{
			name:             "成功: 正しいパスワードで検証",
			manager:          encryptedManager,
			passwordToVerify: originalPassword,
			expectedErr:      nil,
		},
		{
			name:             "失敗: 間違ったパスワードで検証",
			manager:          encryptedManager,
			passwordToVerify: wrongPassword,
			expectedErr:      bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:             "失敗: 空のパスワードで検証",
			manager:          encryptedManager,
			passwordToVerify: "",
			expectedErr:      bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:             "失敗: 暗号化されていないパスワードで検証を試みる",
			manager:          plainManager,
			passwordToVerify: originalPassword,
			// bcrypt.CompareHashAndPasswordは不正なハッシュに対してこのエラーを返します
			expectedErr: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.IsVerifyPassword(tc.passwordToVerify)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr, "期待されるエラーと実際のエラーが一致すべきです")
			} else {
				assert.NoError(t, err, "エラーは期待されません")
			}
		})
	}
}

// TestManager_PasswordWorkflow は、一連の正常な利用シナリオをテストするもので、非常に有用です。
// 構造を少し整理し、コメントを加えて意図を明確にしました。
func TestManager_PasswordWorkflow(t *testing.T) {
	// このテストは、Managerの作成からパスワードの暗号化、検証までの一連の正常な流れを確認します。
	email := "admin@example.com"
	password := "securePassword123!"

	// 1. Managerを作成
	manager := NewManager(email, password)
	require.NotNil(t, manager)
	assert.Equal(t, email, manager.Email)
	assert.Equal(t, password, manager.Password)

	// 2. パスワードを暗号化
	originalPassword := manager.Password
	err := manager.ToEncryptPassword()
	require.NoError(t, err)
	assert.NotEqual(t, originalPassword, manager.Password, "パスワードが暗号化されていること")

	// 3. 正しいパスワードで検証
	err = manager.IsVerifyPassword(originalPassword)
	assert.NoError(t, err, "正しいパスワードで検証が成功すること")

	// 4. 間違ったパスワードで検証
	err = manager.IsVerifyPassword("wrongPassword")
	assert.Error(t, err, "間違ったパスワードで検証が失敗すること")
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}

// TestManager_ToEncryptPassword_DoubleEncryption は、複数回暗号化した場合の挙動をテストします。
// 元のコードの「複数回の暗号化」テストを、意図が明確になるように独立させました。
func TestManager_ToEncryptPassword_DoubleEncryption(t *testing.T) {
	// このテストは、ToEncryptPasswordが冪等ではなく、呼び出されるたびに新しいハッシュを生成し、
	// 結果として元のパスワードでの検証が不可能になる（二重暗号化）ことを確認します。
	manager := NewManager("test@example.com", "password123")
	originalPassword := "password123"

	// 1回目の暗号化
	err := manager.ToEncryptPassword()
	require.NoError(t, err)
	firstHash := manager.Password

	// この時点では元のパスワードで検証できるはずです
	require.NoError(t, manager.IsVerifyPassword(originalPassword))

	// 2回目の暗号化（既にハッシュ化された文字列を再度暗号化）
	err = manager.ToEncryptPassword()
	require.NoError(t, err)
	secondHash := manager.Password

	// bcryptのソルト機能により、ハッシュは毎回異なるものになります
	assert.NotEqual(t, firstHash, secondHash, "再暗号化によりハッシュは変更されるべきです")

	// 二重に暗号化されたため、元のパスワードでの検証は失敗します
	err = manager.IsVerifyPassword(originalPassword)
	assert.Error(t, err, "二重に暗号化されたパスワードは元のパスワードでは検証できません")
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
