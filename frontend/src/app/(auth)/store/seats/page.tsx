'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Typography,
  Space,
  Row,
  Col,
  Button,
  Input,
  Select,
  Tag,
  Modal,
  Tooltip,
  Card as AntCard,
  Grid,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  QrcodeOutlined,
  TableOutlined,
  UserOutlined,
} from '@ant-design/icons';

import { Card } from '@/components';
import { useSeats, useDeleteSeat, useUpdateSeatStatus, Seat, SeatStatus } from '@/hooks/api/useSeats';
import { usePermissions } from '@/hooks';
import { PERMISSIONS } from '@/lib/auth/permissions';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;
const { useBreakpoint } = Grid;

export default function SeatsPage() {
  const router = useRouter();
  const screens = useBreakpoint();
  const { hasPermission } = usePermissions();
  
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<SeatStatus | undefined>();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);

  // API hooks
  const {
    data: seatsData,
    isLoading,
    error,
  } = useSeats({
    page: currentPage,
    limit: pageSize,
    search: searchTerm || undefined,
    status: statusFilter,
  });

  const deleteSeats = useDeleteSeat();
  const updateSeatStatus = useUpdateSeatStatus();

  // Permissions
  const canManageSeats = hasPermission(PERMISSIONS.STORE.MANAGE_SEATS);

  // Handlers
  const handleSearch = (value: string) => {
    setSearchTerm(value);
    setCurrentPage(1);
  };

  const handleStatusFilter = (status: SeatStatus | undefined) => {
    setStatusFilter(status);
    setCurrentPage(1);
  };

  const handleCreateSeat = () => {
    router.push('/store/seats/create');
  };

  const handleViewSeat = (seatId: string) => {
    router.push(`/store/seats/${seatId}`);
  };

  const handleViewQR = (seatId: string) => {
    router.push(`/store/seats/${seatId}/qr`);
  };

  const handleDeleteSeat = (seat: Seat) => {
    Modal.confirm({
      title: '座席を削除しますか？',
      content: `「${seat.name}」を削除します。この操作は取り消せません。`,
      okText: '削除',
      okType: 'danger',
      cancelText: 'キャンセル',
      onOk: () => {
        deleteSeats.mutate(seat.id);
      },
    });
  };

  const handleStatusChange = (seat: Seat, newStatus: SeatStatus) => {
    updateSeatStatus.mutate({
      seatId: seat.id,
      status: newStatus,
    });
  };

  // Status options
  const statusOptions = [
    { value: 'available', label: '利用可能', color: 'green' },
    { value: 'occupied', label: '利用中', color: 'red' },
    { value: 'reserved', label: '予約済み', color: 'orange' },
    { value: 'cleaning', label: '清掃中', color: 'blue' },
    { value: 'maintenance', label: 'メンテナンス', color: 'purple' },
  ];

  const getStatusConfig = (status: SeatStatus) => {
    return statusOptions.find(option => option.value === status) || statusOptions[0];
  };

  if (error) {
    return (
      <div className="p-6">
        <Card>
          <div className="text-center py-8">
            <Text type="danger">座席データの読み込みに失敗しました</Text>
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
            <Title level={2}>座席管理</Title>
            <Text type="secondary">
              店舗の座席とQRコードを管理します
            </Text>
          </Col>
          <Col>
            {canManageSeats && (
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleCreateSeat}
              >
                新しい座席を作成
              </Button>
            )}
          </Col>
        </Row>

        {/* Search and Filters */}
        <Card>
          <Row gutter={16}>
            <Col xs={24} sm={12} md={8}>
              <Search
                placeholder="座席名・番号で検索..."
                allowClear
                onSearch={handleSearch}
                onChange={(e) => !e.target.value && handleSearch('')}
              />
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Select
                placeholder="ステータスで絞り込み"
                allowClear
                className="w-full"
                onChange={handleStatusFilter}
              >
                {statusOptions.map(option => (
                  <Option key={option.value} value={option.value}>
                    <Tag color={option.color} className="mr-2">
                      {option.label}
                    </Tag>
                  </Option>
                ))}
              </Select>
            </Col>
          </Row>
        </Card>

        {/* Seats Grid */}
        <Row gutter={[16, 16]}>
          {seatsData?.seats.map((seat) => {
            const statusConfig = getStatusConfig(seat.status);
            
            return (
              <Col
                key={seat.id}
                xs={24}
                sm={12}
                md={8}
                lg={6}
                xl={4}
              >
                <AntCard
                  size="small"
                  className="h-full"
                  actions={[
                    <Tooltip title="詳細表示" key="view">
                      <EyeOutlined onClick={() => handleViewSeat(seat.id)} />
                    </Tooltip>,
                    <Tooltip title="QRコード" key="qr">
                      <QrcodeOutlined onClick={() => handleViewQR(seat.id)} />
                    </Tooltip>,
                    ...(canManageSeats ? [
                      <Tooltip title="削除" key="delete">
                        <DeleteOutlined 
                          style={{ color: '#ff4d4f' }}
                          onClick={() => handleDeleteSeat(seat)}
                        />
                      </Tooltip>
                    ] : [])
                  ]}
                >
                  <Space direction="vertical" size="small" className="w-full">
                    {/* Seat Header */}
                    <div className="flex justify-between items-start">
                      <div>
                        <Text strong className="text-lg">
                          {seat.name}
                        </Text>
                        <div>
                          <Text type="secondary" className="text-sm">
                            #{seat.number}
                          </Text>
                        </div>
                      </div>
                      <TableOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
                    </div>

                    {/* Seat Info */}
                    <div>
                      <Space>
                        <UserOutlined />
                        <Text>{seat.capacity}名</Text>
                      </Space>
                    </div>

                    {/* Status */}
                    <div>
                      {canManageSeats ? (
                        <Select
                          value={seat.status}
                          size="small"
                          className="w-full"
                          onChange={(status) => handleStatusChange(seat, status)}
                        >
                          {statusOptions.map(option => (
                            <Option key={option.value} value={option.value}>
                              <Tag color={option.color} className="mr-1">
                                {option.label}
                              </Tag>
                            </Option>
                          ))}
                        </Select>
                      ) : (
                        <Tag color={statusConfig.color} className="w-full text-center">
                          {statusConfig.label}
                        </Tag>
                      )}
                    </div>

                    {/* Description */}
                    {seat.description && (
                      <Text type="secondary" className="text-xs" ellipsis>
                        {seat.description}
                      </Text>
                    )}

                    {/* Active Status */}
                    <div className="flex justify-between items-center">
                      <Tag color={seat.isActive ? 'green' : 'red'}>
                        {seat.isActive ? '有効' : '無効'}
                      </Tag>
                      <Text type="secondary" className="text-xs">
                        作成: {new Date(seat.createdAt).toLocaleDateString('ja-JP')}
                      </Text>
                    </div>
                  </Space>
                </AntCard>
              </Col>
            );
          })}
        </Row>

        {/* Empty State */}
        {!isLoading && (!seatsData?.seats || seatsData.seats.length === 0) && (
          <Card>
            <div className="text-center py-8">
              <TableOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
              <Title level={4} type="secondary">座席がありません</Title>
              <Text type="secondary">
                {searchTerm || statusFilter ? '検索条件に一致する座席がありません' : '最初の座席を作成してください'}
              </Text>
              {canManageSeats && !searchTerm && !statusFilter && (
                <div className="mt-4">
                  <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateSeat}>
                    座席を作成
                  </Button>
                </div>
              )}
            </div>
          </Card>
        )}

        {/* Pagination */}
        {seatsData && seatsData.total > pageSize && (
          <div className="text-center">
            <Button
              disabled={currentPage === 1}
              onClick={() => setCurrentPage(currentPage - 1)}
            >
              前のページ
            </Button>
            <span className="mx-4">
              {currentPage} / {Math.ceil(seatsData.total / pageSize)}
            </span>
            <Button
              disabled={currentPage >= Math.ceil(seatsData.total / pageSize)}
              onClick={() => setCurrentPage(currentPage + 1)}
            >
              次のページ
            </Button>
          </div>
        )}
      </Space>
    </div>
  );
}