package usecases

import (
	"backend/models"
	"time"
)

type Session struct {
	Store     *models.Store
	Seat      *models.Seat
	ExpiredAt time.Time
}

// NewSession は新しいセッションを作成します。
// storeID: ストアID
// seatID: 席ID
// expiredAt: セッションの有効期限
func NewSession(storeID, seatID string, expiredAt time.Time) *Session {
	return &Session{
		Store: &models.Store{
			ID: storeID,
		},
		Seat: &models.Seat{
			ID: seatID,
		},
		ExpiredAt: expiredAt,
	}
}

// CreateJWT はセッション情報からJWTトークンを生成します。
// 戻り値: JWTトークン文字列、エラー
func (r *Session) CreateJWT() (string, error) {
	return models.NewSessionClaims(r.Store, r.Seat, r.ExpiredAt).ToJwtToken()
}
