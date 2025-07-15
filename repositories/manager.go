package repositories

// manager.go は管理者 (Manager) モデルに関するデータベース操作を実装します。
// 管理者とは、店舗の発行・管理を行うユーザーを指します。

import (
	"context"

	"backend/models"

	"cloud.google.com/go/firestore"
)

// ManagerRepository は Firestore の "managers" コレクションとやり取りするためのリポジトリです。
type ManagerRepository struct {
	client     *firestore.Client
	collection string
}

// NewManagerRepository は新しい ManagerRepository のインスタンスを生成します。
func NewManagerRepository(client *firestore.Client) Repository[models.Manager] {
	if client == nil {
		return NewMockManagerRepository()
	}
	return &ManagerRepository{
		client:     client,
		collection: "managers",
	}
}

// Create は新しい管理者ドキュメントを Firestore に作成します。
func (r *ManagerRepository) Create(ctx context.Context, manager *models.Manager) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(manager.Email).Set(ctx, manager)
	return err
}

// Read は Firestore からすべての管理者ドキュメントを取得します。
func (r *ManagerRepository) Read(ctx context.Context) ([]*models.Manager, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return []*models.Manager{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	managers := make([]*models.Manager, len(docs))
	for i, doc := range docs {
		manager := &models.Manager{}
		if err := doc.DataTo(manager); err != nil {
			return nil, err
		}
		managers[i] = manager
	}

	return managers, nil
}

// FindByID は指定された ID を持つ管理者ドキュメントを Firestore から取得します。
func (r *ManagerRepository) FindByID(ctx context.Context, id string) (*models.Manager, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	manager := &models.Manager{}
	if err := doc.DataTo(manager); err != nil {
		return nil, err
	}

	return manager, nil
}

// FindByField は指定されたフィールドと値に一致する管理者ドキュメントを Firestore から検索します。
func (r *ManagerRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Manager, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return []*models.Manager{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	managers := make([]*models.Manager, len(docs))
	for i, doc := range docs {
		manager := &models.Manager{}
		if err := doc.DataTo(manager); err != nil {
			return nil, err
		}
		managers[i] = manager
	}

	return managers, nil
}

// UpdateByID は指定された ID を持つ管理者ドキュメントを Firestore 上で更新します。
func (r *ManagerRepository) UpdateByID(ctx context.Context, id string, manager *models.Manager) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Set(ctx, manager)
	return err
}

// DeleteByID は指定された ID を持つ管理者ドキュメントを Firestore から削除します。
func (r *ManagerRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Delete(ctx)
	return err
}

// Count は Firestore 内の管理者ドキュメントの総数を返します。
func (r *ManagerRepository) Count(ctx context.Context) (int, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

// Exists は指定された ID を持つ管理者ドキュメントが Firestore に存在するかどうかを確認します。
func (r *ManagerRepository) Exists(ctx context.Context, id string) (bool, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		// ドキュメントが存在しない場合もエラーが返るため、ここでは false を返す
		return false, nil
	}
	return doc.Exists(), nil
}
