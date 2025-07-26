'use client';

import React from 'react';
import { Typography, Space, Button } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="max-w-lg w-full text-center bg-white p-8 rounded-lg shadow-lg">
        <Space direction="vertical" size="large" className="w-full">
          <CheckCircleOutlined style={{ fontSize: '48px', color: '#52c41a' }} />
          
          <Title level={2}>POS QR System</Title>
          
          <Paragraph>
            フロントエンド開発が完了しました。
            <br />
            システムの各機能をご確認ください。
          </Paragraph>

          <Space direction="vertical" className="w-full" size="middle">
            <Button 
              type="primary" 
              size="large"
              onClick={() => window.location.href = '/auth/admin-login'}
              className="w-full"
            >
              管理者ログイン
            </Button>
            
            <Button 
              type="default" 
              size="large"
              onClick={() => window.location.href = '/auth/store-login'}
              className="w-full"
            >
              店舗ログイン
            </Button>
            
            <Button 
              type="dashed" 
              size="large"
              onClick={() => window.location.href = '/order/test-session-123'}
              className="w-full"
            >
              顧客注文画面（テスト）
            </Button>
          </Space>
        </Space>
      </div>
    </div>
  );
}