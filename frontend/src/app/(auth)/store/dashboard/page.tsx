'use client';

import { Typography, Space, Row, Col } from 'antd';
import { 
  TableOutlined, 
  ShoppingCartOutlined, 
  DollarOutlined,
  ClockCircleOutlined 
} from '@ant-design/icons';

import { Card } from '@/components';
import { useAuth } from '@/hooks';

const { Title, Text } = Typography;

export default function StoreDashboardPage() {
  const { user } = useAuth();

  const stats = [
    {
      title: '総座席数',
      value: '24',
      icon: <TableOutlined style={{ fontSize: '24px', color: '#1890ff' }} />,
      color: '#1890ff',
    },
    {
      title: '本日の注文数',
      value: '67',
      icon: <ShoppingCartOutlined style={{ fontSize: '24px', color: '#52c41a' }} />,
      color: '#52c41a',
    },
    {
      title: '本日の売上',
      value: '¥89,500',
      icon: <DollarOutlined style={{ fontSize: '24px', color: '#faad14' }} />,
      color: '#faad14',
    },
    {
      title: '平均注文時間',
      value: '12分',
      icon: <ClockCircleOutlined style={{ fontSize: '24px', color: '#f5222d' }} />,
      color: '#f5222d',
    },
  ];

  return (
    <div className="p-6">
      <Space direction="vertical" size="large" className="w-full">
        <div>
          <Title level={2}>店舗ダッシュボード</Title>
          <Text type="secondary">
            ようこそ、{user?.name}さん。本日の店舗状況をご確認ください。
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
              <Title level={4}>最近の注文</Title>
              <Space direction="vertical" className="w-full">
                <div className="flex justify-between items-center py-2 border-b">
                  <div>
                    <Text strong>テーブル 5</Text>
                    <br />
                    <Text type="secondary">ハンバーガーセット x2</Text>
                  </div>
                  <div className="text-right">
                    <Text strong>¥2,400</Text>
                    <br />
                    <Text type="secondary">5分前</Text>
                  </div>
                </div>
                <div className="flex justify-between items-center py-2 border-b">
                  <div>
                    <Text strong>テーブル 12</Text>
                    <br />
                    <Text type="secondary">パスタセット x1</Text>
                  </div>
                  <div className="text-right">
                    <Text strong>¥1,800</Text>
                    <br />
                    <Text type="secondary">8分前</Text>
                  </div>
                </div>
                <div className="flex justify-between items-center py-2 border-b">
                  <div>
                    <Text strong>テーブル 3</Text>
                    <br />
                    <Text type="secondary">ピザセット x3</Text>
                  </div>
                  <div className="text-right">
                    <Text strong>¥4,200</Text>
                    <br />
                    <Text type="secondary">12分前</Text>
                  </div>
                </div>
              </Space>
            </Card>
          </Col>

          <Col xs={24} lg={12}>
            <Card>
              <Title level={4}>座席状況</Title>
              <Space direction="vertical" className="w-full">
                <div className="flex justify-between items-center py-2">
                  <Text>利用中の座席</Text>
                  <Text style={{ color: '#f5222d' }}>18 / 24</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>空席</Text>
                  <Text style={{ color: '#52c41a' }}>6</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>清掃中</Text>
                  <Text style={{ color: '#faad14' }}>0</Text>
                </div>
                <div className="flex justify-between items-center py-2">
                  <Text>稼働率</Text>
                  <Text strong>75%</Text>
                </div>
              </Space>
            </Card>
          </Col>
        </Row>
      </Space>
    </div>
  );
}