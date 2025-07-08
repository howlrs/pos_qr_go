package usecases

import (
	"backend/models"
	"backend/repositories"
	"context"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

type Store struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) AsValidatedStore() (*models.Store, error) {
	store := &models.Store{
		Name:     s.Name,
		Email:    s.Email,
		Password: s.Password,
		Address:  s.Address,
		Phone:    s.Phone,
	}
	if err := store.ValidateRequiredFields(); err != nil {
		return nil, err // バリデーションエラーを返す
	}

	return store, nil
}

// Register は、Firestore データベースに新しい店舗を登録するメソッドです。
//
// このメソッドは以下の処理を行います:
//  1. Store 構造体の必須フィールドを検証します。
//  2. Store のメタ情報（ID, CreatedAt, UpdatedAt）をリセットします。
//  3. パスワードをハッシュ化します。
//  4. Firestore に店舗情報を保存します。
//
// 引数:
//   - db: Firestore クライアントインスタンス
//
// 戻り値:
//   - error: バリデーションや保存処理でエラーが発生した場合に返します。
func (s *Store) Register(ctx context.Context, db *firestore.Client) error {
	// APIで受ける不完全な店舗情報をチェック
	// ID, CreatedAt, UpdatedAtを付与してデータベース登録モデルとして整形する
	store, err := s.AsValidatedStore()
	if err != nil {
		return err
	}

	store.ResetMetaFields() // メタフィールドをリセット
	store.PasswordToHash()  // パスワードをハッシュ化

	log.Info().Msgf("Registering store: %+v", store)

	// add firestore logic here
	return repositories.NewStoreRepository(db).Create(ctx, store)
}

// GetAll は、Firestore データベースから店舗情報を取得するメソッドです。
func (s *Store) GetAll(ctx context.Context, db *firestore.Client) ([]*models.Store, error) {
	// Firestoreから店舗情報を取得
	return repositories.NewStoreRepository(db).Read(ctx)
}

// GetByID は、Firestore データベースから特定の店舗情報を取得するメソッドです。
func (s *Store) GetByID(ctx context.Context, db *firestore.Client, id string) (*models.Store, error) {
	// Store IDが空でないことを確認
	if id == "" {
		return nil, models.ErrStoreIDRequired
	}

	// Firestoreから特定の店舗情報を取得
	return repositories.NewStoreRepository(db).FindByID(ctx, id)
}

// Update は、Firestore データベースの店舗情報を更新するメソッドです。
func (s *Store) Update(ctx context.Context, db *firestore.Client, id string) error {
	// Store IDが空でないことを確認
	if id == "" {
		return models.ErrStoreIDRequired
	}

	// APIで受ける不完全な店舗情報をチェック
	store, err := s.AsValidatedStore()
	if err != nil {
		return err // バリデーションエラーを返す
	}

	store.ID = id           // 更新する店舗のIDを設定
	store.ResetMetaFields() // メタフィールドをリセット
	store.PasswordToHash()  // パスワードをハッシュ化

	// Firestoreの店舗情報を更新
	return repositories.NewStoreRepository(db).UpdateByID(ctx, store.ID, store)
}

// Delete は、Firestore データベースから店舗を削除するメソッドです。
func (s *Store) Delete(ctx context.Context, db *firestore.Client, id string) error {
	// Store IDが空でないことを確認
	if id == "" {
		return models.ErrStoreIDRequired
	}

	// Firestoreから店舗を削除
	return repositories.NewStoreRepository(db).DeleteByID(ctx, id)
}
