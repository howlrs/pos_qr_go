'use client';

import React from 'react';
import { FloatButton, Badge } from 'antd';
import { ShoppingCartOutlined } from '@ant-design/icons';
import { FloatingCartProps } from './FloatingCart.types';

export const FloatingCart: React.FC<FloatingCartProps> = ({
  itemCount = 0,
  totalAmount = 0,
  onClick,
  visible = true,
}) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price);
  };

  if (!visible || itemCount === 0) {
    return null;
  }

  return (
    <div className="floating-cart">
      <Badge count={itemCount} size="small" offset={[-8, 8]}>
        <FloatButton
          icon={<ShoppingCartOutlined />}
          type="primary"
          onClick={onClick}
          className="floating-cart-button"
          style={{
            width: 64,
            height: 64,
            fontSize: 24,
          }}
          description={
            <div className="text-xs text-white mt-1">
              {formatPrice(totalAmount)}
            </div>
          }
        />
      </Badge>

      <style jsx global>{`
        .floating-cart {
          position: fixed;
          bottom: 24px;
          right: 24px;
          z-index: 1000;
        }
        
        .floating-cart-button {
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
        
        .floating-cart-button:hover {
          transform: scale(1.05);
          transition: transform 0.2s ease;
        }
        
        @media (max-width: 768px) {
          .floating-cart {
            bottom: 16px;
            right: 16px;
          }
          
          .floating-cart-button {
            width: 56px;
            height: 56px;
            font-size: 20px;
          }
        }
        
        @media (max-width: 480px) {
          .floating-cart {
            bottom: 12px;
            right: 12px;
          }
          
          .floating-cart-button {
            width: 48px;
            height: 48px;
            font-size: 18px;
          }
        }
      `}</style>
    </div>
  );
};