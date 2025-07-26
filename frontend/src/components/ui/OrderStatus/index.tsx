'use client';

import React, { useEffect, useState } from 'react';
import { Card, Typography, Progress, Space, Tag, Button, Alert } from 'antd';
import {
  ClockCircleOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { OrderStatusProps } from './OrderStatus.types';
import { ORDER_STATUS_LABELS, ORDER_STATUS_COLORS } from '@/types/models';

const { Title, Text } = Typography;

export const OrderStatus: React.FC<OrderStatusProps> = ({
  order,
  onRefresh,
  refreshing = false,
}) => {
  const [timeElapsed, setTimeElapsed] = useState(0);
  const [estimatedProgress, setEstimatedProgress] = useState(0);

  useEffect(() => {
    if (!order) return;

    const interval = setInterval(() => {
      const now = new Date();
      const placedAt = new Date(order.placedAt);
      const elapsed = Math.floor((now.getTime() - placedAt.getTime()) / 1000 / 60); // minutes
      setTimeElapsed(elapsed);

      // Calculate estimated progress based on status and time
      let progress = 0;
      switch (order.status) {
        case 'pending':
          progress = Math.min(10, elapsed * 2);
          break;
        case 'confirmed':
          progress = Math.min(25, 15 + elapsed * 1.5);
          break;
        case 'preparing':
          progress = Math.min(80, 30 + elapsed * 2);
          break;
        case 'ready':
          progress = 100;
          break;
        case 'served':
          progress = 100;
          break;
        case 'cancelled':
          progress = 0;
          break;
        default:
          progress = 0;
      }
      setEstimatedProgress(progress);
    }, 1000);

    return () => clearInterval(interval);
  }, [order]);

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
      case 'confirmed':
      case 'preparing':
        return <ClockCircleOutlined />;
      case 'ready':
      case 'served':
        return <CheckCircleOutlined />;
      case 'cancelled':
        return <ExclamationCircleOutlined />;
      default:
        return <ClockCircleOutlined />;
    }
  };

  const getStatusMessage = (status: string) => {
    switch (status) {
      case 'pending':
        return 'ご注文を受け付けました。確認中です...';
      case 'confirmed':
        return 'ご注文を確認しました。調理を開始します。';
      case 'preparing':
        return '調理中です。もうしばらくお待ちください。';
      case 'ready':
        return 'お料理が完成しました！スタッフがお持ちします。';
      case 'served':
        return 'お料理をお渡ししました。ありがとうございます！';
      case 'cancelled':
        return 'ご注文がキャンセルされました。';
      default:
        return '状況を確認中です...';
    }
  };

  if (!order) {
    return null;
  }

  const isCompleted = order.status === 'served';
  const isCancelled = order.status === 'cancelled';

  return (
    <Card className="order-status-card">
      <div className="space-y-4">
        {/* Header */}
        <div className="flex justify-between items-start">
          <div>
            <Title level={4} className="m-0">
              注文番号: {order.orderNumber}
            </Title>
            <Text type="secondary" className="text-sm">
              注文時刻: {formatDateTime(order.placedAt)}
            </Text>
          </div>
          <Space>
            <Tag
              color={ORDER_STATUS_COLORS[order.status]}
              icon={getStatusIcon(order.status)}
              className="text-base px-3 py-1"
            >
              {ORDER_STATUS_LABELS[order.status]}
            </Tag>
            <Button
              type="text"
              icon={<ReloadOutlined />}
              onClick={onRefresh}
              loading={refreshing}
              size="small"
            />
          </Space>
        </div>

        {/* Progress Bar */}
        {!isCancelled && (
          <div>
            <Progress
              percent={estimatedProgress}
              status={isCompleted ? 'success' : 'active'}
              strokeColor={isCompleted ? '#52c41a' : '#1890ff'}
              trailColor="#f0f0f0"
              strokeWidth={8}
              showInfo={false}
            />
            <div className="flex justify-between text-xs text-gray-500 mt-1">
              <span>注文受付</span>
              <span>調理中</span>
              <span>完成</span>
            </div>
          </div>
        )}

        {/* Status Message */}
        <Alert
          message={getStatusMessage(order.status)}
          type={
            isCompleted
              ? 'success'
              : isCancelled
              ? 'error'
              : order.status === 'ready'
              ? 'warning'
              : 'info'
          }
          showIcon
          className="text-center"
        />

        {/* Time Information */}
        <div className="bg-gray-50 p-3 rounded-lg">
          <div className="flex justify-between items-center text-sm">
            <Text type="secondary">経過時間</Text>
            <Text strong>{timeElapsed}分</Text>
          </div>
          {order.estimatedReadyAt && !isCompleted && (
            <div className="flex justify-between items-center text-sm mt-1">
              <Text type="secondary">完成予定</Text>
              <Text strong>{formatDateTime(order.estimatedReadyAt)}</Text>
            </div>
          )}
          {order.completedAt && (
            <div className="flex justify-between items-center text-sm mt-1">
              <Text type="secondary">完了時刻</Text>
              <Text strong>{formatDateTime(order.completedAt)}</Text>
            </div>
          )}
        </div>

        {/* Special Instructions */}
        {order.specialInstructions && (
          <div className="bg-yellow-50 p-3 rounded-lg border border-yellow-200">
            <Text strong className="text-yellow-800 text-sm">
              特別な要望:
            </Text>
            <Text className="text-yellow-700 text-sm block mt-1">
              {order.specialInstructions}
            </Text>
          </div>
        )}
      </div>

      <style jsx global>{`
        .order-status-card {
          border: 1px solid #d9d9d9;
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        
        .order-status-card .ant-progress-bg {
          transition: width 0.3s ease;
        }
        
        @media (max-width: 768px) {
          .order-status-card .ant-typography h4 {
            font-size: 16px;
          }
          
          .order-status-card .ant-tag {
            font-size: 12px;
            padding: 2px 8px;
          }
        }
      `}</style>
    </Card>
  );
};