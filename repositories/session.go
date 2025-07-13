package repositories

// SessionRepository はセッションの永続化を管理します。
import (
	"backend/models"
	"context"
	"time"

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

type Session struct {
	ID          string  `firestore:"id"`
	StoreID     string  `firestore:"store_id"`
	SeatID      string  `firestore:"seat_id"`
	Items       []Order `firestore:"items"`
	TotalAmount float64 `firestore:"total_amount"`
	Status      Status  `firestore:"status"`

	ExpiresAt time.Time `firestore:"expires_at"`
	IssuedAt  time.Time `firestore:"issued_at"`

	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

type Order struct {
	OrderID   string    `firestore:"order_id"`
	ProductID string    `firestore:"product_id"`
	Quantity  int       `firestore:"quantity"`
	Price     float64   `firestore:"price"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

type Status string

func ToSetSession(s *models.Session) *Session {
	return &Session{
		ID:          s.ID,
		StoreID:     s.StoreID,
		SeatID:      s.SeatID,
		Items:       ToSetOrders(s.Items),
		TotalAmount: s.TotalAmount,
		Status:      Status(s.Status),
		ExpiresAt:   s.ExpiresAt,
		IssuedAt:    s.IssuedAt,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (s *Session) ToUpdate() *Session {
	s.UpdatedAt = time.Now().UTC()
	return s
}

func (s *Session) ToModel() *models.Session {
	return &models.Session{
		ID:          s.ID,
		StoreID:     s.StoreID,
		SeatID:      s.SeatID,
		Items:       ToModelOrders(s.Items),
		TotalAmount: s.TotalAmount,
		Status:      models.Status(s.Status),
		ExpiresAt:   s.ExpiresAt,
		IssuedAt:    s.IssuedAt,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func ToSetOrders(orders []models.Order) []Order {
	setOrders := make([]Order, len(orders))
	for i, o := range orders {
		setOrders[i] = Order{
			OrderID:   o.OrderID,
			ProductID: o.ProductID,
			Quantity:  o.Quantity,
			Price:     o.Price,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
		}
	}
	return setOrders
}

func ToModelOrders(orders []Order) []models.Order {
	modelOrders := make([]models.Order, len(orders))
	for i, o := range orders {
		modelOrders[i] = models.Order{
			OrderID:   o.OrderID,
			ProductID: o.ProductID,
			Quantity:  o.Quantity,
			Price:     o.Price,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
		}
	}
	return modelOrders
}

// Create は新しいセッションをFirestoreに作成します。
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(session.ID).Set(ctx, ToSetSession(session))
	return err
}

// Read はすべてのセッションをFirestoreから読み取ります。
func (r *SessionRepository) Read(ctx context.Context) ([]*models.Session, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return []*models.Session{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	sessions := make([]*models.Session, len(docs))
	for i, doc := range docs {
		var session Session
		if err := doc.DataTo(&session); err != nil {
			return nil, err
		}
		sessions[i] = session.ToModel()
	}

	return sessions, nil
}

// FindByID はIDを使用してセッションを検索します。
func (r *SessionRepository) FindByID(ctx context.Context, id string) (*models.Session, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	session := &Session{}
	if err := doc.DataTo(session); err != nil {
		return nil, err
	}

	return session.ToModel(), nil
}

// FindByField は指定されたフィールドと値に一致するセッションを検索します。
func (r *SessionRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Session, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return []*models.Session{}, nil // ドキュメントが存在しない場合は空のスライスを返す
	}

	sessions := make([]*models.Session, len(docs))
	for i, doc := range docs {
		session := &Session{}
		if err := doc.DataTo(&session); err != nil {
			return nil, err
		}
		sessions[i] = session.ToModel()
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
