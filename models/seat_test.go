package models

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSeat(t *testing.T) {
	seatName := "Table 1"

	seat := NewSeat(seatName)

	assert.NotNil(t, seat, "NewSeatはnilを返すべきではありません")

	// IDのプレフィックスと長さを確認
	assert.True(t, strings.HasPrefix(seat.ID, UserSeatPrefix), "IDは正しいプレフィックスで始まる必要があります")
	assert.Len(t, seat.ID, len(UserSeatPrefix)+20, "IDの長さが正しくありません") // xidの長さを考慮

	// Nameが正しく設定されているか確認
	assert.Equal(t, seatName, seat.Name, "Nameが正しく設定されていません")

	// 時間が正しく設定されているか確認
	assert.WithinDuration(t, time.Now().UTC(), seat.CreatedAt, time.Second, "CreatedAtが現在時刻に設定されていません")
	assert.Equal(t, seat.CreatedAt, seat.UpdatedAt, "CreatedAtとUpdatedAtは同じである必要があります")
}

func TestSeat_IDGeneration(t *testing.T) {
	// 連続して生成してもIDがユニークであることを確認
	seat1 := NewSeat("Seat A")
	seat2 := NewSeat("Seat B")

	assert.NotNil(t, seat1)
	assert.NotNil(t, seat2)
	assert.NotEqual(t, seat1.ID, seat2.ID, "連続して生成されたIDはユニークである必要があります")
	assert.NotEqual(t, seat1.Name, seat2.Name, "名前も異なるはずです")
}

func TestNewSeat_EmptyName(t *testing.T) {
	seat := NewSeat("")

	assert.NotNil(t, seat)
	assert.Equal(t, "", seat.Name, "空の名前でも正しく設定されるべきです")
	assert.True(t, strings.HasPrefix(seat.ID, UserSeatPrefix), "IDは空の名前でも生成されるべきです")
}
