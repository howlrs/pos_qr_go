package repositories

import (
	"context"

	"backend/models"

	"cloud.google.com/go/firestore"
)

// SeatRepository は Firestore の seats コレクションを操作するためのリポジトリです。
type SeatRepository struct {
	client     *firestore.Client
	collection string
}

// NewSeatRepository は新しい SeatRepository のインスタンスを生成します。
func NewSeatRepository(client *firestore.Client) *SeatRepository {
	return &SeatRepository{
		client:     client,
		collection: "seats",
	}
}

// Create は新しい座席を Firestore に作成します。
func (r *SeatRepository) Create(ctx context.Context, seat *models.Seat) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(seat.ID).Set(ctx, seat)
	return err
}

// Read はすべての座席情報を Firestore から読み取ります。
func (r *SeatRepository) Read(ctx context.Context) ([]*models.Seat, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var seats []*models.Seat
	for _, doc := range docs {
		var seat models.Seat
		if err := doc.DataTo(&seat); err != nil {
			return nil, err
		}
		seats = append(seats, &seat)
	}
	return seats, nil
}

// FindByID は指定されたIDの座席情報を Firestore から検索します。
func (r *SeatRepository) FindByID(ctx context.Context, id string) (*models.Seat, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var seat models.Seat
	if err := doc.DataTo(&seat); err != nil {
		return nil, err
	}
	return &seat, nil
}

// FindByField は指定されたフィールドと値に一致する座席情報を Firestore から検索します。
func (r *SeatRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Seat, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Where(field, "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var seats []*models.Seat
	for _, doc := range docs {
		var seat models.Seat
		if err := doc.DataTo(&seat); err != nil {
			return nil, err
		}
		seats = append(seats, &seat)
	}
	return seats, nil
}

// UpdateByID は指定されたIDの座席情報を Firestore で更新します。
func (r *SeatRepository) UpdateByID(ctx context.Context, id string, seat *models.Seat) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Set(ctx, seat)
	return err
}

// DeleteByID は指定されたIDの座席情報を Firestore から削除します。
func (r *SeatRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Delete(ctx)
	return err
}

// Count は Firestore に保存されている座席の総数を返します。
func (r *SeatRepository) Count(ctx context.Context) (int, error) {
	docs, err := r.client.Collection(GetCollectionName(r.collection)).Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

// Exists は指定されたIDの座席が Firestore に存在するかどうかを確認します。
func (r *SeatRepository) Exists(ctx context.Context, id string) (bool, error) {
	doc, err := r.client.Collection(GetCollectionName(r.collection)).Doc(id).Get(ctx)
	if err != nil {
		// ドキュメントが存在しない場合もエラーが返るため、falseを返す
		return false, nil
	}
	return doc.Exists(), nil
}
