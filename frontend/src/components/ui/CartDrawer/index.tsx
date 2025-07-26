'use client';

import React from 'react';
import {
  List,
  Typography,
  Button,
  InputNumber,
  Space,
  Divider,
  Empty,
  Input,
  message,
} from 'antd';
import { DeleteOutlined, ShoppingCartOutlined } from '@ant-design/icons';
import { CartDrawerProps } from './CartDrawer.types';
import { useUpdateCartItem, useRemoveFromCart, useClearCart } from '@/hooks/api/useOrders';

const { Title, Text } = Typography;
const { TextArea } = Input;

export const CartDrawer: React.FC<CartDrawerProps> = ({
  sessionId,
  cart,
  onOrderConfirm,
  loading = false,
}) => {
  const updateCartItemMutation = useUpdateCartItem(sessionId, '');
  const removeFromCartMutation = useRemoveFromCart(sessionId, '');
  const clearCartMutation = useClearCart(sessionId);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price);
  };

  const handleQuantityChange = async (itemId: string, quantity: number) => {
    if (quantity <= 0) {
      await handleRemoveItem(itemId);
      return;
    }

    try {
      await updateCartItemMutation.mutateAsync({ quantity });
    } catch (error) {
      message.error('数量の更新に失敗しました');
    }
  };

  const handleRemoveItem = async (itemId: string) => {
    try {
      await removeFromCartMutation.mutateAsync();
      message.success('商品をカートから削除しました');
    } catch (error) {
      message.error('商品の削除に失敗しました');
    }
  };

  const handleClearCart = async () => {
    try {
      await clearCartMutation.mutateAsync();
      message.success('カートを空にしました');
    } catch (error) {
      message.error('カートのクリアに失敗しました');
    }
  };

  const handleOrderConfirm = () => {
    if (!cart || cart.items.length === 0) {
      message.warning('カートが空です');
      return;
    }
    onOrderConfirm?.();
  };

  if (!cart || cart.items.length === 0) {
    return (
      <div className="h-full flex flex-col">
        <Empty
          image={<ShoppingCartOutlined style={{ fontSize: 64, color: '#d9d9d9' }} />}
          description="カートが空です"
          className="flex-1 flex flex-col justify-center"
        >
          <Text type="secondary">メニューから商品を選んでください</Text>
        </Empty>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col">
      {/* Cart Items */}
      <div className="flex-1 overflow-auto">
        <List
          dataSource={cart.items}
          renderItem={(item) => (
            <List.Item className="px-0">
              <div className="w-full">
                <div className="flex justify-between items-start mb-2">
                  <div className="flex-1">
                    <Title level={5} className="m-0">
                      {item.menuItem.name}
                    </Title>
                    <Text type="secondary" className="text-sm">
                      {formatPrice(item.unitPrice)} × {item.quantity}
                    </Text>
                  </div>
                  <Text strong className="text-lg">
                    {formatPrice(item.totalPrice)}
                  </Text>
                </div>

                {item.specialInstructions && (
                  <Text type="secondary" className="text-sm block mb-2">
                    備考: {item.specialInstructions}
                  </Text>
                )}

                <Space className="w-full justify-between">
                  <Space>
                    <InputNumber
                      min={1}
                      max={99}
                      value={item.quantity}
                      onChange={(value) => handleQuantityChange(item.id, value || 0)}
                      size="small"
                      className="w-16"
                    />
                  </Space>
                  <Button
                    type="text"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => handleRemoveItem(item.id)}
                    size="small"
                  >
                    削除
                  </Button>
                </Space>
              </div>
            </List.Item>
          )}
        />
      </div>

      <Divider className="my-4" />

      {/* Order Summary */}
      <div className="bg-gray-50 p-4 rounded-lg mb-4">
        <div className="flex justify-between items-center mb-2">
          <Text>商品数</Text>
          <Text>{cart.totalItems}点</Text>
        </div>
        <div className="flex justify-between items-center">
          <Title level={4} className="m-0">
            合計金額
          </Title>
          <Title level={4} className="m-0 text-primary">
            {formatPrice(cart.totalAmount)}
          </Title>
        </div>
      </div>

      {/* Action Buttons */}
      <Space direction="vertical" className="w-full">
        <Button
          type="primary"
          size="large"
          onClick={handleOrderConfirm}
          className="w-full h-12"
          disabled={cart.items.length === 0}
        >
          注文内容を確認する ({formatPrice(cart.totalAmount)})
        </Button>
        <Button
          type="default"
          onClick={handleClearCart}
          loading={clearCartMutation.isPending}
          className="w-full"
        >
          カートを空にする
        </Button>
      </Space>

      <style jsx global>{`
        .ant-list-item {
          border-bottom: 1px solid #f0f0f0;
          padding: 16px 0;
        }
        
        .ant-list-item:last-child {
          border-bottom: none;
        }
        
        @media (max-width: 768px) {
          .ant-typography h5 {
            font-size: 14px;
          }
          
          .ant-btn-lg {
            height: 48px;
            font-size: 16px;
          }
        }
      `}</style>
    </div>
  );
};