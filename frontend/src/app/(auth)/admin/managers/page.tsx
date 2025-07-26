'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Typography,
  Space,
  Table,
  Button,
  Input,
  Switch,
  Tag,
  Modal,
  Tooltip,
  Row,
  Col,
  Avatar,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  UserOutlined,
  CrownOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

import { Card } from '@/components';
import { useManagers, useDeleteManager, useToggleManagerStatus, Manager } from '@/hooks/api/useManagers';
import { usePermissions } from '@/hooks';
import { PERMISSIONS, permissionUtils } from '@/lib/auth/permissions';

const { Title, Text } = Typography;
const { Search } = Input;

export default function ManagersPage() {
  const router = useRouter();
  const { hasPermission } = usePermissions();
  
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  // API hooks
  const {
    data: managersData,
    isLoading,
    error,
  } = useManagers({
    page: currentPage,
    limit: pageSize,
    search: searchTerm || undefined,
  });

  const deleteManagerMutation = useDeleteManager();
  const toggleStatusMutation = useToggleManagerStatus();

  // Permissions
  const canManageManagers = hasPermission(PERMISSIONS.ADMIN.MANAGE_MANAGERS);

  // Handlers
  const handleSearch = (value: string) => {
    setSearchTerm(value);
    setCurrentPage(1);
  };

  const handleCreateManager = () => {
    router.push('/admin/managers/create');
  };

  const handleViewManager = (managerId: string) => {
    router.push(`/admin/managers/${managerId}`);
  };

  const handleDeleteManager = (manager: Manager) => {
    Modal.confirm({
      title: '管理者を削除しますか？',
      content: `「${manager.name}」を削除します。この操作は取り消せません。`,
      okText: '削除',
      okType: 'danger',
      cancelText: 'キャンセル',
      onOk: () => {
        deleteManagerMutation.mutate(manager.id);
      },
    });
  };

  const handleToggleStatus = (manager: Manager) => {
    const newStatus = !manager.isActive;
    const actionText = newStatus ? '有効' : '無効';
    
    Modal.confirm({
      title: `管理者を${actionText}にしますか？`,
      content: `「${manager.name}」を${actionText}にします。`,
      okText: actionText,
      cancelText: 'キャンセル',
      onOk: () => {
        toggleStatusMutation.mutate({
          managerId: manager.id,
          isActive: newStatus,
        });
      },
    });
  };

  // Table columns
  const columns: ColumnsType<Manager> = [
    {
      title: '管理者',
      key: 'manager',
      render: (_, record: Manager) => (
        <Space>
          <Avatar icon={<UserOutlined />} />
          <div>
            <div className="font-medium flex items-center gap-2">
              {record.name}
              {record.role === 'admin' && (
                <CrownOutlined style={{ color: '#faad14' }} />
              )}
            </div>
            <Text type="secondary" className="text-sm">
              {record.email}
            </Text>
          </div>
        </Space>
      ),
    },
    {
      title: 'ステータス',
      dataIndex: 'isActive',
      key: 'isActive',
      render: (isActive: boolean, record: Manager) => (
        <Space>
          <Tag color={isActive ? 'green' : 'red'}>
            {isActive ? '有効' : '無効'}
          </Tag>
          {canManageManagers && (
            <Switch
              size="small"
              checked={isActive}
              loading={toggleStatusMutation.isPending}
              onChange={() => handleToggleStatus(record)}
            />
          )}
        </Space>
      ),
    },
    {
      title: '権限',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <div>
          <Text className="text-sm">
            {permissions.length}個の権限
          </Text>
          <div className="mt-1">
            {permissions.slice(0, 2).map((permission) => (
              <Tag key={permission} className="text-xs">
                {permissionUtils.getPermissionDescription(permission)}
              </Tag>
            ))}
            {permissions.length > 2 && (
              <Tag className="text-xs">
                +{permissions.length - 2}
              </Tag>
            )}
          </div>
        </div>
      ),
    },
    {
      title: '最終ログイン',
      dataIndex: 'lastLoginAt',
      key: 'lastLoginAt',
      render: (date?: string) => (
        <Text type="secondary">
          {date ? new Date(date).toLocaleString('ja-JP') : '未ログイン'}
        </Text>
      ),
    },
    {
      title: '作成日',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleDateString('ja-JP'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record: Manager) => (
        <Space>
          <Tooltip title="詳細表示">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleViewManager(record.id)}
            />
          </Tooltip>
          {canManageManagers && (
            <Tooltip title="削除">
              <Button
                type="text"
                danger
                icon={<DeleteOutlined />}
                loading={deleteManagerMutation.isPending}
                onClick={() => handleDeleteManager(record)}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  if (error) {
    return (
      <div className="p-6">
        <Card>
          <div className="text-center py-8">
            <Text type="danger">管理者データの読み込みに失敗しました</Text>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="p-6">
      <Space direction="vertical" size="large" className="w-full">
        {/* Header */}
        <Row justify="space-between" align="middle">
          <Col>
            <Title level={2}>管理者管理</Title>
            <Text type="secondary">
              システム管理者を管理します
            </Text>
          </Col>
          <Col>
            {canManageManagers && (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleCreateManager}
              >
                新しい管理者を作成
              </Button>
            )}
          </Col>
        </Row>

        {/* Search and Filters */}
        <Card>
          <Row gutter={16}>
            <Col xs={24} sm={12} md={8}>
              <Search
                placeholder="管理者名またはメールアドレスで検索..."
                allowClear
                onSearch={handleSearch}
                onChange={(e) => !e.target.value && handleSearch('')}
              />
            </Col>
          </Row>
        </Card>

        {/* Managers Table */}
        <Card>
          <Table
            columns={columns}
            dataSource={managersData?.managers || []}
            rowKey="id"
            loading={isLoading}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              total: managersData?.total || 0,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) =>
                `${range[0]}-${range[1]} / ${total}件`,
              onChange: (page, size) => {
                setCurrentPage(page);
                setPageSize(size || 10);
              },
            }}
            scroll={{ x: 1000 }}
          />
        </Card>
      </Space>
    </div>
  );
}