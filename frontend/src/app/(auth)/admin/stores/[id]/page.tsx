'use client';

import { useRouter } from 'next/navigation';
import {
  Typography,
  Space,
  Row,
  Col,
  Tag,
  Descriptions,
  Statistic,
  Button,
  Spin,
} from 'antd';
import {
  ArrowLeftOutlined,
  EditOutlined,
  ShopOutlined,
  PhoneOutlined,
  MailOutlined,
  EnvironmentOutlined,
  SettingOutlined,
} from '@ant-design/icons';

import { Card } from '@/components';
import { useStore, useStoreStats } from '@/hooks/api/useStores';
import { usePermissions } from '@/hooks';
import { PERMISSIONS } from '@/lib/auth/permissions';

const { Title, Text } = Typography;

interface StoreDetailPageProps {
  params: {
    id: string;
  };
}

export default function StoreDetailPage({ params }: StoreDetailPageProps) {
  const router = useRouter();
  const { hasPermission } = usePermissions();
  const storeId = params.id;

  // API hooks
  const {
    data: store,
    isLoading: storeLoading,
    error: storeError,
  } = useStore(storeId);

  const {
    data: stats,
    isLoading: statsLoading,
  } = useStoreStats(storeId);

  // Permissions
  const canManageStores = hasPermission(PERMISSIONS.ADMIN.MANAGE_STORES);
  const canViewAnalytics = hasPermission(PERMISSIONS.ADMIN.VIEW_ANALYTICS);

  // Handlers
  const handleBack = () => {
    router.push('/admin/stores');
  };

  const handleEdit = () => {
    router.push(`/admin/stores/${storeId}/edit`);
  };

  if (storeLoading) {
    return (
      <div className="p-6 flex justify-center items-center min-h-96">
        <Spin size="large" />
      </div>
    );
  }

  if (storeError || !store) {
    return (
      <div className="p-6">
        <Card>
          <div className="text-center py-8">
            <Text type="danger">店舗データの読み込みに失敗しました</Text>
          </div>
        </Card>
      </div>
    );
  }

  const getFeatureLabel = (feature: string) => {
    const labels: Record<string, string> = {
      qr_ordering: 'QR注文',
      table_service: 'テーブルサービス',
      takeaway: 'テイクアウト',
      delivery: 'デリバリー',
      payment_integration: '決済連携',
    };
    return labels[feature] || feature;
  };

  return (
    <div className="p-6">
      <Space direction="vertical" size="large" className="w-full">
        {/* Header */}
        <div>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={handleBack}
            className="mb-4"
          >
            店舗一覧に戻る
          </Button>
          <Row justify="space-between" align="middle">
            <Col>
              <Space>
                <ShopOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
                <div>
                  <Title level={2} className="!mb-0">
                    {store.name}
                  </Title>
                  <Text type="secondary">{store.description}</Text>
                </div>
              </Space>
            </Col>
            <Col>
              <Space>
                <Tag color={store.isActive ? 'green' : 'red'} className="text-sm px-3 py-1">
                  {store.isActive ? '有効' : '無効'}
                </Tag>
                {canManageStores && (
                  <Button
                    type="primary"
                    icon={<EditOutlined />}
                    onClick={handleEdit}
                  >
                    編集
                  </Button>
                )}
              </Space>
            </Col>
          </Row>
        </div>

        {/* Statistics */}
        {canViewAnalytics && stats && (
          <Card title="統計情報" loading={statsLoading}>
            <Row gutter={16}>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="総注文数"
                  value={stats.totalOrders}
                  suffix="件"
                />
              </Col>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="総売上"
                  value={stats.totalRevenue}
                  prefix="¥"
                  precision={0}
                />
              </Col>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="本日の注文"
                  value={stats.ordersToday}
                  suffix="件"
                />
              </Col>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="本日の売上"
                  value={stats.revenueToday}
                  prefix="¥"
                  precision={0}
                />
              </Col>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="稼働座席"
                  value={stats.activeSeats}
                  suffix={`/ ${store.settings.maxSeats}`}
                />
              </Col>
              <Col xs={12} sm={8} md={4}>
                <Statistic
                  title="平均注文額"
                  value={stats.averageOrderValue}
                  prefix="¥"
                  precision={0}
                />
              </Col>
            </Row>
          </Card>
        )}

        <Row gutter={16}>
          {/* Basic Information */}
          <Col xs={24} lg={12}>
            <Card title="基本情報" extra={<ShopOutlined />}>
              <Descriptions column={1} size="small">
                <Descriptions.Item
                  label={<Space><EnvironmentOutlined />住所</Space>}
                >
                  {store.address}
                </Descriptions.Item>
                <Descriptions.Item
                  label={<Space><PhoneOutlined />電話番号</Space>}
                >
                  {store.phone}
                </Descriptions.Item>
                <Descriptions.Item
                  label={<Space><MailOutlined />メールアドレス</Space>}
                >
                  {store.email}
                </Descriptions.Item>
                <Descriptions.Item label="作成日">
                  {new Date(store.createdAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
                <Descriptions.Item label="更新日">
                  {new Date(store.updatedAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>

          {/* Settings */}
          <Col xs={24} lg={12}>
            <Card title="店舗設定" extra={<SettingOutlined />}>
              <Descriptions column={1} size="small">
                <Descriptions.Item label="タイムゾーン">
                  {store.settings.timezone}
                </Descriptions.Item>
                <Descriptions.Item label="通貨">
                  {store.settings.currency}
                </Descriptions.Item>
                <Descriptions.Item label="言語">
                  {store.settings.language === 'ja' ? '日本語' : store.settings.language}
                </Descriptions.Item>
                <Descriptions.Item label="注文タイムアウト">
                  {store.settings.orderTimeout}分
                </Descriptions.Item>
                <Descriptions.Item label="最大座席数">
                  {store.settings.maxSeats}席
                </Descriptions.Item>
                <Descriptions.Item label="利用可能機能">
                  <Space wrap>
                    {store.settings.features.map((feature) => (
                      <Tag key={feature} color="blue">
                        {getFeatureLabel(feature)}
                      </Tag>
                    ))}
                  </Space>
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>
        </Row>
      </Space>
    </div>
  );
}