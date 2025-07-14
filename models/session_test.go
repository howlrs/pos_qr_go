package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ヘルパー関数 ---

// テスト用の基本的なSessionインスタンスを作成
func newTestSession(t *testing.T) *Session {
	items := []Order{
		*NewOrder("prod_1", 2, 100),
		*NewOrder("prod_2", 1, 50),
	}
	session, err := NewSession("store_123", "seat_456", items)
	require.NoError(t, err)
	require.NotNil(t, session)
	return session
}

// --- テストケース ---

func TestNewSession(t *testing.T) {
	t.Run("正常なケース", func(t *testing.T) {
		items := []Order{
			*NewOrder("prod_A", 1, 150),
			*NewOrder("prod_B", 3, 200),
		}
		storeID := "store_abc"
		seatID := "seat_xyz"

		session, err := NewSession(storeID, seatID, items)

		assert.NoError(t, err)
		assert.NotNil(t, session)

		assert.Equal(t, storeID, session.StoreID)
		assert.Equal(t, seatID, session.SeatID)
		assert.Len(t, session.Items, 2)
		assert.Equal(t, 150.0+600.0, session.TotalAmount)
		assert.Equal(t, StatusCreated, session.Status)
		assert.NotEmpty(t, session.ID)

		// 時間関連のフィールドを確認
		now := time.Now().UTC()
		assert.WithinDuration(t, now, session.IssuedAt, time.Second)
		assert.WithinDuration(t, now, session.CreatedAt, time.Second)
		assert.WithinDuration(t, now, session.UpdatedAt, time.Second)
		assert.WithinDuration(t, now.Add(15*time.Minute), session.ExpiresAt, time.Second)

		// 各アイテムにOrderIDが設定されているか確認
		for _, item := range session.Items {
			assert.Equal(t, session.ID, item.OrderID)
		}
	})

	t.Run("引数が不正なケース", func(t *testing.T) {
		items := []Order{*NewOrder("p1", 1, 1)}
		_, err := NewSession("", "seat1", items)
		assert.ErrorIs(t, err, ErrInvalidArgument)

		_, err = NewSession("store1", "", items)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("アイテムが空のケース", func(t *testing.T) {
		_, err := NewSession("store1", "seat1", []Order{})
		assert.ErrorIs(t, err, ErrNoItems)

		_, err = NewSession("store1", "seat1", nil)
		assert.ErrorIs(t, err, ErrNoItems)
	})
}

func TestSession_AddItem(t *testing.T) {
	t.Run("正常に追加できるケース", func(t *testing.T) {
		session := newTestSession(t)
		initialTotal := session.TotalAmount
		initialUpdatedAt := session.UpdatedAt
		time.Sleep(10 * time.Millisecond) // 更新時刻が変わるように少し待つ

		newItem := *NewOrder("prod_3", 1, 300)
		err := session.AddItem(newItem)

		assert.NoError(t, err)
		assert.Len(t, session.Items, 3)
		assert.Equal(t, initialTotal+300, session.TotalAmount)
		assert.Equal(t, session.ID, session.Items[2].OrderID)
		assert.True(t, session.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("追加できないステータスのケース", func(t *testing.T) {
		session := newTestSession(t)
		session.Status = StatusCompleted // 最終ステータスに設定

		err := session.AddItem(*NewOrder("p", 1, 1))
		assert.Error(t, err)
		var e *CannotAddItemError
		assert.ErrorAs(t, err, &e)
		assert.Equal(t, StatusCompleted, e.Status)
	})

	t.Run("有効期限切れのケース", func(t *testing.T) {
		session := newTestSession(t)
		session.ExpiresAt = time.Now().UTC().Add(-time.Minute) // 期限切れに設定

		err := session.AddItem(*NewOrder("p", 1, 1))
		assert.ErrorIs(t, err, ErrOrderExpired)
	})
}

func TestSession_UpdateStatus(t *testing.T) {
	t.Run("正常な状態遷移", func(t *testing.T) {
		session := newTestSession(t)
		initialUpdatedAt := session.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		err := session.UpdateStatus(StatusConfirmed)
		assert.NoError(t, err)
		assert.Equal(t, StatusConfirmed, session.Status)
		assert.True(t, session.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("不正な状態遷移", func(t *testing.T) {
		session := newTestSession(t)
		err := session.UpdateStatus(StatusCompleted) // CreatedからCompletedへは直接遷移できない
		assert.Error(t, err)
		var e *InvalidStatusTransitionError
		assert.ErrorAs(t, err, &e)
		assert.Equal(t, StatusCreated, e.From)
		assert.Equal(t, StatusCompleted, e.To)
	})

	t.Run("最終ステータスからの更新", func(t *testing.T) {
		session := newTestSession(t)
		session.Status = StatusCancelled
		err := session.UpdateStatus(StatusConfirmed)
		assert.ErrorIs(t, err, ErrOrderAlreadyFinal)
	})
}

func TestSession_RecalculateTotalAmount(t *testing.T) {
	session := newTestSession(t)
	session.Items = append(session.Items, *NewOrder("prod_4", 1, 1000))
	// この時点ではTotalAmountは古いまま
	assert.NotEqual(t, 250.0+1000.0, session.TotalAmount)

	session.RecalculateTotalAmount()
	assert.Equal(t, 250.0+1000.0, session.TotalAmount)
}

func TestSession_StatusHelpers(t *testing.T) {
	t.Run("MarkConfirmOrder", func(t *testing.T) {
		s := newTestSession(t)
		err := s.MarkConfirmOrder()
		assert.NoError(t, err)
		assert.Equal(t, StatusConfirmed, s.Status)
	})

	t.Run("MarkCompleteOrder", func(t *testing.T) {
		s := newTestSession(t)
		// 途中のステータスを経由する必要がある
		s.Status = StatusDelivered
		err := s.MarkCompleteOrder()
		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, s.Status)
	})

	t.Run("MarkCancelOrder", func(t *testing.T) {
		s := newTestSession(t)
		err := s.MarkCancelOrder()
		assert.NoError(t, err)
		assert.Equal(t, StatusCancelled, s.Status)
	})
}

func TestSession_MarkRefundPartially(t *testing.T) {
	t.Run("正常な部分返金", func(t *testing.T) {
		s := newTestSession(t)
		s.Status = StatusCompleted // 返金は完了後などから行われる想定
		err := s.MarkRefundPartially(100.0)
		assert.NoError(t, err)
		assert.Equal(t, StatusPartiallyRefunded, s.Status)
	})

	t.Run("返金額が合計を超えるケース", func(t *testing.T) {
		s := newTestSession(t)
		s.Status = StatusCompleted // 返金は完了後などから行われる想定
		err := s.MarkRefundPartially(s.TotalAmount + 1)
		assert.ErrorIs(t, err, ErrRefundAmountExceedsTotal)
	})
}
