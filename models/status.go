package models

// --- Status Enumとメソッド ---

// Status は飲食店における注文の現在の状態を表します。
// Rustのenum `Status` に相当します。
type Status string

const (
	// --- 初期状態と支払い ---
	StatusCreated         Status = "created"
	StatusPendingPayment  Status = "pending_payment"
	StatusPaymentReceived Status = "payment_received"
	StatusPaymentFailed   Status = "payment_failed"

	// --- レストラン内部での処理 ---
	StatusPendingConfirmation Status = "pending_confirmation"
	StatusConfirmed           Status = "confirmed"
	StatusPreparing           Status = "preparing"
	StatusOnHold              Status = "on_hold"

	// --- 提供・引き渡し・配送 ---
	StatusReadyForPickup        Status = "ready_for_pickup"
	StatusReadyForDelivery      Status = "ready_for_delivery"
	StatusOutForDelivery        Status = "out_for_delivery"
	StatusDeliveryAttemptFailed Status = "delivery_attempt_failed"
	StatusDelivered             Status = "delivered"
	StatusPickedUp              Status = "picked_up"
	StatusServed                Status = "served"

	// --- 完了と例外 ---
	StatusCompleted         Status = "completed"
	StatusCancelled         Status = "cancelled"
	StatusDeclined          Status = "declined"
	StatusRefunded          Status = "refunded"
	StatusPartiallyRefunded Status = "partially_refunded"
	StatusFailed            Status = "failed"
)

// IsFinal は注文が最終状態にあるかどうかを判定します。
// 最終状態からは通常、それ以上の主要な状態遷移はありません。
func (s Status) IsFinal() bool {
	switch s {
	case StatusCompleted, StatusCancelled, StatusDeclined, StatusRefunded,
		StatusFailed, StatusPaymentFailed:
		// 配達を外部に委託する場合
		// StatusDelivered, StatusServed, StatusPickedUp: // PickedUpも最終状態と見なすことが多い
		return true
	case StatusPartiallyRefunded:
		// 部分返金は限定的な遷移が可能なため、完全には最終状態ではない
		return false
	default:
		return false
	}
}

// IsReadyForService は注文が顧客に提供される準備ができていることを示します。
func (s Status) IsReadyForService() bool {
	switch s {
	case StatusReadyForPickup, StatusReadyForDelivery, StatusOutForDelivery:
		return true
	default:
		return false
	}
}

// CanAddItem は注文にアイテムの追加が許可されているかどうかを確認します。
func (s Status) CanAddItem() bool {
	switch s {
	case StatusCreated, StatusPendingPayment, StatusPendingConfirmation, StatusConfirmed, StatusOnHold:
		return true
	default:
		return false
	}
}

// CanCancel は注文がキャンセル可能かどうかを判定します。
// 通常、調理開始後や配送準備後はキャンセル不可とされます。
func (s Status) CanCancel() bool {
	switch s {
	case StatusCreated, StatusPendingPayment, StatusPendingConfirmation, StatusPreparing, StatusConfirmed, StatusOnHold,
		StatusReadyForPickup, StatusReadyForDelivery:
		return true
	default:
		return false
	}
}

// IsFulfilled は注文が顧客に引き渡されたか、提供された状態かどうかを判定します。
func (s Status) IsFulfilled() bool {
	switch s {
	case StatusDelivered, StatusPickedUp, StatusServed:
		return true
	default:
		return false
	}
}

// IsInPreparation は注文がまだレストランで処理中である状態かどうかを判定します。
func (s Status) IsInPreparation() bool {
	switch s {
	case StatusConfirmed, StatusPreparing, StatusOnHold:
		return true
	default:
		return false
	}
}

// CanTransitionTo は特定のステータスから別のステータスへ有効に遷移できるかを判定します。
// このメソッドはビジネスロジックの核となる部分です。
func (s Status) CanTransitionTo(newStatus Status) bool {
	switch s {
	case StatusCompleted, StatusCancelled, StatusDeclined, StatusRefunded:
		return false // 最終状態からは遷移不可
	case StatusPartiallyRefunded:
		// 部分返金からは限定的な遷移を許可
		return newStatus == StatusRefunded || newStatus == StatusCompleted || newStatus == StatusCancelled
	case StatusCreated:
		return newStatus == StatusPendingPayment || newStatus == StatusConfirmed || newStatus == StatusCancelled || newStatus == StatusDeclined
	case StatusPendingPayment:
		return newStatus == StatusPaymentReceived || newStatus == StatusPaymentFailed || newStatus == StatusCancelled || newStatus == StatusDeclined
	case StatusPaymentReceived:
		return newStatus == StatusPendingConfirmation || newStatus == StatusCancelled || newStatus == StatusDeclined
	case StatusPendingConfirmation:
		return newStatus == StatusConfirmed || newStatus == StatusCancelled || newStatus == StatusDeclined
	case StatusConfirmed:
		return newStatus == StatusPreparing || newStatus == StatusOnHold || newStatus == StatusCancelled
	case StatusPreparing:
		return newStatus == StatusReadyForPickup || newStatus == StatusReadyForDelivery || newStatus == StatusOnHold || newStatus == StatusCancelled
	case StatusOnHold:
		return newStatus == StatusConfirmed || newStatus == StatusPreparing || newStatus == StatusReadyForPickup || newStatus == StatusReadyForDelivery || newStatus == StatusCancelled
	case StatusReadyForPickup:
		return newStatus == StatusPickedUp || newStatus == StatusCancelled || newStatus == StatusServed
	case StatusReadyForDelivery:
		return newStatus == StatusOutForDelivery || newStatus == StatusCancelled || newStatus == StatusServed
	case StatusOutForDelivery:
		return newStatus == StatusDelivered || newStatus == StatusDeliveryAttemptFailed || newStatus == StatusCancelled
	case StatusDeliveryAttemptFailed:
		return newStatus == StatusOutForDelivery || newStatus == StatusDelivered || newStatus == StatusCancelled
	case StatusDelivered, StatusPickedUp, StatusServed:
		return newStatus == StatusCompleted || newStatus == StatusRefunded || newStatus == StatusPartiallyRefunded
	case StatusPaymentFailed:
		return newStatus == StatusCancelled
	case StatusFailed:
		return newStatus == StatusCancelled

	default:
		return false
	}
}
