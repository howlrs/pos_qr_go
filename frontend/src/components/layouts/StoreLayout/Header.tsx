'use client';

import React from 'react';
import { Layout, Avatar, Dropdown, Typography, Space, Button, Badge } from 'antd';
import {
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  BellOutlined,
  ShopOutlined,
} from '@ant-design/icons';

import { useAuth } from '@/hooks';
import { useSidebar } from '@/store';

const { Header: AntHeader } = Layout;
const { Text } = Typography;

interface StoreHeaderProps {
  title?: string;
}

export const StoreHeader: React.FC<StoreHeaderProps> = ({ 
  title = '店舗管理システム' 
}) => {
  const { user, logout } = useAuth();
  const { collapsed, toggle } = useSidebar();

  const handleLogout = () => {
    logout.mutate();
  };

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: 'プロフィール',
    },
    {
      key: 'store-settings',
      icon: <ShopOutlined />,
      label: '店舗設定',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '個人設定',
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'ログアウト',
      onClick: handleLogout,
    },
  ];

  return (
    <AntHeader className="bg-white shadow-sm border-b border-gray-200 px-4 flex items-center justify-between">
      <div className="flex items-center space-x-4">
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={toggle}
          className="text-gray-600 hover:text-gray-800"
        />
        <Typography.Title level={4} className="m-0 text-gray-800">
          {title}
        </Typography.Title>
      </div>

      <div className="flex items-center space-x-4">
        {/* Notifications with badge */}
        <Badge count={3} size="small">
          <Button
            type="text"
            icon={<BellOutlined />}
            className="text-gray-600 hover:text-gray-800"
          />
        </Badge>

        {/* User Menu */}
        <Dropdown
          menu={{ items: userMenuItems }}
          placement="bottomRight"
          trigger={['click']}
        >
          <div className="flex items-center space-x-2 cursor-pointer hover:bg-gray-50 px-2 py-1 rounded">
            <Avatar
              size="small"
              icon={<UserOutlined />}
              className="bg-green-500"
            />
            <Space direction="vertical" size={0}>
              <Text strong className="text-sm">
                {user?.name || 'Store Manager'}
              </Text>
              <Text type="secondary" className="text-xs">
                店舗管理者
              </Text>
            </Space>
          </div>
        </Dropdown>
      </div>
    </AntHeader>
  );
};

export default StoreHeader;