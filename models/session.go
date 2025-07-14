package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

// --- エラー定義 ---

// OrderErrorに対応するGoのカスタムエラー
var (
	ErrInvalidArgument          = errors.New("store_idまたはseat_idが無効です")
	ErrNoItems                  = errors.New("注文に商品が含まれていません")
	ErrOrderExpired             = errors.New("注文の有効期限が切れています")
	ErrOrderAlreadyFinal        = errors.New("注文はすでに最終状態のため更新できません")
	ErrRefundAmountExceedsTotal = errors.New("返金額が合計金額を超えています")
)

// 動的なエラーを生成するためのカスタムエラー型
type CannotAddItemError struct {
	Status Status
}

func (e *CannotAddItemError) Error() string {
	return fmt.Sprintf("現在のステータス(%s)では商品を追加できません", e.Status)
}

type InvalidStatusTransitionError struct {
	From Status
	To   Status
}

func (e *InvalidStatusTransitionError) Error() string {
	return fmt.Sprintf("ステータスを '%s' から '%s' に遷移させることはできません", e.From, e.To)
}

// --- Session 構造体とメソッド ---

// Session は注文のモデルを表す構造体です。
// 注文は、店舗、ユーザー席、商品、合計金額、ステータスなどの情報を含みます。
type Session struct {
	ID          string
	StoreID     string
	SeatID      string
	Items       []Order
	TotalAmount float64
	Status      Status

	ExpiresAt time.Time
	IssuedAt  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSession は新しい注文を作成するコンストラクタです。
// 初期ステータスは`Created`に設定されます。
func NewSession(storeID, seatID string, items []Order) (*Session, error) {
	if storeID == "" || seatID == "" {
		return nil, ErrInvalidArgument
	}
	if len(items) == 0 {
		return nil, ErrNoItems
	}

	orderID := xid.New().String()
	now := time.Now().UTC()
	var totalAmount float64

	// 各OrderItemに親である注文IDを設定し、合計金額を計算
	for i := range items {
		items[i].OrderID = orderID
		totalAmount += items[i].Subtotal()
	}

	return &Session{
		ID:          orderID,
		StoreID:     storeID,
		SeatID:      seatID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      StatusCreated,
		ExpiresAt:   now.Add(15 * time.Minute), // 注文の有効期限は15分後
		IssuedAt:    now,
		CreatedAt:   now, // FirestoreのserverTimestampが使えない場合のフォールバック
		UpdatedAt:   now,
	}, nil
}

// AddItem は注文に新しいアイテムを追加します。
// `Created`, `PendingPayment`, `PendingConfirmation`, `Confirmed`, `OnHold` ステータスでのみ追加可能です。
func (s *Session) AddItem(newItem Order) error {
	if !s.Status.CanAddItem() {
		return &CannotAddItemError{Status: s.Status}
	}
	if s.ExpiresAt.Before(time.Now().UTC()) {
		return ErrOrderExpired
	}

	newItem.OrderID = s.ID
	s.Items = append(s.Items, newItem)
	s.RecalculateTotalAmount() // アイテム追加後に合計金額を更新
	s.setUpdatedAt()
	return nil
}

// UpdateStatus は注文のステータスを更新します。
// 不正な状態遷移をチェックします。
func (s *Session) UpdateStatus(newStatus Status) error {
	if s.Status.IsFinal() {
		return fmt.Errorf("%w: 現在のステータスは '%s'", ErrOrderAlreadyFinal, s.Status)
	}
	if !s.Status.CanTransitionTo(newStatus) {
		return &InvalidStatusTransitionError{From: s.Status, To: newStatus}
	}

	s.Status = newStatus
	s.setUpdatedAt()
	return nil
}

// ExceptionUpdateStatus は注文のステータスを更新します。
// ステータスが最終状態でも更新可能です。
func (s *Session) ExceptionUpdateStatus(newStatus Status) error {
	s.Status = newStatus
	s.setUpdatedAt()
	return nil
}

// RecalculateTotalAmount は注文の合計金額を、現在のアイテムリストに基づいて再計算します。
func (s *Session) RecalculateTotalAmount() {
	var total float64
	for _, item := range s.Items {
		total += item.Subtotal()
	}
	s.TotalAmount = total
	s.setUpdatedAt()
}

// setUpdatedAt は `updated_at` フィールドを現在時刻に設定します。
// FirestoreのserverTimestampが利用されるため、これは主にGoのコード内での状態を反映させるためです。
func (s *Session) setUpdatedAt() {
	s.UpdatedAt = time.Now().UTC()
}

// --- ステータス固有のヘルパーメソッド（ビジネスロジックをカプセル化） ---

// MarkConfirmOrder は注文を店舗側確認済みにします。
func (s *Session) MarkConfirmOrder() error {
	return s.UpdateStatus(StatusConfirmed)
}

// MarkPreparing は注文を準備中にします。
func (s *Session) MarkPreparing() error {
	return s.UpdateStatus(StatusPreparing)
}

// MarkAsReadyForPickup は注文を店舗側提供可能にします。
func (s *Session) MarkAsReadyForPickup() error {
	return s.UpdateStatus(StatusReadyForPickup)
}

// MarkAsReadyForDelivery は注文を店舗側提供及び配送準備完了にします。
func (s *Session) MarkAsReadyForDelivery() error {
	return s.UpdateStatus(StatusReadyForDelivery)
}

// MarkOutForDelivery は注文を配送中にします。
func (s *Session) MarkOutForDelivery() error {
	return s.UpdateStatus(StatusOutForDelivery)
}

// MarkAsDelivered は注文の配送が完了したことを示します。
func (s *Session) MarkAsDelivered() error {
	return s.UpdateStatus(StatusDelivered)
}

// MarkAsPickedUp は注文が配膳または配送者に渡ったことを示します。
func (s *Session) MarkAsPickedUp() error {
	return s.UpdateStatus(StatusPickedUp)
}

// MarkAsServed は注文が店舗側で提供され、顧客テーブルに納品確認されたことを示します。
func (s *Session) MarkAsServed() error {
	return s.UpdateStatus(StatusServed)
}

// MarkCompleteOrder は注文の全プロセスが完了したことを示します。
func (s *Session) MarkCompleteOrder() error {
	return s.UpdateStatus(StatusCompleted)
}

// MarkCancelOrder は注文をキャンセルします。
func (s *Session) MarkCancelOrder() error {
	return s.UpdateStatus(StatusCancelled)
}

// MarkFailPayment は注文のお支払いが失敗したことを示します。
func (s *Session) MarkFailPayment() error {
	return s.UpdateStatus(StatusPaymentFailed)
}

// MarkAsOnHold は注文を一時保留にします。
func (s *Session) MarkAsOnHold() error {
	return s.UpdateStatus(StatusOnHold)
}

// MarkRefundFully は注文を全額返金済みにします。
func (s *Session) MarkRefundFully() error {
	return s.UpdateStatus(StatusRefunded)
}

// MarkRefundPartially は注文を部分返金済みにします。
// ステータスは `PartiallyRefunded` に更新されます。
// 完了状態でも返金処理が可能です。
func (s *Session) MarkRefundPartially(amount float64) error {
	if amount > s.TotalAmount {
		return ErrRefundAmountExceedsTotal
	}
	// ここに実際の返金処理（金額の記録など）を実装
	// 関数を通すと完了後のステータス変更ができないため、例外処理
	return s.ExceptionUpdateStatus(StatusPartiallyRefunded)
}
