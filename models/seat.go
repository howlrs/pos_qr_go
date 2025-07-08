package models

import (
	"time"
)

// GenerateUniqueID関数は、Rustの `generate_unique_id` に相当し、
// 外部パッケージで定義されていることを想定しています。
// 例:
// func GenerateUniqueID(prefix string, seed *string) string { /* ... */ }

// UserSeatPrefix はSeatエンティティのIDを生成する際のプレフィックスです。
// Rustの `WhereIsIdPrefix::UserSeat` に相当します。
const UserSeatPrefix = "seat_" // 実際のプレフィックス文字列に置き換えてください

// Seat は座席エンティティを表します。
type Seat struct {
	ID        string    `json:"id" db:"id" firestore:"id"`
	Name      string    `json:"name" db:"name" firestore:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" firestore:"updated_at"`
}

// NewSeat は新しいSeatインスタンスを作成します。
//
// # 引数
//   - `name`: 座席の名前。
//
// # 戻り値
//
// 新しいSeatインスタンスへのポインタ。
//
// # 例
//
//	seat := models.NewSeat("My Seat")
//	fmt.Println(seat.Name) // 出力: "My Seat"
func NewSeat(name string) *Seat {
	uid := GenerateUniqueID(UserSeatPrefix)
	now := time.Now().UTC()

	return &Seat{
		ID:        uid,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
