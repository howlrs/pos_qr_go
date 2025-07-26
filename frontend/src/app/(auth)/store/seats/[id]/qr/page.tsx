'use client';

import { useRouter } from 'next/navigation';
import {
  Typography,
  Space,
  Row,
  Col,
  Button,
  Spin,
  Image,
  Modal,
  message,
} from 'antd';
import {
  ArrowLeftOutlined,
  PrinterOutlined,
  DownloadOutlined,
  ReloadOutlined,
  CopyOutlined,
} from '@ant-design/icons';

import { Card } from '@/components';
import { useSeat, useSeatQR, useRegenerateQR } from '@/hooks/api/useSeats';
import { usePermissions } from '@/hooks';
import { PERMISSIONS } from '@/lib/auth/permissions';

const { Title, Text } = Typography;

interface SeatQRPageProps {
  params: {
    id: string;
  };
}

export default function SeatQRPage({ params }: SeatQRPageProps) {
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
    error: qrError,
  } = useSeatQR(seatId);

  const regenerateQRMutation = useRegenerateQR();

  // Permissions
  const canManageSeats = hasPermission(PERMISSIONS.STORE.MANAGE_SEATS);

  // Handlers
  const handleBack = () => {
    router.push(`/store/seats/${seatId}`);
  };

  const handlePrint = () => {
    if (qrData?.qrCodeUrl && seat) {
      const printWindow = window.open('', '_blank');
      if (printWindow) {
        printWindow.document.write(`
          <html>
            <head>
              <title>QRコード - ${seat.name}</title>
              <style>
                @media print {
                  body { margin: 0; }
                  .no-print { display: none; }
                }
                body { 
                  font-family: Arial, sans-serif; 
                  text-align: center; 
                  padding: 40px 20px; 
                  background: white;
                }
                .qr-container { 
                  border: 3px solid #000; 
                  padding: 30px; 
                  display: inline-block; 
                  margin: 20px;
                  background: white;
                  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
                }
                .seat-info { 
                  margin-bottom: 25px; 
                  font-size: 24px; 
                  font-weight: bold; 
                  color: #333;
                }
                .seat-number {
                  font-size: 18px;
                  color: #666;
                  margin-bottom: 20px;
                }
                .qr-code { 
                  max-width: 300px; 
                  height: auto; 
                  border: 1px solid #ddd;
                }
                .instructions { 
                  margin-top: 25px; 
                  font-size: 16px; 
                  color: #666; 
                  line-height: 1.5;
                }
                .footer {
                  margin-top: 30px;
                  font-size: 12px;
                  color: #999;
                }
              </style>
            </head>
            <body>
              <div class="qr-container">
                <div class="seat-info">${seat.name}</div>
                <div class="seat-number">座席番号: ${seat.number}</div>
                <img src="${qrData.qrCodeUrl}" alt="QRコード" class="qr-code" />
                <div class="instructions">
                  <strong>ご注文方法</strong><br>
                  スマートフォンでこのQRコードを読み取り、<br>
                  表示されたページからご注文ください
                </div>
                <div class="footer">
                  定員: ${seat.capacity}名 | 生成日時: ${new Date().toLocaleString('ja-JP')}
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

  const handleDownload = async () => {
    if (qrData?.qrCodeUrl) {
      try {
        const response = await fetch(qrData.qrCodeUrl);
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `qr-code-${seat?.number || seatId}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
        message.success('QRコードをダウンロードしました');
      } catch (error) {
        message.error('ダウンロードに失敗しました');
      }
    }
  };

  const handleCopyURL = async () => {
    if (qrData?.sessionUrl) {
      try {
        await navigator.clipboard.writeText(qrData.sessionUrl);
        message.success('URLをクリップボードにコピーしました');
      } catch (error) {
        message.error('コピーに失敗しました');
      }
    }
  };

  const handleRegenerateQR = () => {
    Modal.confirm({
      title: 'QRコードを再生成しますか？',
      content: '現在のQRコードは無効になり、新しいQRコードが生成されます。印刷済みのQRコードは使用できなくなります。',
      okText: '再生成',
      cancelText: 'キャンセル',
      onOk: () => {
        regenerateQRMutation.mutate(seatId);
      },
    });
  };

  if (seatLoading || qrLoading) {
    return (
      <div className="p-6 flex justify-center items-center min-h-96">
        <Spin size="large" />
      </div>
    );
  }

  if (seatError || qrError || !seat || !qrData) {
    return (
      <div className="p-6">
        <Card>
          <div className="text-center py-8">
            <Text type="danger">データの読み込みに失敗しました</Text>
          </div>
        </Card>
      </div>
    );
  }

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
            座席詳細に戻る
          </Button>
          <Row justify="space-between" align="middle">
            <Col>
              <Title level={2}>QRコード - {seat.name}</Title>
              <Text type="secondary">
                座席番号: #{seat.number} | 定員: {seat.capacity}名
              </Text>
            </Col>
            <Col>
              <Space>
                <Button
                  icon={<PrinterOutlined />}
                  onClick={handlePrint}
                >
                  印刷
                </Button>
                <Button
                  icon={<DownloadOutlined />}
                  onClick={handleDownload}
                >
                  ダウンロード
                </Button>
                {canManageSeats && (
                  <Button
                    icon={<ReloadOutlined />}
                    onClick={handleRegenerateQR}
                    loading={regenerateQRMutation.isPending}
                  >
                    再生成
                  </Button>
                )}
              </Space>
            </Col>
          </Row>
        </div>

        <Row gutter={16}>
          {/* QR Code Display */}
          <Col xs={24} lg={12}>
            <Card title="QRコード" className="text-center">
              <Space direction="vertical" size="large" className="w-full">
                <div className="p-8 bg-white border-2 border-gray-200 rounded-lg inline-block">
                  <Image
                    src={qrData.qrCodeUrl}
                    alt="QRコード"
                    width={300}
                    height={300}
                    className="border border-gray-100"
                    preview={false}
                  />
                </div>
                
                <div>
                  <Title level={4}>{seat.name}</Title>
                  <Text type="secondary">座席番号: #{seat.number}</Text>
                </div>

                <div className="text-left bg-blue-50 p-4 rounded">
                  <Title level={5} className="!mb-2">📱 お客様へのご案内</Title>
                  <ul className="text-sm mb-0">
                    <li>スマートフォンのカメラでQRコードを読み取ってください</li>
                    <li>表示されたページからメニューを選択してご注文ください</li>
                    <li>注文内容の確認後、決済を行ってください</li>
                  </ul>
                </div>
              </Space>
            </Card>
          </Col>

          {/* QR Code Information */}
          <Col xs={24} lg={12}>
            <Space direction="vertical" size="middle" className="w-full">
              {/* Session URL */}
              <Card title="注文URL">
                <Space direction="vertical" className="w-full">
                  <div>
                    <Text strong>セッションURL:</Text>
                    <div className="mt-2 p-3 bg-gray-50 rounded border">
                      <Text code className="text-sm break-all">
                        {qrData.sessionUrl}
                      </Text>
                    </div>
                    <div className="mt-2">
                      <Button
                        size="small"
                        icon={<CopyOutlined />}
                        onClick={handleCopyURL}
                      >
                        URLをコピー
                      </Button>
                    </div>
                  </div>
                </Space>
              </Card>

              {/* QR Code Info */}
              <Card title="QRコード情報">
                <Space direction="vertical" className="w-full">
                  <div>
                    <Text strong>QRコードID:</Text>
                    <div className="mt-1">
                      <Text code>{qrData.qrCode}</Text>
                    </div>
                  </div>
                  <div>
                    <Text strong>生成日時:</Text>
                    <div className="mt-1">
                      <Text>{new Date().toLocaleString('ja-JP')}</Text>
                    </div>
                  </div>
                </Space>
              </Card>

              {/* Usage Instructions */}
              <Card title="使用方法">
                <Space direction="vertical" className="w-full">
                  <div>
                    <Text strong>印刷時の注意:</Text>
                    <ul className="text-sm mt-2 mb-3">
                      <li>A4サイズで印刷することをお勧めします</li>
                      <li>QRコードが鮮明に印刷されることを確認してください</li>
                      <li>汚れや損傷を防ぐため、ラミネート加工をお勧めします</li>
                    </ul>
                  </div>
                  <div>
                    <Text strong>設置場所:</Text>
                    <ul className="text-sm mt-2 mb-0">
                      <li>テーブルの見やすい位置に設置してください</li>
                      <li>照明が十分にある場所を選んでください</li>
                      <li>お客様が座った状態で読み取りやすい角度に調整してください</li>
                    </ul>
                  </div>
                </Space>
              </Card>
            </Space>
          </Col>
        </Row>
      </Space>
    </div>
  );
}