'use client';

import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import {
  DashboardOutlined,
  TableOutlined,
  ShoppingCartOutlined,
  MenuOutlined,
  QrcodeOutlined,
  BarChartOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useRouter, usePathname } from 'next/navigation';

import { useSidebar } from '@/store';

const { Sider } = Layout;
const { Text } = Typography;

export const StoreSidebar: React.FC = () => {
  const router = useRouter();
  const pathname = usePathname();
  const { collapsed } = useSidebar();

  const menuItems = [
    {
      key: '/store/dashboard',
      icon: <DashboardOutlined />,
      label: 'ダッシュボード',
    },
    {
      key: '/store/seats',
      icon: <TableOutlined />,
      label: '座席管理',
    },
    {
      key: '/store/orders',
      icon: <ShoppingCartOutlined />,
      label: '注文管理',
    },
    {
      key: '/store/menu',
      icon: <MenuOutlined />,
      label: 'メニュー管理',
    },
    {
      key: '/store/qr-codes',
      icon: <QrcodeOutlined />,
      label: 'QRコード管理',
    },
    {
      key: '/store/analytics',
      icon: <BarChartOutlined />,
      label: '売上分析',
    },
    {
      key: '/store/settings',
      icon: <SettingOutlined />,
      label: '店舗設定',
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
          <div className="w-8 h-8 bg-green-500 rounded flex items-center justify-center">
            <Text className="text-white font-bold text-sm">S</Text>
          </div>
        ) : (
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-green-500 rounded flex items-center justify-center">
              <Text className="text-white font-bold text-sm">S</Text>
            </div>
            <Text strong className="text-gray-800">
              Store Manager
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

export default StoreSidebar;