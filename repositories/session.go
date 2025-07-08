package repositories

// SessionRepository はセッションの永続化を管理します。
import (
	"backend/models"
	"context"

	"cloud.google.com/go/firestore"
)

type SessionRepository struct {
	client     *firestore.Client
	collection string
}

// NewSessionRepository は新しいSessionRepositoryを生成します。
func NewSessionRepository(client *firestore.Client) Repository[models.Session] {

	return &SessionRepository{
		client:     client,
		collection: "sessions",
	}
}

// Create は新しいセッションをFirestoreに作成します。
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(session.ID).Set(ctx, session)
	return err
}

// Read はすべてのセッションをFirestoreから読み取ります。
func (r *SessionRepository) Read(ctx context.Context) ([]*models.Session, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session
	for _, doc := range docs {
		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// FindByID はIDを使用してセッションを検索します。
func (r *SessionRepository) FindByID(ctx context.Context, id string) (*models.Session, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := doc.DataTo(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByField は指定されたフィールドと値に一致するセッションを検索します。
func (r *SessionRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Session, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session
	for _, doc := range docs {
		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// UpdateByID はIDを使用してセッションを更新します。
func (r *SessionRepository) UpdateByID(ctx context.Context, id string, session *models.Session) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Set(ctx, session)
	return err
}

// DeleteByID はIDを使用してセッションを削除します。
func (r *SessionRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Delete(ctx)
	return err
}

// Count はセッションの総数を返します。
func (r *SessionRepository) Count(ctx context.Context) (int, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

// Exists は指定されたIDのセッションが存在するかどうかを確認します。
func (r *SessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return false, nil
	}
	return doc.Exists(), nil
}
