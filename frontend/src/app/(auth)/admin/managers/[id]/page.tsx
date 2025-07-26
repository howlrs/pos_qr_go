'use client';

import { useRouter } from 'next/navigation';
import {
  Typography,
  Space,
  Row,
  Col,
  Tag,
  Descriptions,
  Button,
  Spin,
  Avatar,
} from 'antd';
import {
  ArrowLeftOutlined,
  UserOutlined,
  MailOutlined,
  CrownOutlined,
  SafetyOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';

import { Card } from '@/components';
import { useManager, Manager } from '@/hooks/api/useManagers';
import { usePermissions } from '@/hooks';
import { PERMISSIONS, permissionUtils } from '@/lib/auth/permissions';

const { Title, Text } = Typography;

interface ManagerDetailPageProps {
  params: {
    id: string;
  };
}

export default function ManagerDetailPage({ params }: ManagerDetailPageProps) {
  const router = useRouter();
  const { hasPermission } = usePermissions();
  const managerId = params.id;

  // API hooks
  const {
    data: manager,
    isLoading: managerLoading,
    error: managerError,
  } = useManager(managerId);

  // Permissions
  const canManageManagers = hasPermission(PERMISSIONS.ADMIN.MANAGE_MANAGERS);

  // Handlers
  const handleBack = () => {
    router.push('/admin/managers');
  };

  if (managerLoading) {
    return (
      <div className="p-6 flex justify-center items-center min-h-96">
        <Spin size="large" />
      </div>
    );
  }

  if (managerError || !manager) {
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

  // Group permissions by category
  const adminPermissions = manager.permissions.filter((p: string) => p.startsWith('admin:'));
  const commonPermissions = manager.permissions.filter((p: string) => p.startsWith('common:'));

  return (
    <div className="p-6">
      <Space direction="vertical" size="large" className="w-full">
        {/* Header */}
        <div>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={handleBack}
            className="mb-4"
          >
            管理者一覧に戻る
          </Button>
          <Row justify="space-between" align="middle">
            <Col>
              <Space>
                <Avatar size={64} icon={<UserOutlined />} />
                <div>
                  <Title level={2} className="!mb-0 flex items-center gap-2">
                    {manager.name}
                    {manager.role === 'admin' && (
                      <CrownOutlined style={{ color: '#faad14' }} />
                    )}
                  </Title>
                  <Text type="secondary">{manager.email}</Text>
                  <div className="mt-2">
                    <Tag color={manager.isActive ? 'green' : 'red'}>
                      {manager.isActive ? '有効' : '無効'}
                    </Tag>
                  </div>
                </div>
              </Space>
            </Col>
          </Row>
        </div>

        <Row gutter={16}>
          {/* Basic Information */}
          <Col xs={24} lg={12}>
            <Card title="基本情報" extra={<UserOutlined />}>
              <Descriptions column={1} size="small">
                <Descriptions.Item
                  label={<Space><MailOutlined />メールアドレス</Space>}
                >
                  {manager.email}
                </Descriptions.Item>
                <Descriptions.Item
                  label={<Space><SafetyOutlined />ロール</Space>}
                >
                  <Tag color="blue">
                    {manager.role === 'admin' ? '管理者' : manager.role}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item
                  label={<Space><ClockCircleOutlined />最終ログイン</Space>}
                >
                  {manager.lastLoginAt 
                    ? new Date(manager.lastLoginAt).toLocaleString('ja-JP')
                    : '未ログイン'
                  }
                </Descriptions.Item>
                <Descriptions.Item label="作成日">
                  {new Date(manager.createdAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
                <Descriptions.Item label="更新日">
                  {new Date(manager.updatedAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>

          {/* Permissions */}
          <Col xs={24} lg={12}>
            <Card title="権限情報" extra={<SafetyOutlined />}>
              <Space direction="vertical" className="w-full">
                <div>
                  <Text strong>総権限数: </Text>
                  <Tag color="blue">{manager.permissions.length}個</Tag>
                </div>
                
                {adminPermissions.length > 0 && (
                  <div>
                    <Text strong className="block mb-2">管理者権限</Text>
                    <Space wrap>
                      {adminPermissions.map((permission: string) => (
                        <Tag key={permission} color="red">
                          {permissionUtils.getPermissionDescription(permission)}
                        </Tag>
                      ))}
                    </Space>
                  </div>
                )}
                
                {commonPermissions.length > 0 && (
                  <div>
                    <Text strong className="block mb-2">共通権限</Text>
                    <Space wrap>
                      {commonPermissions.map((permission: string) => (
                        <Tag key={permission} color="green">
                          {permissionUtils.getPermissionDescription(permission)}
                        </Tag>
                      ))}
                    </Space>
                  </div>
                )}
              </Space>
            </Card>
          </Col>
        </Row>

        {/* Detailed Permissions */}
        <Card title="詳細権限一覧">
          <Row gutter={[16, 16]}>
            {manager.permissions.map((permission: string) => (
              <Col xs={24} sm={12} md={8} lg={6} key={permission}>
                <Card size="small" className="h-full">
                  <Space direction="vertical" size="small" className="w-full">
                    <Tag 
                      color={permissionUtils.isAdminPermission(permission) ? 'red' : 'green'}
                      className="w-full text-center"
                    >
                      {permissionUtils.getPermissionDescription(permission)}
                    </Tag>
                    <Text type="secondary" className="text-xs">
                      {permission}
                    </Text>
                  </Space>
                </Card>
              </Col>
            ))}
          </Row>
        </Card>
      </Space>
    </div>
  );
}