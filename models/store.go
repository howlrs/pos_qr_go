package models

import (
	"encoding/json"
	"strings"
	"time"
)

const StorePrefix = "store"

// Store はストアエンティティを表します
type Store struct {
	ID        string    `json:"id" db:"id" firestore:"id"`
	Name      string    `json:"name" db:"name" firestore:"name"`
	Email     string    `json:"email" db:"email" firestore:"email"`
	Password  string    `json:"password" db:"password" firestore:"password"`
	Address   string    `json:"address" db:"address" firestore:"address"`
	Phone     string    `json:"phone" db:"phone" firestore:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" firestore:"updated_at"`
}

// NewStore は新しいStoreインスタンスを作成します
func NewStore(name, email, password, address, phone string) *Store {
	uid := GenerateUniqueID(StorePrefix) // この関数が存在することを想定
	now := time.Now().UTC()

	return &Store{
		ID:        uid,
		Name:      name,
		Email:     email,
		Password:  password,
		Address:   address,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ResetMetaFields はID, CreatedAt, UpdatedAtを新たに設定します
func (s *Store) ResetMetaFields() {
	s.ID = GenerateUniqueID(StorePrefix) // この関数が存在することを想定
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
}

// FromValue はJSON値から新しいStoreインスタンスを作成します
func (s *Store) FromValue(data []byte) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	name := getStringFromJSON(jsonData, "name")
	email := getStringFromJSON(jsonData, "email")
	password := getStringFromJSON(jsonData, "password")
	address := getStringFromJSON(jsonData, "address")
	phone := getStringFromJSON(jsonData, "phone")

	newStore := NewStore(name, email, password, address, phone)
	*s = *newStore
	return nil
}

// ValidateRequiredFields は必須フィールドを検証します
func (s *Store) ValidateRequiredFields() error {
	emptyFields := []string{}

	fields := map[string]string{
		"name":     s.Name,
		"email":    s.Email,
		"password": s.Password,
		"address":  s.Address,
		"phone":    s.Phone,
	}

	for fieldName, fieldValue := range fields {
		if strings.TrimSpace(fieldValue) == "" {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	if len(emptyFields) > 0 {
		return NewValidationError(emptyFields...)
	}

	return nil
}

// PasswordToHash はパスワードをハッシュ化します
func (s *Store) PasswordToHash() error {
	password, err := HashPassword(s.Password) // この関数が存在することを想定
	if err != nil {
		return err
	}

	s.Password = password
	return nil
}

// NormalizeTimestamps はタイムスタンプを正規化します
func (s *Store) NormalizeTimestamps() {
	now := time.Now().UTC()

	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = now
	}
}

// getStringFromJSON はJSONマップから文字列を取得するヘルパー関数です
func getStringFromJSON(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
