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
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  ShopOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

import { Card } from '@/components';
import { useStores, useDeleteStore, useToggleStoreStatus } from '@/hooks/api/useStores';
import { usePermissions } from '@/hooks';
import { PERMISSIONS } from '@/lib/auth/permissions';
import { Store } from '@/types';

const { Title, Text } = Typography;
const { Search } = Input;

export default function StoresPage() {
  const router = useRouter();
  const { hasPermission } = usePermissions();
  
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  // API hooks
  const {
    data: storesData,
    isLoading,
    error,
  } = useStores({
    page: currentPage,
    limit: pageSize,
    search: searchTerm || undefined,
  });

  const deleteStoreMutation = useDeleteStore();
  const toggleStatusMutation = useToggleStoreStatus();

  // Permissions
  const canManageStores = hasPermission(PERMISSIONS.ADMIN.MANAGE_STORES);
  const canViewAnalytics = hasPermission(PERMISSIONS.ADMIN.VIEW_ANALYTICS);

  // Handlers
  const handleSearch = (value: string) => {
    setSearchTerm(value);
    setCurrentPage(1);
  };

  const handleCreateStore = () => {
    router.push('/admin/stores/create');
  };

  const handleViewStore = (storeId: string) => {
    router.push(`/admin/stores/${storeId}`);
  };

  const handleEditStore = (storeId: string) => {
    router.push(`/admin/stores/${storeId}/edit`);
  };

  const handleDeleteStore = (store: Store) => {
    Modal.confirm({
      title: '店舗を削除しますか？',
      content: `「${store.name}」を削除すると、関連するデータもすべて削除されます。この操作は取り消せません。`,
      okText: '削除',
      okType: 'danger',
      cancelText: 'キャンセル',
      onOk: () => {
        deleteStoreMutation.mutate(store.id);
      },
    });
  };

  const handleToggleStatus = (store: Store) => {
    const newStatus = !store.isActive;
    const actionText = newStatus ? '有効' : '無効';
    
    Modal.confirm({
      title: `店舗を${actionText}にしますか？`,
      content: `「${store.name}」を${actionText}にします。`,
      okText: actionText,
      cancelText: 'キャンセル',
      onOk: () => {
        toggleStatusMutation.mutate({
          storeId: store.id,
          isActive: newStatus,
        });
      },
    });
  };

  // Table columns
  const columns: ColumnsType<Store> = [
    {
      title: '店舗名',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Store) => (
        <Space>
          <ShopOutlined />
          <div>
            <div className="font-medium">{name}</div>
            {record.description && (
              <Text type="secondary" className="text-sm">
                {record.description}
              </Text>
            )}
          </div>
        </Space>
      ),
    },
    {
      title: '住所',
      dataIndex: 'address',
      key: 'address',
      ellipsis: true,
    },
    {
      title: '連絡先',
      key: 'contact',
      render: (_, record: Store) => (
        <div>
          <div>{record.phone}</div>
          <Text type="secondary" className="text-sm">
            {record.email}
          </Text>
        </div>
      ),
    },
    {
      title: 'ステータス',
      dataIndex: 'isActive',
      key: 'isActive',
      render: (isActive: boolean, record: Store) => (
        <Space>
          <Tag color={isActive ? 'green' : 'red'}>
            {isActive ? '有効' : '無効'}
          </Tag>
          {canManageStores && (
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
      title: '最大座席数',
      dataIndex: ['settings', 'maxSeats'],
      key: 'maxSeats',
      align: 'center',
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
      render: (_, record: Store) => (
        <Space>
          <Tooltip title="詳細表示">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleViewStore(record.id)}
            />
          </Tooltip>
          {canManageStores && (
            <>
              <Tooltip title="編集">
                <Button
                  type="text"
                  icon={<EditOutlined />}
                  onClick={() => handleEditStore(record.id)}
                />
              </Tooltip>
              <Tooltip title="削除">
                <Button
                  type="text"
                  danger
                  icon={<DeleteOutlined />}
                  loading={deleteStoreMutation.isPending}
                  onClick={() => handleDeleteStore(record)}
                />
              </Tooltip>
            </>
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
            <Text type="danger">店舗データの読み込みに失敗しました</Text>
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
            <Title level={2}>店舗管理</Title>
            <Text type="secondary">
              システム内の全店舗を管理します
            </Text>
          </Col>
          <Col>
            {canManageStores && (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleCreateStore}
              >
                新しい店舗を作成
              </Button>
            )}
          </Col>
        </Row>

        {/* Search and Filters */}
        <Card>
          <Row gutter={16}>
            <Col xs={24} sm={12} md={8}>
              <Search
                placeholder="店舗名で検索..."
                allowClear
                onSearch={handleSearch}
                onChange={(e) => !e.target.value && handleSearch('')}
              />
            </Col>
          </Row>
        </Card>

        {/* Stores Table */}
        <Card>
          <Table
            columns={columns}
            dataSource={storesData?.stores || []}
            rowKey="id"
            loading={isLoading}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              total: storesData?.total || 0,
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