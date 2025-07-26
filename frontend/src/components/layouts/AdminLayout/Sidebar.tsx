'use client';

import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import {
  DashboardOutlined,
  ShopOutlined,
  UserOutlined,
  SettingOutlined,
  BarChartOutlined,
} from '@ant-design/icons';
import { useRouter, usePathname } from 'next/navigation';

import { useSidebar } from '@/store';

const { Sider } = Layout;
const { Text } = Typography;

export const AdminSidebar: React.FC = () => {
  const router = useRouter();
  const pathname = usePathname();
  const { collapsed } = useSidebar();

  const menuItems = [
    {
      key: '/admin/dashboard',
      icon: <DashboardOutlined />,
      label: 'ダッシュボード',
    },
    {
      key: '/admin/stores',
      icon: <ShopOutlined />,
      label: '店舗管理',
    },
    {
      key: '/admin/managers',
      icon: <UserOutlined />,
      label: '管理者管理',
    },
    {
      key: '/admin/analytics',
      icon: <BarChartOutlined />,
      label: '分析・レポート',
    },
    {
      key: '/admin/settings',
      icon: <SettingOutlined />,
      label: 'システム設定',
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    router.push(key);
  };

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      className="bg-white shadow-lg border-r border-gray-200"
      width={240}
      collapsedWidth={80}
    >
      {/* Logo */}
      <div className="h-16 flex items-center justify-center border-b border-gray-200">
        {collapsed ? (
          <div className="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
            <Text className="text-white font-bold text-sm">P</Text>
          </div>
        ) : (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
              <Text className="text-white font-bold text-sm">P</Text>
            </div>
            <Text strong className="text-gray-800">
              POS QR System
            </Text>
          </div>
        )}
      </div>

      {/* Menu */}
      <Menu
        mode="inline"
        selectedKeys={[pathname]}
        items={menuItems}
        onClick={handleMenuClick}
        className="border-none"
        style={{ height: 'calc(100vh - 64px)', borderRight: 0 }}
      />
    </Sider>
  );
};

export default AdminSidebar;