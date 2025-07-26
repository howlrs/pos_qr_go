'use client';

import React, { useState, useMemo } from 'react';
import { useParams } from 'next/navigation';
import {
  Row,
  Col,
  Typography,
  Tabs,
  Input,
  Space,
  Alert,
  Spin,
  message,
  Modal,
  Drawer,
} from 'antd';
import {
  SearchOutlined,
  ShoppingCartOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import { OrderLayout } from '@/components/layouts';
import {
  MenuCard,
  CartDrawer,
  OrderHistory,
  OrderConfirmation,
  OrderSuccess,
  FloatingCart,
  OrderStatus,
  StoreInfo,
  PageLoading,
  ErrorBoundary,
} from '@/components/ui';
import {
  useOrderPage,
  useAddToCart,
  useOrderHistory,
  usePlaceOrder,
} from '@/hooks/api/useOrders';
import type { MenuItem, MenuCategory, Order } from '@/types/models';

const { Title, Text } = Typography;
const { TabPane } = Tabs;

interface OrderPageProps {}

const OrderPage: React.FC<OrderPageProps> = () => {
  const params = useParams();
  const sessionId = params.sessionId as string;

  // State management
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const [showCartDrawer, setShowCartDrawer] = useState(false);
  const [showHistoryModal, setShowHistoryModal] = useState(false);
  const [showOrderConfirmation, setShowOrderConfirmation] = useState(false);
  const [showOrderSuccess, setShowOrderSuccess] = useState(false);
  const [showStoreInfo, setShowStoreInfo] = useState(false);
  const [specialInstructions, setSpecialInstructions] = useState('');
  const [completedOrder, setCompletedOrder] = useState<Order | null>(null);
  const [estimatedWaitTime, setEstimatedWaitTime] = useState<number | null>(null);
  const [activeOrder, setActiveOrder] = useState<Order | null>(null);

  // API hooks
  const { session, menu, cart, isLoading, isError, error } = useOrderPage(sessionId);
  const addToCartMutation = useAddToCart(sessionId);
  const placeOrderMutation = usePlaceOrder(sessionId);
  const { data: historyData, isLoading: historyLoading } = useOrderHistory(sessionId);

  // Memoized filtered menu items
  const filteredMenuItems = useMemo(() => {
    if (!menu.data?.items) return [];

    let items = menu.data.items;

    // Filter by category
    if (selectedCategory !== 'all') {
      items = items.filter((item) => item.categoryId === selectedCategory);
    }

    // Filter by search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      items = items.filter(
        (item) =>
          item.name.toLowerCase().includes(query) ||
          item.description?.toLowerCase().includes(query)
      );
    }

    return items;
  }, [menu.data?.items, selectedCategory, searchQuery]);

  // Event handlers
  const handleQuantityChange = (menuItemId: string, quantity: number) => {
    setQuantities((prev) => ({
      ...prev,
      [menuItemId]: quantity,
    }));
  };

  const handleAddToCart = async (menuItemId: string, quantity: number) => {
    try {
      await addToCartMutation.mutateAsync({
        menuItemId,
        quantity,
      });
      
      // Reset quantity after successful add
      setQuantities((prev) => ({
        ...prev,
        [menuItemId]: 0,
      }));
      
      message.success('カートに追加しました');
    } catch (error) {
      message.error('カートへの追加に失敗しました');
    }
  };

  const handleOrderConfirm = () => {
    setShowCartDrawer(false);
    setShowOrderConfirmation(true);
  };

  const handleConfirmOrder = async () => {
    try {
      const result = await placeOrderMutation.mutateAsync({
        specialInstructions: specialInstructions.trim() || undefined,
      });
      
      setCompletedOrder(result.order);
      setEstimatedWaitTime(result.estimatedWaitTime);
      setShowOrderConfirmation(false);
      setShowOrderSuccess(true);
      setSpecialInstructions('');
    } catch (error) {
      message.error('注文の送信に失敗しました');
    }
  };

  const handleContinueOrdering = () => {
    setShowOrderSuccess(false);
    setCompletedOrder(null);
    setEstimatedWaitTime(null);
  };

  const handleViewHistory = () => {
    setShowOrderSuccess(false);
    setShowHistoryModal(true);
  };

  const handleStoreInfoClick = () => {
    setShowStoreInfo(true);
  };

  const refreshOrderStatus = () => {
    // This would typically refetch the active order status
    // For now, we'll just simulate a refresh
    console.log('Refreshing order status...');
  };

  // Check for active orders
  React.useEffect(() => {
    if (historyData?.orders && historyData.orders.length > 0) {
      const latestOrder = historyData.orders[0];
      if (['pending', 'confirmed', 'preparing', 'ready'].includes(latestOrder.status)) {
        setActiveOrder(latestOrder);
      } else {
        setActiveOrder(null);
      }
    }
  }, [historyData]);

  // Loading state
  if (isLoading) {
    return <PageLoading />;
  }

  // Error state
  if (isError) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Alert
          message="エラーが発生しました"
          description={error?.message || 'セッションの読み込みに失敗しました'}
          type="error"
          showIcon
        />
      </div>
    );
  }

  // Session validation
  if (!session.data?.session) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Alert
          message="無効なセッション"
          description="このQRコードは無効または期限切れです"
          type="warning"
          showIcon
        />
      </div>
    );
  }

  const sessionData = session.data.session;
  const menuData = menu.data;
  const cartData = cart.data?.cart;

  return (
    <ErrorBoundary>
      <OrderLayout
        sessionId={sessionId}
        session={sessionData}
        onCartClick={() => setShowCartDrawer(true)}
        onHistoryClick={() => setShowHistoryModal(true)}
        onStoreInfoClick={handleStoreInfoClick}
        showCartDrawer={showCartDrawer}
        onCartDrawerClose={() => setShowCartDrawer(false)}
        cartContent={
          <CartDrawer
            sessionId={sessionId}
            cart={cartData}
            onOrderConfirm={handleOrderConfirm}
          />
        }
      >
        <div className="p-4 max-w-7xl mx-auto">
          {/* Active Order Status */}
          {activeOrder && (
            <div className="mb-6">
              <OrderStatus
                order={activeOrder}
                onRefresh={refreshOrderStatus}
                refreshing={false}
              />
            </div>
          )}

          {/* Welcome Section */}
          <div className="mb-6 text-center">
            <Title level={2} className="mb-2">
              {sessionData.store.name}へようこそ
            </Title>
            <Text type="secondary" className="text-lg">
              座席 {sessionData.seat.name || sessionData.seat.number} からご注文いただけます
            </Text>
          </div>

          {/* Search and Filter */}
          <div className="mb-6">
            <Space direction="vertical" className="w-full">
              <Input
                placeholder="メニューを検索..."
                prefix={<SearchOutlined />}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                size="large"
                className="max-w-md"
              />
            </Space>
          </div>

          {/* Menu Categories */}
          {menuData?.categories && menuData.categories.length > 0 && (
            <Tabs
              activeKey={selectedCategory}
              onChange={setSelectedCategory}
              className="mb-6"
              size="large"
            >
              <TabPane tab="すべて" key="all" />
              {menuData.categories
                .filter((category) => category.isActive)
                .sort((a, b) => a.displayOrder - b.displayOrder)
                .map((category) => (
                  <TabPane tab={category.name} key={category.id} />
                ))}
            </Tabs>
          )}

          {/* Menu Items */}
          <div className="menu-grid">
            {filteredMenuItems.length > 0 ? (
              <Row gutter={[16, 16]}>
                {filteredMenuItems.map((item) => (
                  <Col
                    key={item.id}
                    xs={24}
                    sm={12}
                    md={8}
                    lg={6}
                    xl={6}
                  >
                    <MenuCard
                      menuItem={item}
                      quantity={quantities[item.id] || 0}
                      onQuantityChange={(quantity) =>
                        handleQuantityChange(item.id, quantity)
                      }
                      onAddToCart={(quantity) =>
                        handleAddToCart(item.id, quantity)
                      }
                      loading={addToCartMutation.isPending}
                    />
                  </Col>
                ))}
              </Row>
            ) : (
              <div className="text-center py-12">
                <Text type="secondary" className="text-lg">
                  {searchQuery || selectedCategory !== 'all'
                    ? '条件に一致するメニューが見つかりません'
                    : 'メニューがありません'}
                </Text>
              </div>
            )}
          </div>
        </div>

        {/* Order Confirmation Modal */}
        <OrderConfirmation
          visible={showOrderConfirmation}
          cart={cartData}
          specialInstructions={specialInstructions}
          onSpecialInstructionsChange={setSpecialInstructions}
          onConfirm={handleConfirmOrder}
          onCancel={() => setShowOrderConfirmation(false)}
          loading={placeOrderMutation.isPending}
        />

        {/* Order Success Modal */}
        <OrderSuccess
          visible={showOrderSuccess}
          order={completedOrder || undefined}
          estimatedWaitTime={estimatedWaitTime || undefined}
          onClose={() => setShowOrderSuccess(false)}
          onContinueOrdering={handleContinueOrdering}
          onViewHistory={handleViewHistory}
        />

        {/* Store Info Modal */}
        <StoreInfo
          visible={showStoreInfo}
          store={sessionData.store}
          onClose={() => setShowStoreInfo(false)}
        />

        {/* Order History Modal */}
        <Modal
          title="注文履歴"
          open={showHistoryModal}
          onCancel={() => setShowHistoryModal(false)}
          footer={null}
          width={600}
          className="order-history-modal"
        >
          <OrderHistory
            orders={historyData?.orders}
            loading={historyLoading}
          />
        </Modal>

        {/* Floating Cart Button (Mobile) */}
        <FloatingCart
          itemCount={cartData?.totalItems}
          totalAmount={cartData?.totalAmount}
          onClick={() => setShowCartDrawer(true)}
          visible={!showCartDrawer && (cartData?.totalItems || 0) > 0}
        />
      </OrderLayout>

      <style jsx global>{`
        .menu-grid .ant-col {
          display: flex;
        }
        
        .menu-grid .menu-card {
          height: 100%;
          display: flex;
          flex-direction: column;
        }
        
        .menu-grid .ant-card-body {
          flex: 1;
          display: flex;
          flex-direction: column;
        }
        
        .order-history-modal .ant-modal-body {
          max-height: 60vh;
          overflow-y: auto;
        }
        
        @media (max-width: 768px) {
          .ant-tabs-tab {
            font-size: 14px;
            padding: 8px 12px;
          }
          
          .menu-grid .ant-col {
            margin-bottom: 16px;
          }
        }
        
        @media (max-width: 480px) {
          .ant-input-affix-wrapper-lg {
            font-size: 16px;
          }
          
          .ant-typography h2 {
            font-size: 24px;
          }
        }
      `}</style>
    </ErrorBoundary>
  );
};

export default OrderPage;