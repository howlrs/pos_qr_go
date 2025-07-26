'use client';

import React from 'react';
import { Layout, Typography, Badge, Button, Space, Drawer } from 'antd';
import { ShoppingCartOutlined, HistoryOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { OrderLayoutProps } from './OrderLayout.types';
import { useCart } from '@/hooks/api/useOrders';

const { Header, Content } = Layout;
const { Title, Text } = Typography;

export const OrderLayout: React.FC<OrderLayoutProps> = ({
  children,
  sessionId,
  session,
  onCartClick,
  onHistoryClick,
  onStoreInfoClick,
  showCartDrawer = false,
  onCartDrawerClose,
  cartContent,
}) => {
  const { data: cartData } = useCart(sessionId);
  const cartItemCount = cartData?.cart?.totalItems || 0;

  return (
    <Layout className="min-h-screen bg-gray-50">
      {/* Header */}
      <Header className="bg-white shadow-sm px-4 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div>
            <Title level={4} className="m-0 text-gray-800">
              {session?.store?.name || 'レストラン'}
            </Title>
            <Text type="secondary" className="text-sm">
              座席: {session?.seat?.name || session?.seat?.number}
            </Text>
          </div>
        </div>

        <Space size="middle">
          {/* Order History Button */}
          <Button
            type="text"
            icon={<HistoryOutlined />}
            onClick={onHistoryClick}
            className="flex items-center"
          >
            履歴
          </Button>

          {/* Store Info Button */}
          <Button
            type="text"
            icon={<InfoCircleOutlined />}
            onClick={onStoreInfoClick}
            className="flex items-center"
          >
            店舗情報
          </Button>

          {/* Cart Button */}
          <Badge count={cartItemCount} size="small" offset={[-5, 5]}>
            <Button
              type="primary"
              icon={<ShoppingCartOutlined />}
              onClick={onCartClick}
              className="flex items-center"
              size="large"
            >
              カート
            </Button>
          </Badge>
        </Space>
      </Header>

      {/* Main Content */}
      <Content className="flex-1">
        {children}
      </Content>

      {/* Cart Drawer */}
      <Drawer
        title="カート"
        placement="right"
        onClose={onCartDrawerClose}
        open={showCartDrawer}
        width={400}
        className="cart-drawer"
      >
        {cartContent}
      </Drawer>

      {/* Mobile-optimized styles */}
      <style jsx global>{`
        @media (max-width: 768px) {
          .ant-layout-header {
            padding: 0 16px;
            height: 64px;
          }
          
          .cart-drawer .ant-drawer-content-wrapper {
            width: 100vw !important;
          }
          
          .ant-typography h4 {
            font-size: 16px;
          }
          
          .ant-btn {
            font-size: 14px;
          }
        }
        
        @media (max-width: 480px) {
          .ant-layout-header {
            padding: 0 12px;
          }
          
          .ant-space-item .ant-btn span {
            display: none;
          }
          
          .ant-space-item .ant-btn {
            padding: 4px 8px;
          }
        }
      `}</style>
    </Layout>
  );
};