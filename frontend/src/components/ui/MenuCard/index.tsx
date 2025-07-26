'use client';

import React from 'react';
import { Card, Typography, Button, InputNumber, Space, Tag, Image } from 'antd';
import { PlusOutlined, MinusOutlined } from '@ant-design/icons';
import { MenuCardProps } from './MenuCard.types';

const { Title, Text, Paragraph } = Typography;

export const MenuCard: React.FC<MenuCardProps> = ({
  menuItem,
  quantity = 0,
  onQuantityChange,
  onAddToCart,
  loading = false,
  disabled = false,
}) => {
  const handleQuantityChange = (value: number | null) => {
    const newQuantity = Math.max(0, value || 0);
    onQuantityChange?.(newQuantity);
  };

  const handleAddToCart = () => {
    if (quantity > 0) {
      onAddToCart?.(quantity);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price);
  };

  return (
    <Card
      className={`menu-card ${!menuItem.isAvailable ? 'unavailable' : ''}`}
      cover={
        menuItem.imageUrl ? (
          <div className="relative h-48 overflow-hidden">
            <Image
              alt={menuItem.name}
              src={menuItem.imageUrl}
              className="w-full h-full object-cover"
              preview={false}
            />
            {!menuItem.isAvailable && (
              <div className="absolute inset-0 bg-gray-500 bg-opacity-50 flex items-center justify-center">
                <Tag color="red" className="text-lg">
                  売り切れ
                </Tag>
              </div>
            )}
          </div>
        ) : (
          <div className="h-48 bg-gray-100 flex items-center justify-center">
            <Text type="secondary">画像なし</Text>
          </div>
        )
      }
      actions={
        menuItem.isAvailable
          ? [
              <Space key="quantity" className="w-full justify-center">
                <Button
                  icon={<MinusOutlined />}
                  size="small"
                  onClick={() => handleQuantityChange(quantity - 1)}
                  disabled={quantity <= 0 || disabled}
                />
                <InputNumber
                  min={0}
                  max={99}
                  value={quantity}
                  onChange={handleQuantityChange}
                  size="small"
                  className="w-16 text-center"
                  disabled={disabled}
                />
                <Button
                  icon={<PlusOutlined />}
                  size="small"
                  onClick={() => handleQuantityChange(quantity + 1)}
                  disabled={disabled}
                />
              </Space>,
              <Button
                key="add"
                type="primary"
                onClick={handleAddToCart}
                loading={loading}
                disabled={quantity <= 0 || disabled}
                className="w-full"
              >
                カートに追加 ({formatPrice(menuItem.price * quantity)})
              </Button>,
            ]
          : [
              <Button key="unavailable" disabled className="w-full">
                売り切れ
              </Button>,
            ]
      }
    >
      <Card.Meta
        title={
          <div className="flex justify-between items-start">
            <Title level={5} className="m-0 flex-1">
              {menuItem.name}
            </Title>
            <Text strong className="text-lg text-primary ml-2">
              {formatPrice(menuItem.price)}
            </Text>
          </div>
        }
        description={
          <div className="space-y-2">
            {menuItem.description && (
              <Paragraph className="text-gray-600 text-sm m-0">
                {menuItem.description}
              </Paragraph>
            )}

            {/* Allergens */}
            {menuItem.allergens && menuItem.allergens.length > 0 && (
              <div>
                <Text type="secondary" className="text-xs">
                  アレルゲン:
                </Text>
                <div className="mt-1">
                  {menuItem.allergens.map((allergen) => (
                    <Tag key={allergen} color="orange" className="text-xs">
                      {allergen}
                    </Tag>
                  ))}
                </div>
              </div>
            )}

            {/* Nutrition Info */}
            {menuItem.nutritionInfo && (
              <div className="text-xs text-gray-500">
                {menuItem.nutritionInfo.calories && (
                  <span>カロリー: {menuItem.nutritionInfo.calories}kcal</span>
                )}
              </div>
            )}
          </div>
        }
      />

      <style jsx global>{`
        .menu-card.unavailable {
          opacity: 0.6;
        }
        
        .menu-card .ant-card-actions {
          background: #fafafa;
        }
        
        .menu-card .ant-card-actions > li {
          margin: 8px 0;
        }
        
        @media (max-width: 768px) {
          .menu-card .ant-card-cover {
            height: 200px;
          }
          
          .menu-card .ant-typography h5 {
            font-size: 16px;
          }
          
          .menu-card .ant-card-actions .ant-btn {
            font-size: 12px;
            padding: 4px 8px;
          }
        }
      `}</style>
    </Card>
  );
};