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
  Image,
} from 'antd';
import {
  ArrowLeftOutlined,
  EditOutlined,
  TableOutlined,
  QrcodeOutlined,
  UserOutlined,
  PrinterOutlined,
} from '@ant-design/icons';

import { Card } from '@/components';
import { useSeat, useSeatQR, Seat, SeatStatus } from '@/hooks/api/useSeats';
import { usePermissions } from '@/hooks';
import { PERMISSIONS } from '@/lib/auth/permissions';

const { Title, Text } = Typography;

interface SeatDetailPageProps {
  params: {
    id: string;
  };
}

export default function SeatDetailPage({ params }: SeatDetailPageProps) {
  const router = useRouter();
  const { hasPermission } = usePermissions();
  const seatId = params.id;

  // API hooks
  const {
    data: seat,
    isLoading: seatLoading,
    error: seatError,
  } = useSeat(seatId);

  const {
    data: qrData,
    isLoading: qrLoading,
  } = useSeatQR(seatId);

  // Permissions
  const canManageSeats = hasPermission(PERMISSIONS.STORE.MANAGE_SEATS);

  // Handlers
  const handleBack = () => {
    router.push('/store/seats');
  };

  const handleViewQR = () => {
    router.push(`/store/seats/${seatId}/qr`);
  };

  const handlePrintQR = () => {
    if (qrData?.qrCodeUrl) {
      const printWindow = window.open('', '_blank');
      if (printWindow) {
        printWindow.document.write(`
          <html>
            <head>
              <title>QRコード - ${seat?.name}</title>
              <style>
                body { 
                  font-family: Arial, sans-serif; 
                  text-align: center; 
                  padding: 20px; 
                }
                .qr-container { 
                  border: 2px solid #000; 
                  padding: 20px; 
                  display: inline-block; 
                  margin: 20px;
                }
                .seat-info { 
                  margin-bottom: 20px; 
                  font-size: 18px; 
                  font-weight: bold; 
                }
                .qr-code { 
                  max-width: 300px; 
                  height: auto; 
                }
                .instructions { 
                  margin-top: 20px; 
                  font-size: 14px; 
                  color: #666; 
                }
              </style>
            </head>
            <body>
              <div class="qr-container">
                <div class="seat-info">${seat?.name} (${seat?.number})</div>
                <img src="${qrData.qrCodeUrl}" alt="QRコード" class="qr-code" />
                <div class="instructions">
                  スマートフォンでQRコードを読み取って注文してください
                </div>
              </div>
            </body>
          </html>
        `);
        printWindow.document.close();
        printWindow.print();
      }
    }
  };

  if (seatLoading) {
    return (
      <div className="p-6 flex justify-center items-center min-h-96">
        <Spin size="large" />
      </div>
    );
  }

  if (seatError || !seat) {
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

  const getStatusConfig = (status: SeatStatus) => {
    const configs = {
      available: { label: '利用可能', color: 'green' },
      occupied: { label: '利用中', color: 'red' },
      reserved: { label: '予約済み', color: 'orange' },
      cleaning: { label: '清掃中', color: 'blue' },
      maintenance: { label: 'メンテナンス', color: 'purple' },
    };
    return configs[status] || configs.available;
  };

  const statusConfig = getStatusConfig(seat.status);

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
            座席一覧に戻る
          </Button>
          <Row justify="space-between" align="middle">
            <Col>
              <Space>
                <TableOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
                <div>
                  <Title level={2} className="!mb-0">
                    {seat.name}
                  </Title>
                  <Text type="secondary">#{seat.number}</Text>
                </div>
              </Space>
            </Col>
            <Col>
              <Space>
                <Tag color={statusConfig.color} className="text-sm px-3 py-1">
                  {statusConfig.label}
                </Tag>
                <Tag color={seat.isActive ? 'green' : 'red'} className="text-sm px-3 py-1">
                  {seat.isActive ? '有効' : '無効'}
                </Tag>
                <Button
                  icon={<QrcodeOutlined />}
                  onClick={handleViewQR}
                >
                  QRコード表示
                </Button>
              </Space>
            </Col>
          </Row>
        </div>

        <Row gutter={16}>
          {/* Basic Information */}
          <Col xs={24} lg={12}>
            <Card title="基本情報" extra={<TableOutlined />}>
              <Descriptions column={1} size="small">
                <Descriptions.Item label="座席番号">
                  #{seat.number}
                </Descriptions.Item>
                <Descriptions.Item label="座席名">
                  {seat.name}
                </Descriptions.Item>
                <Descriptions.Item
                  label={<Space><UserOutlined />定員</Space>}
                >
                  {seat.capacity}名
                </Descriptions.Item>
                <Descriptions.Item label="説明">
                  {seat.description || '説明なし'}
                </Descriptions.Item>
                <Descriptions.Item label="ステータス">
                  <Tag color={statusConfig.color}>
                    {statusConfig.label}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="作成日">
                  {new Date(seat.createdAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
                <Descriptions.Item label="更新日">
                  {new Date(seat.updatedAt).toLocaleString('ja-JP')}
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>

          {/* QR Code Preview */}
          <Col xs={24} lg={12}>
            <Card 
              title="QRコードプレビュー" 
              extra={<QrcodeOutlined />}
              loading={qrLoading}
            >
              {qrData ? (
                <Space direction="vertical" className="w-full text-center">
                  <Image
                    src={qrData.qrCodeUrl}
                    alt="QRコード"
                    width={200}
                    height={200}
                    className="border border-gray-200 rounded"
                  />
                  <div>
                    <Text type="secondary" className="text-sm">
                      顧客はこのQRコードを読み取って注文できます
                    </Text>
                  </div>
                  <Space>
                    <Button
                      type="primary"
                      icon={<QrcodeOutlined />}
                      onClick={handleViewQR}
                    >
                      詳細表示
                    </Button>
                    <Button
                      icon={<PrinterOutlined />}
                      onClick={handlePrintQR}
                    >
                      印刷
                    </Button>
                  </Space>
                </Space>
              ) : (
                <div className="text-center py-8">
                  <Text type="secondary">QRコードを読み込み中...</Text>
                </div>
              )}
            </Card>
          </Col>
        </Row>

        {/* Session URL Info */}
        {qrData && (
          <Card title="注文URL情報">
            <Space direction="vertical" className="w-full">
              <div>
                <Text strong>セッションURL:</Text>
                <div className="mt-2 p-3 bg-gray-50 rounded border">
                  <Text code className="text-sm break-all">
                    {qrData.sessionUrl}
                  </Text>
                </div>
              </div>
              <div className="text-sm text-gray-600">
                <p>• このURLは顧客がQRコードを読み取った際にアクセスされます</p>
                <p>• URLには座席情報とセッション情報が含まれています</p>
                <p>• セキュリティのため、定期的にQRコードを再生成することをお勧めします</p>
              </div>
            </Space>
          </Card>
        )}
      </Space>
    </div>
  );
}