package repositories

// storesコレクションへのアクセスを管理するリポジトリです。

import (
	"backend/models"
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type StoreRepository struct {
	client     *firestore.Client
	collection string
}

// NewStoreRepositoryは、StoreRepositoryの新しいインスタンスを生成します。
func NewStoreRepository(client *firestore.Client) Repository[models.Store] {
	if client == nil {
		return NewMockStoreRepository()
	}

	return &StoreRepository{
		client:     client,
		collection: "stores",
	}
}

type Store struct {
	ID        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	Email     string    `firestore:"email"`
	Password  string    `firestore:"password"`
	Address   string    `firestore:"address"`
	Phone     string    `firestore:"phone"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

func ToSetStore(store *models.Store) *Store {
	return &Store{
		ID:        store.ID,
		Name:      store.Name,
		Email:     store.Email,
		Password:  store.Password,
		Address:   store.Address,
		Phone:     store.Phone,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}

func (s *Store) ToUpdate() *Store {
	s.UpdatedAt = time.Now()
	return s
}

func (s *Store) ToModel() *models.Store {
	return &models.Store{
		ID:        s.ID,
		Name:      s.Name,
		Email:     s.Email,
		Password:  s.Password,
		Address:   s.Address,
		Phone:     s.Phone,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// Createは、新しい店舗ドキュメントをFirestoreに作成します。
func (r *StoreRepository) Create(ctx context.Context, store *models.Store) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(store.ID).Set(ctx, ToSetStore(store))
	return err
}

// Readは、すべての店舗ドキュメントをFirestoreから取得します。
func (r *StoreRepository) Read(ctx context.Context) ([]*models.Store, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return []*models.Store{}, err
	}

	if len(docs) == 0 {
		return []*models.Store{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	stores := make([]*models.Store, len(docs))
	for i, doc := range docs {
		store := &Store{}
		if err := doc.DataTo(store); err != nil {
			return []*models.Store{}, err
		}
		stores[i] = store.ToModel()
	}

	return stores, nil
}

// FindByIDは、指定されたIDを持つ店舗ドキュメントをFirestoreから検索します。
func (r *StoreRepository) FindByID(ctx context.Context, id string) (*models.Store, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	store := &Store{}
	if err := doc.DataTo(store); err != nil {
		return nil, err
	}

	return store.ToModel(), nil
}

// FindByFieldは、指定されたフィールドと値に一致する店舗ドキュメントをFirestoreから検索します。
func (r *StoreRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Store, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return []*models.Store{}, err
	}

	if len(docs) == 0 {
		return []*models.Store{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	stores := make([]*models.Store, len(docs))
	for i, doc := range docs {
		store := &Store{}
		if err := doc.DataTo(store); err != nil {
			return []*models.Store{}, err
		}
		stores[i] = store.ToModel()
	}

	return stores, nil
}

// UpdateByIDは、指定されたIDの店舗ドキュメントをFirestoreで更新します。
func (r *StoreRepository) UpdateByID(ctx context.Context, id string, store *models.Store) error {
	fields := []firestore.Update{
		{Path: "updated_at", Value: time.Now()},
	}

	// storeの値があるフィールドのみを更新
	updateFields := map[string]interface{}{
		"name":     store.Name,
		"email":    store.Email,
		"password": store.Password,
		"address":  store.Address,
		"phone":    store.Phone,
	}

	for path, value := range updateFields {
		if str, ok := value.(string); ok && str != "" {
			fields = append(fields, firestore.Update{Path: path, Value: value})
		}
	}

	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Update(ctx, fields)
	return err
}

// DeleteByIDは、指定されたIDの店舗ドキュメントをFirestoreから削除します。
func (r *StoreRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Delete(ctx)
	return err
}

// Countは、Firestore内の店舗ドキュメントの総数を返します。
func (r *StoreRepository) Count(ctx context.Context) (int, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

// Existsは、指定されたIDの店舗ドキュメントがFirestoreに存在するかどうかを確認します。
func (r *StoreRepository) Exists(ctx context.Context, id string) (bool, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return false, nil
	}
	return doc.Exists(), nil
}
