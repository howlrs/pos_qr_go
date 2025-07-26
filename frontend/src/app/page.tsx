'use client';

import { Button, Card, Typography, Space, message } from 'antd';
import { CheckCircleOutlined, ApiOutlined } from '@ant-design/icons';

import { env } from '@/lib/config/env';
import { api, API_ENDPOINTS } from '@/lib/api';

const { Title, Paragraph, Text } = Typography;

export default function Home() {
  const testApiConnection = async () => {
    try {
      message.loading('API接続テスト中...', 0);
      
      // Test API connection
      const response = await api.get(API_ENDPOINTS.COMMON.HEALTH);
      
      message.destroy();
      message.success('API接続成功！');
      // eslint-disable-next-line no-console
      console.log('API Response:', response.data);
    } catch (error) {
      message.destroy();
      message.error('API接続に失敗しました');
      // eslint-disable-next-line no-console
      console.error('API Error:', error);
    }
  };

  return (
    <div className='min-h-screen bg-gray-50 flex items-center justify-center p-4'>
      <Card className='max-w-lg w-full text-center'>
        <Space direction='vertical' size='large' className='w-full'>
          <CheckCircleOutlined style={{ fontSize: '48px', color: '#52c41a' }} />

          <Title level={2}>POS QR System</Title>

          <Paragraph>
            フロントエンド基盤構築が完了しました。
            <br />
            環境変数とAPI設定が正常に動作しています。
          </Paragraph>

          <Card size='small' className='text-left'>
            <Text strong>環境設定:</Text>
            <br />
            <Text code>API_URL: {env.API_URL}</Text>
            <br />
            <Text code>APP_NAME: {env.APP_NAME}</Text>
            <br />
            <Text code>ENVIRONMENT: {env.ENVIRONMENT}</Text>
            <br />
            <Text code>DEBUG: {env.DEBUG.toString()}</Text>
          </Card>

          <Space direction='vertical' className='w-full'>
            <Button
              type='primary'
              size='large'
              icon={<ApiOutlined />}
              onClick={testApiConnection}
              block
            >
              API接続テスト
            </Button>

            <Space>
              <Button type='primary' size='large'>
                管理者ログイン
              </Button>
              <Button size='large'>店舗ログイン</Button>
            </Space>
          </Space>
        </Space>
      </Card>
    </div>
  );
}
