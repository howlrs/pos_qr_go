'use client';

import React from 'react';
import {
  Modal,
  Result,
  Typography,
  Space,
  Button,
  Card,
  Divider,
  Tag,
} from 'antd';
import {
  CheckCircleOutlined,
  ClockCircleOutlined,
  ShoppingOutlined,
} from '@ant-design/icons';
import { OrderSuccessProps } from './OrderSuccess.types';
import { ORDER_STATUS_LABELS, ORDER_STATUS_COLORS } from '@/types/models';

const { Title, Text, Paragraph } = Typography;

export const OrderSuccess: React.FC<OrderSuccessProps> = ({
  visible,
  order,
  estimatedWaitTime,
  onClose,
  onContinueOrdering,
  onViewHistory,
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
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  if (!order) {
    return null;
  }

  return (
    <Modal
      open={visible}
      onCancel={onClose}
      footer={null}
      width={600}
      centered
      className="order-success-modal"
      closable={false}
    >
      <Result
        icon={<CheckCircleOutlined className="text-green-500" />}
        title="注文が完了しました！"
        subTitle={
          <Space direction="vertical" className="text-center">
            <Text className="text-lg">
              注文番号: <Text strong className="text-xl">{order.orderNumber}</Text>
            </Text>
            <Text type="secondary">
              注文時刻: {formatDateTime(order.placedAt)}
            </Text>
          </Space>
        }
      />

      <div className="space-y-6 mt-6">
        {/* Order Status */}
        <Card size="small" className="text-center">
          <Space direction="vertical" size="small">
            <Tag
              color={ORDER_STATUS_COLORS[order.status]}
              icon={<ClockCircleOutlined />}
              className="text-lg px-4 py-2"
            >
              {ORDER_STATUS_LABELS[order.status]}
            </Tag>
            {estimatedWaitTime && (
              <Text type="secondary">
                調理時間の目安: 約{estimatedWaitTime}分
              </Text>
            )}
            {order.estimatedReadyAt && (
              <Text type="secondary">
                完成予定時刻: {formatDateTime(order.estimatedReadyAt)}
              </Text>
            )}
          </Space>
        </Card>

        {/* Order Summary */}
        <Card title="注文内容" size="small">
          <div className="space-y-3">
            {order.items.map((item) => (
              <div key={item.id} className="flex justify-between items-start">
                <div className="flex-1">
                  <Text strong>{item.menuItem.name}</Text>
                  <Text type="secondary" className="text-sm block">
                    {formatPrice(item.unitPrice)} × {item.quantity}
                  </Text>
                  {item.specialInstructions && (
                    <Text type="secondary" className="text-xs block">
                      備考: {item.specialInstructions}
                    </Text>
                  )}
                </div>
                <Text strong>{formatPrice(item.totalPrice)}</Text>
              </div>
            ))}
            
            <Divider className="my-3" />
            
            <div className="flex justify-between items-center">
              <Title level={4} className="m-0">
                合計金額
              </Title>
              <Title level={4} className="m-0 text-primary">
                {formatPrice(order.totalAmount)}
              </Title>
            </div>
          </div>
        </Card>

        {/* Special Instructions */}
        {order.specialInstructions && (
          <Card title="特別な要望" size="small">
            <Paragraph className="m-0 text-sm">
              {order.specialInstructions}
            </Paragraph>
          </Card>
        )}

        {/* Important Notice */}
        <Card className="bg-blue-50 border-blue-200">
          <Space direction="vertical" size="small" className="w-full">
            <Text strong className="text-blue-800">
              <ClockCircleOutlined className="mr-2" />
              お知らせ
            </Text>
            <ul className="text-sm text-blue-700 m-0 pl-4">
              <li>調理が完了しましたら、店舗スタッフがお席までお持ちします</li>
              <li>混雑状況により、調理時間が前後する場合がございます</li>
              <li>ご不明な点がございましたら、お気軽にスタッフまでお声がけください</li>
            </ul>
          </Space>
        </Card>

        {/* Action Buttons */}
        <div className="space-y-3">
          <Button
            type="primary"
            size="large"
            onClick={onContinueOrdering}
            icon={<ShoppingOutlined />}
            className="w-full h-12"
          >
            追加で注文する
          </Button>
          
          <Space className="w-full">
            <Button
              type="default"
              onClick={onViewHistory}
              className="flex-1"
            >
              注文履歴を見る
            </Button>
            <Button
              type="default"
              onClick={onClose}
              className="flex-1"
            >
              閉じる
            </Button>
          </Space>
        </div>
      </div>

      <style jsx global>{`
        .order-success-modal .ant-modal-body {
          padding: 24px;
        }
        
        .order-success-modal .ant-result-title {
          color: #52c41a;
          font-size: 24px;
          font-weight: 600;
        }
        
        .order-success-modal .ant-result-subtitle {
          font-size: 16px;
        }
        
        @media (max-width: 768px) {
          .order-success-modal .ant-modal-content {
            margin: 16px;
          }
          
          .order-success-modal .ant-result-title {
            font-size: 20px;
          }
          
          .order-success-modal .ant-result-subtitle {
            font-size: 14px;
          }
          
          .order-success-modal .ant-btn-lg {
            height: 48px;
            font-size: 16px;
          }
        }
      `}</style>
    </Modal>
  );
};