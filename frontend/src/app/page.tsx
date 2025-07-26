'use client';

import { Button, Card, Typography, Space } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export default function Home() {
  return (
    <div className='min-h-screen bg-gray-50 flex items-center justify-center p-4'>
      <Card className='max-w-md w-full text-center'>
        <Space direction='vertical' size='large' className='w-full'>
          <CheckCircleOutlined style={{ fontSize: '48px', color: '#52c41a' }} />

          <Title level={2}>POS QR System</Title>

          <Paragraph>
            Next.js + TypeScript + Ant Design の統合が完了しました。
            開発環境が正常に動作しています。
          </Paragraph>

          <Space>
            <Button type='primary' size='large'>
              管理者ログイン
            </Button>
            <Button size='large'>店舗ログイン</Button>
          </Space>
        </Space>
      </Card>
    </div>
  );
}
