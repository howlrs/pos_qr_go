'use client';

import { Typography, Space, Row, Col } from 'antd';
import { 
  ShopOutlined, 
  UserOutlined, 
  BarChartOutlined,
  DollarOutlined 
} from '@ant-design/icons';

import { Card } from '@/components';
import { useAuth } from '@/hooks';

const { Title, Text } = Typography;

export default function AdminDashboardPage() {
  const { user } = useAuth();

  const stats = [
    {
      title: '総店舗数',
      value: '12',
      icon: <ShopOutlined style={{ fontSize: '24px', color: '#1890ff' }} />,
      color: '#1890ff',
    },
    {
      title: '総管理者数',
      value: '8',
      icon: <UserOutlined style={{ fontSize: '24px', color: '#52c41a' }} />,
      color: '#52c41a',
    },
    {
      title: '月間売上',
      value: '¥2,450,000',
      icon: <DollarOutlined style={{ fontSize: '24px', color: '#faad14' }} />,
      color: '#faad14',
    },
    {
      title: '月間注文数',
      value: '1,234',
      icon: <BarChartOutlined style={{ fontSize: '24px', color: '#f5222d' }} />,
      color: '#f5222d',
    },
  ];

  return (
    <div className="p-6">
      <Space direction="vertical" size="large" className="w-full">
        <div>
          <Title level={2}>管理者ダッシュボード</Title>
          <Text type="secondary">
            ようこそ、{user?.name}さん。システム全体の状況をご確認ください。
          </Text>
        </div>

        <Row gutter={[16, 16]}>
          {stats.map((stat, index) => (
            <Col xs={24} sm={12} lg={6} key={index}>
              <Card className="text-center">
                <Space direction="vertical" size="small" className="w-full">
                  {stat.icon}
                  <Title level={3} className="!mb-0" style={{ color: stat.color }}>
                    {stat.value}
                  </Title>
                  <Text type="secondary">{stat.title}</Text>
                </Space>
              </Card>
            </Col>
          ))}
        </Row>

        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <Card>
              <Title level={4}>最近の活動</Title>
              <Space direction="vertical" className="w-full">
                <div className="flex justify-between items-center py-2 border-b">
                  <Text>新しい店舗が登録されました</Text>
                  <Text type="secondary">2時間前</Text>
                </div>
                <div className="flex justify-between items-center py-2 border-b">
                  <Text>管理者権限が更新されました</Text>
                  <Text type="secondary">4時間前</Text>
                </div>
                <div className="flex justify-between items-center py-2 border-b">
                  <Text>システムメンテナンスが完了しました</Text>
                  <Text type="secondary">1日前</Text>
                </div>
              </Space>
            </Card>
          </Col>

          <Col xs={24} lg={12}>
            <Card>
              <Title level={4}>システム状況</Title>
              <Space direction="vertical" className="w-full">
                <div className="flex justify-between items-center py-2">
                  <Text>サーバー状態</Text>
                  <Text style={{ color: '#52c41a' }}>正常</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>データベース</Text>
                  <Text style={{ color: '#52c41a' }}>正常</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>API応答時間</Text>
                  <Text style={{ color: '#52c41a' }}>125ms</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>アクティブセッション</Text>
                  <Text>47</Text>
                </div>
              </Space>
            </Card>
          </Col>
        </Row>
      </Space>
    </div>
  );
}