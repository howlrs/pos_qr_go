'use client';

import React from 'react';
import { Layout, Typography, Space } from 'antd';

import { ErrorBoundary } from '@/components/ui/ErrorBoundary';
import { Card } from '@/components/ui';

const { Content } = Layout;
const { Title, Text } = Typography;

interface AuthLayoutProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
  showLogo?: boolean;
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({
  children,
  title = 'POS QR System',
  subtitle = 'レストラン注文管理システム',
  showLogo = true,
}) => {
  return (
    <ErrorBoundary>
      <Layout className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
        <Content className="flex items-center justify-center p-4">
          <div className="w-full max-w-md">
            {showLogo && (
              <div className="text-center mb-8">
                <Space direction="vertical" size="small">
                  <div className="w-16 h-16 bg-blue-500 rounded-full flex items-center justify-center mx-auto mb-4">
                    <Text className="text-white font-bold text-2xl">P</Text>
                  </div>
                  <Title level={2} className="text-gray-800 mb-0">
                    {title}
                  </Title>
                  <Text type="secondary" className="text-base">
                    {subtitle}
                  </Text>
                </Space>
              </div>
            )}

            <Card
              shadow="large"
              padding="large"
              className="bg-white"
            >
              {children}
            </Card>

            <div className="text-center mt-6">
              <Text type="secondary" className="text-sm">
                © 2025 POS QR System. All rights reserved.
              </Text>
            </div>
          </div>
        </Content>
      </Layout>
    </ErrorBoundary>
  );
};

export default AuthLayout;