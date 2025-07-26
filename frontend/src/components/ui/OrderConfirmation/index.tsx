'use client';

import React from 'react';
import {
  Modal,
  Typography,
  List,
  Divider,
  Space,
  Button,
  Input,
  Alert,
  Tag,
} from 'antd';
import {
  CheckCircleOutlined,
  ClockCircleOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons';
import { OrderConfirmationProps } from './OrderConfirmation.types';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

export const OrderConfirmation: React.FC<OrderConfirmationProps> = ({
  visible,
  cart,
  specialInstructions,
  onSpecialInstructionsChange,
  onConfirm,
  onCancel,
  loading = false,
}) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price);
  };

  if (!cart || cart.items.length === 0) {
    return null;
  }

  return (
    <Modal
      title={
        <Space>
          <CheckCircleOutlined className="text-green-500" />
          <span>注文内容の確認</span>
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      width={600}
      footer={[
        <Button key="cancel" onClick={onCancel} disabled={loading}>
          キャンセル
        </Button>,
        <Button
          key="confirm"
          type="primary"
          onClick={onConfirm}
          loading={loading}
          size="large"
          className="h-12"
        >
          注文を確定する ({formatPrice(cart.totalAmount)})
        </Button>,
      ]}
      className="order-confirmation-modal"
    >
      <div className="space-y-6">
        {/* Order Summary Alert */}
        <Alert
          message="注文内容をご確認ください"
          description="注文確定後は内容の変更ができませんので、よくご確認ください。"
          type="info"
          icon={<InfoCircleOutlined />}
          showIcon
        />

        {/* Order Items */}
        <div>
          <Title level={4} className="mb-3">
            注文商品
          </Title>
          <List
            dataSource={cart.items}
            renderItem={(item) => (
              <List.Item className="px-0">
                <div className="w-full">
                  <div className="flex justify-between items-start mb-2">
                    <div className="flex-1">
                      <Text strong className="text-base">
                        {item.menuItem.name}
                      </Text>
                      <div className="mt-1">
                        <Text type="secondary" className="text-sm">
                          {formatPrice(item.unitPrice)} × {item.quantity}
                        </Text>
                      </div>
                      {item.menuItem.allergens && item.menuItem.allergens.length > 0 && (
                        <div className="mt-1">
                          <Text type="secondary" className="text-xs">
                            アレルゲン:
                          </Text>
                          {item.menuItem.allergens.map((allergen) => (
                            <Tag key={allergen} color="orange" className="ml-1 text-xs">
                              {allergen}
                            </Tag>
                          ))}
                        </div>
                      )}
                    </div>
                    <Text strong className="text-lg">
                      {formatPrice(item.totalPrice)}
                    </Text>
                  </div>
                  {item.specialInstructions && (
                    <div className="mt-2 p-2 bg-gray-50 rounded">
                      <Text type="secondary" className="text-sm">
                        <strong>備考:</strong> {item.specialInstructions}
                      </Text>
                    </div>
                  )}
                </div>
              </List.Item>
            )}
          />
        </div>

        <Divider />

        {/* Special Instructions */}
        <div>
          <Title level={4} className="mb-3">
            特別な要望（任意）
          </Title>
          <TextArea
            placeholder="アレルギー対応、調理方法の希望、その他のご要望をお書きください..."
            value={specialInstructions}
            onChange={(e) => onSpecialInstructionsChange?.(e.target.value)}
            rows={4}
            maxLength={500}
            showCount
            disabled={loading}
          />
          <Text type="secondary" className="text-xs mt-1 block">
            ※ 内容によってはご対応できない場合がございます
          </Text>
        </div>

        <Divider />

        {/* Order Total */}
        <div className="bg-gray-50 p-4 rounded-lg">
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <Text>商品数</Text>
              <Text>{cart.totalItems}点</Text>
            </div>
            <div className="flex justify-between items-center">
              <Title level={3} className="m-0">
                合計金額
              </Title>
              <Title level={3} className="m-0 text-primary">
                {formatPrice(cart.totalAmount)}
              </Title>
            </div>
          </div>
        </div>

        {/* Estimated Time */}
        <Alert
          message={
            <Space>
              <ClockCircleOutlined />
              <span>調理時間の目安: 15-25分</span>
            </Space>
          }
          type="warning"
          showIcon={false}
          className="text-center"
        />

        {/* Terms */}
        <div className="text-center">
          <Paragraph type="secondary" className="text-xs mb-0">
            注文確定により、当店の利用規約に同意したものとみなします。
            <br />
            混雑状況により調理時間が前後する場合がございます。
          </Paragraph>
        </div>
      </div>

      <style jsx global>{`
        .order-confirmation-modal .ant-modal-body {
          max-height: 70vh;
          overflow-y: auto;
        }
        
        .order-confirmation-modal .ant-list-item {
          border-bottom: 1px solid #f0f0f0;
          padding: 16px 0;
        }
        
        .order-confirmation-modal .ant-list-item:last-child {
          border-bottom: none;
        }
        
        @media (max-width: 768px) {
          .order-confirmation-modal .ant-modal-content {
            margin: 16px;
          }
          
          .order-confirmation-modal .ant-typography h4 {
            font-size: 16px;
          }
          
          .order-confirmation-modal .ant-btn-lg {
            height: 48px;
            font-size: 16px;
          }
        }
      `}</style>
    </Modal>
  );
};