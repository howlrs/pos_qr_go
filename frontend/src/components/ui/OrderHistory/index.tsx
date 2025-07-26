'use client';

import React from 'react';
import { List, Typography, Tag, Space, Empty, Card, Divider } from 'antd';
import { ClockCircleOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { OrderHistoryProps } from './OrderHistory.types';
import { ORDER_STATUS_LABELS, ORDER_STATUS_COLORS } from '@/types/models';

const { Title, Text } = Typography;

export const OrderHistory: React.FC<OrderHistoryProps> = ({
  orders,
  loading = false,
}) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price);
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
      case 'confirmed':
      case 'preparing':
        return <ClockCircleOutlined />;
      case 'ready':
      case 'served':
        return <CheckCircleOutlined />;
      default:
        return <ClockCircleOutlined />;
    }
  };

  if (!orders || orders.length === 0) {
    return (
      <Empty
        description="注文履歴がありません"
        className="py-8"
      >
        <Text type="secondary">まだ注文をしていません</Text>
      </Empty>
    );
  }

  return (
    <div className="order-history">
      <List
        loading={loading}
        dataSource={orders}
        renderItem={(order) => (
          <List.Item className="px-0">
            <Card className="w-full" size="small">
              {/* Order Header */}
              <div className="flex justify-between items-start mb-3">
                <div>
                  <Title level={5} className="m-0">
                    注文番号: {order.orderNumber}
                  </Title>
                  <Text type="secondary" className="text-sm">
                    {formatDateTime(order.placedAt)}
                  </Text>
                </div>
                <Tag
                  color={ORDER_STATUS_COLORS[order.status]}
                  icon={getStatusIcon(order.status)}
                  className="ml-2"
                >
                  {ORDER_STATUS_LABELS[order.status]}
                </Tag>
              </div>

              {/* Order Items */}
              <div className="mb-3">
                {order.items.map((item, index) => (
                  <div key={item.id} className="flex justify-between items-center py-1">
                    <div className="flex-1">
                      <Text className="text-sm">
                        {item.menuItem.name} × {item.quantity}
                      </Text>
                      {item.specialInstructions && (
                        <Text type="secondary" className="text-xs block">
                          備考: {item.specialInstructions}
                        </Text>
                      )}
                    </div>
                    <Text className="text-sm">
                      {formatPrice(item.totalPrice)}
                    </Text>
                  </div>
                ))}
              </div>

              {/* Special Instructions */}
              {order.specialInstructions && (
                <div className="mb-3">
                  <Text type="secondary" className="text-sm">
                    特別な要望: {order.specialInstructions}
                  </Text>
                </div>
              )}

              <Divider className="my-2" />

              {/* Order Summary */}
              <div className="flex justify-between items-center">
                <Space>
                  {order.estimatedReadyAt && (
                    <Text type="secondary" className="text-sm">
                      <ClockCircleOutlined className="mr-1" />
                      予定時刻: {formatDateTime(order.estimatedReadyAt)}
                    </Text>
                  )}
                </Space>
                <Text strong className="text-lg">
                  合計: {formatPrice(order.totalAmount)}
                </Text>
              </div>

              {/* Completion Time */}
              {order.completedAt && (
                <div className="mt-2">
                  <Text type="secondary" className="text-sm">
                    <CheckCircleOutlined className="mr-1" />
                    完了時刻: {formatDateTime(order.completedAt)}
                  </Text>
                </div>
              )}
            </Card>
          </List.Item>
        )}
      />

      <style jsx global>{`
        .order-history .ant-list-item {
          border-bottom: none;
          padding: 8px 0;
        }
        
        .order-history .ant-card {
          border: 1px solid #f0f0f0;
          border-radius: 8px;
        }
        
        .order-history .ant-card:hover {
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        
        @media (max-width: 768px) {
          .order-history .ant-typography h5 {
            font-size: 14px;
          }
          
          .order-history .ant-tag {
            font-size: 11px;
            padding: 2px 6px;
          }
        }
      `}</style>
    </div>
  );
};