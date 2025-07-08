package repositories

// storesコレクションへのアクセスを管理するリポジトリです。

import (
	"backend/models"
	"context"

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

// Createは、新しい店舗ドキュメントをFirestoreに作成します。
func (r *StoreRepository) Create(ctx context.Context, store *models.Store) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(store.ID).Set(ctx, store)
	return err
}

// Readは、すべての店舗ドキュメントをFirestoreから取得します。
func (r *StoreRepository) Read(ctx context.Context) ([]*models.Store, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var stores []*models.Store
	for _, doc := range docs {
		var store models.Store
		if err := doc.DataTo(&store); err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}
	return stores, nil
}

// FindByIDは、指定されたIDを持つ店舗ドキュメントをFirestoreから検索します。
func (r *StoreRepository) FindByID(ctx context.Context, id string) (*models.Store, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var store models.Store
	if err := doc.DataTo(&store); err != nil {
		return nil, err
	}
	return &store, nil
}

// FindByFieldは、指定されたフィールドと値に一致する店舗ドキュメントをFirestoreから検索します。
func (r *StoreRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Store, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var stores []*models.Store
	for _, doc := range docs {
		var store models.Store
		if err := doc.DataTo(&store); err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}
	return stores, nil
}

// UpdateByIDは、指定されたIDの店舗ドキュメントをFirestoreで更新します。
func (r *StoreRepository) UpdateByID(ctx context.Context, id string, store *models.Store) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Set(ctx, store)
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
