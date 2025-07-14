package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStatus_BooleanFlags は、bool値を返す各ステータスメソッドをまとめてテストします。
func TestStatus_BooleanFlags(t *testing.T) {
	testCases := []struct {
		status            Status
		isFinal           bool
		isReadyForService bool
		canAddItem        bool
		canCancel         bool
		isFulfilled       bool
		isInPreparation   bool
	}{
		// --- 初期状態 ---
		{StatusCreated, false, false, true, true, false, false},
		{StatusPendingPayment, false, false, true, true, false, false},
		// 支払いを受けた場合は、通常は次のステータスに進むため、キャンセル不可
		{StatusPaymentReceived, false, false, false, false, false, false},

		// --- 内部処理 ---
		{StatusConfirmed, false, false, true, true, false, true},
		{StatusPreparing, false, false, false, true, false, true},
		{StatusOnHold, false, false, true, true, false, true},

		// --- 提供準備完了 ---
		// 現状では、ピックアップ後である配送準備完了はキャンセル可能
		{StatusReadyForPickup, false, true, false, true, false, false},
		{StatusReadyForDelivery, false, true, false, true, false, false},

		// --- 提供/配送中 ---
		// 現状では、ピックアップ後である配送中はキャンセル不可
		{StatusOutForDelivery, false, true, false, false, false, false},
		{StatusPickedUp, false, false, false, false, true, false},
		{StatusServed, false, false, false, false, true, false},
		{StatusDelivered, false, false, false, false, true, false},

		// --- 最終状態 ---
		// 完了、返金、失敗などの最終状態
		{StatusCompleted, true, false, false, false, false, false},
		{StatusCancelled, true, false, false, false, false, false},
		{StatusRefunded, true, false, false, false, false, false},
		{StatusPartiallyRefunded, false, false, false, false, false, false}, // 部分返金は限定的な遷移が可能
		{StatusPaymentFailed, true, false, false, false, false, false},
		{StatusFailed, true, false, false, false, false, false},
	}

	for _, tc := range testCases {
		t.Run(string(tc.status), func(t *testing.T) {
			assert.Equal(t, tc.isFinal, tc.status.IsFinal(), "IsFinal() の結果が正しくありません")
			assert.Equal(t, tc.isReadyForService, tc.status.IsReadyForService(), "IsReadyForService() の結果が正しくありません")
			assert.Equal(t, tc.canAddItem, tc.status.CanAddItem(), "CanAddItem() の結果が正しくありません")
			assert.Equal(t, tc.canCancel, tc.status.CanCancel(), fmt.Sprintf("CanCancel() の結果が正しくありません: %s", tc.status))
			assert.Equal(t, tc.isFulfilled, tc.status.IsFulfilled(), "IsFulfilled() の結果が正しくありません")
			assert.Equal(t, tc.isInPreparation, tc.status.IsInPreparation(), "IsInPreparation() の結果が正しくありません")
		})
	}
}

// TestStatus_CanTransitionTo は、状態遷移のロジックをテストします。
func TestStatus_CanTransitionTo(t *testing.T) {
	testCases := []struct {
		from     Status
		to       Status
		expected bool
	}{
		// --- 正常な遷移 ---
		{StatusCreated, StatusPendingPayment, true},
		{StatusCreated, StatusConfirmed, true},
		{StatusCreated, StatusCancelled, true},
		{StatusConfirmed, StatusPreparing, true},
		{StatusPreparing, StatusReadyForPickup, true},
		{StatusReadyForPickup, StatusPickedUp, true},
		{StatusPickedUp, StatusCompleted, true},
		{StatusDelivered, StatusCompleted, true},
		{StatusCompleted, StatusRefunded, false}, // 完了から返金は不可（モデル上）
		// --- 返金関連の遷移 ---
		// 部分返金からの許可された遷移
		{StatusPartiallyRefunded, StatusRefunded, true},  // 残りの金額も返金
		{StatusPartiallyRefunded, StatusCompleted, true}, // 返金後に完了扱い
		{StatusPartiallyRefunded, StatusCancelled, true}, // 部分返金後にキャンセル
		// 部分返金からの不許可遷移
		{StatusPartiallyRefunded, StatusCreated, false},
		{StatusPartiallyRefunded, StatusPreparing, false},
		{StatusPartiallyRefunded, StatusConfirmed, false},

		// --- 不正な遷移 ---
		{StatusCreated, StatusCompleted, false},
		{StatusPreparing, StatusCreated, false},
		{StatusPickedUp, StatusPreparing, false},
		{StatusCancelled, StatusConfirmed, false}, // 最終状態からの遷移
		{StatusCompleted, StatusCancelled, false}, // 最終状態からの遷移

		// --- 自己遷移 ---
		{StatusCreated, StatusCreated, false},
		{StatusPreparing, StatusPreparing, false},
		{StatusCompleted, StatusCompleted, false},

		// --- その他エッジケース ---
		{StatusPaymentFailed, StatusCancelled, true},
		{StatusOnHold, StatusPreparing, true},
		{StatusOnHold, StatusCancelled, true},
	}

	// status.goで定義されているすべての遷移をテスト
	allStatuses := []Status{
		StatusCreated, StatusPendingPayment, StatusPaymentReceived, StatusPaymentFailed,
		StatusPendingConfirmation, StatusConfirmed, StatusPreparing, StatusOnHold,
		StatusReadyForPickup, StatusReadyForDelivery, StatusOutForDelivery,
		StatusDeliveryAttemptFailed, StatusDelivered, StatusPickedUp, StatusServed,
		StatusCompleted, StatusCancelled, StatusDeclined, StatusRefunded,
		StatusPartiallyRefunded, StatusFailed,
	}

	// CanTransitionTo の switch-case に基づいて動的にテストケースを生成
	for _, from := range allStatuses {
		for _, to := range allStatuses {
			// 期待値を計算
			expected := isTransitionAllowedByLogic(from, to)
			// テストケースを追加
			testCases = append(testCases, struct {
				from     Status
				to       Status
				expected bool
			}{from, to, expected})
		}
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_to_%s", tc.from, tc.to), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.from.CanTransitionTo(tc.to), fmt.Sprintf("CanTransitionTo(%s, %s) の結果(%t)が正しくありません", tc.from, tc.to, tc.expected))
		})
	}
}

// isTransitionAllowedByLogic は、status.goのCanTransitionToのロジックをミラーリングしたヘルパー関数です。
// これにより、実装とテストが一致していることを保証します。
func isTransitionAllowedByLogic(s, newStatus Status) bool {
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
