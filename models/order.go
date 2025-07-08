package models

import "time"

// OrderPrefix はOrderエンティティのIDを生成する際のプレフィックスです。
const OrderPrefix = "order_" // 実際のプレフィックス文字列に置き換えてください

// --- OrderItem プレースホルダー ---

// Order は注文内の個々の商品を表します。
// この構造体は外部で定義されていることを想定しています。
type Order struct {
	OrderID   string    `json:"order_id" db:"order_id" firestore:"order_id"`
	ProductID string    `json:"product_id" db:"product_id" firestore:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity" firestore:"quantity"`
	Price     float64   `json:"price" db:"price" firestore:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" firestore:"updated_at"`
}

// NewOrder は新しい注文アイテムを作成します。
func NewOrder(productID string, quantity int, price float64) *Order {
	uid := GenerateUniqueID(UserSeatPrefix)
	now := time.Now().UTC()
	return &Order{
		OrderID:   uid,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Subtotal はこの注文アイテムの小計を計算します。
func (oi *Order) Subtotal() float64 {
	return oi.Price * float64(oi.Quantity)
}
