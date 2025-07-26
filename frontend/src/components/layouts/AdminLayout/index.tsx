'use client';

import React from 'react';
import { Layout } from 'antd';

import { AdminGuard } from '@/lib/auth';
import { ErrorBoundary } from '@/components/ui/ErrorBoundary';
import { AdminHeader } from './Header';
import { AdminSidebar } from './Sidebar';

const { Content } = Layout;

interface AdminLayoutProps {
  children: React.ReactNode;
  title?: string;
}

export const AdminLayout: React.FC<AdminLayoutProps> = ({ 
  children, 
  title 
}) => {
  return (
    <AdminGuard>
      <ErrorBoundary>
        <Layout className="min-h-screen">
          <AdminSidebar />
          <Layout>
            <AdminHeader title={title} />
            <Content className="bg-gray-50 p-6 overflow-auto">
              <div className="max-w-7xl mx-auto">
                {children}
              </div>
            </Content>
          </Layout>
        </Layout>
      </ErrorBoundary>
    </AdminGuard>
  );
};

export default AdminLayout;