'use client';

import React from 'react';
import { Layout } from 'antd';

import { StoreGuard } from '@/lib/auth';
import { ErrorBoundary } from '@/components/ui/ErrorBoundary';
import { StoreHeader } from './Header';
import { StoreSidebar } from './Sidebar';

const { Content } = Layout;

interface StoreLayoutProps {
  children: React.ReactNode;
  title?: string;
}

export const StoreLayout: React.FC<StoreLayoutProps> = ({ 
  children, 
  title 
}) => {
  return (
    <StoreGuard>
      <ErrorBoundary>
        <Layout className="min-h-screen">
          <StoreSidebar />
          <Layout>
            <StoreHeader title={title} />
            <Content className="bg-gray-50 p-6 overflow-auto">
              <div className="max-w-7xl mx-auto">
                {children}
              </div>
            </Content>
          </Layout>
        </Layout>
      </ErrorBoundary>
    </StoreGuard>
  );
};

export default StoreLayout;