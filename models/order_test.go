package models

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	productID := "prod_123"
	quantity := 2
	price := 1500.50

	order := NewOrder(productID, quantity, price)

	assert.NotNil(t, order, "NewOrderはnilを返すべきではありません")

	// OrderIDのプレフィックスを確認
	assert.True(t, strings.HasPrefix(order.OrderID, OrderPrefix), "OrderIDは正しいプレフィックスで始まる必要があります")
	assert.Len(t, order.OrderID, len(OrderPrefix)+20, "OrderIDの長さが正しくありません") // xidの長さを考慮

	assert.Equal(t, productID, order.ProductID, "ProductIDが正しく設定されていません")
	assert.Equal(t, quantity, order.Quantity, "Quantityが正しく設定されていません")
	assert.Equal(t, price, order.Price, "Priceが正しく設定されていません")

	// 時間が正しく設定されているか確認
	assert.WithinDuration(t, time.Now().UTC(), order.CreatedAt, time.Second, "CreatedAtが現在時刻に設定されていません")
	assert.Equal(t, order.CreatedAt, order.UpdatedAt, "CreatedAtとUpdatedAtは同じである必要があります")
}

func TestOrder_Subtotal(t *testing.T) {
	testCases := []struct {
		name     string
		quantity int
		price    float64
		expected float64
	}{
		{
			name:     "通常のケース",
			quantity: 2,
			price:    10.5,
			expected: 21.0,
		},
		{
			name:     "数量が0のケース",
			quantity: 0,
			price:    100.0,
			expected: 0.0,
		},
		{
			name:     "価格が0のケース",
			quantity: 5,
			price:    0.0,
			expected: 0.0,
		},
		{
			name:     "数量と価格が両方0のケース",
			quantity: 0,
			price:    0.0,
			expected: 0.0,
		},
		{
			name:     "数量が1のケース",
			quantity: 1,
			price:    99.99,
			expected: 99.99,
		},
		{
			name:     "大きい数量と価格",
			quantity: 1000,
			price:    12345.67,
			expected: 12345670.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := &Order{
				ProductID: "test_product",
				Quantity:  tc.quantity,
				Price:     tc.price,
			}

			subtotal := order.Subtotal()
			assert.InDelta(t, tc.expected, subtotal, 0.001, "小計が正しく計算されていません")
		})
	}
}

func TestOrder_IDGeneration(t *testing.T) {
	// 連続して生成してもIDがユニークであることを確認
	order1 := NewOrder("p1", 1, 100)
	order2 := NewOrder("p2", 2, 200)

	assert.NotNil(t, order1)
	assert.NotNil(t, order2)
	assert.NotEqual(t, order1.OrderID, order2.OrderID, "連続して生成されたOrderIDはユニークである必要があります")
}
