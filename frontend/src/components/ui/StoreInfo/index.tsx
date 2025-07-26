'use client';

import React from 'react';
import { Modal, Typography, Space, Divider, Button } from 'antd';
import {
  PhoneOutlined,
  EnvironmentOutlined,
  ClockCircleOutlined,
  InfoCircleOutlined,
  WifiOutlined,
} from '@ant-design/icons';
import { StoreInfoProps } from './StoreInfo.types';

const { Title, Text, Paragraph } = Typography;

export const StoreInfo: React.FC<StoreInfoProps> = ({
  visible,
  store,
  onClose,
}) => {
  if (!store) {
    return null;
  }

  const handleCallStore = () => {
    if (store.phone) {
      window.location.href = `tel:${store.phone}`;
    }
  };

  return (
    <Modal
      title={
        <Space>
          <InfoCircleOutlined />
          <span>店舗情報</span>
        </Space>
      }
      open={visible}
      onCancel={onClose}
      footer={[
        <Button key="close" onClick={onClose}>
          閉じる
        </Button>,
      ]}
      width={500}
      className="store-info-modal"
    >
      <div className="space-y-6">
        {/* Store Name */}
        <div className="text-center">
          <Title level={2} className="mb-2">
            {store.name}
          </Title>
          <Text type="secondary" className="text-lg">
            ようこそお越しくださいました
          </Text>
        </div>

        <Divider />

        {/* Contact Information */}
        <div className="space-y-4">
          <div className="flex items-start space-x-3">
            <EnvironmentOutlined className="text-blue-500 text-lg mt-1" />
            <div className="flex-1">
              <Text strong className="block">住所</Text>
              <Text className="text-gray-600">{store.address}</Text>
            </div>
          </div>

          {store.phone && (
            <div className="flex items-start space-x-3">
              <PhoneOutlined className="text-green-500 text-lg mt-1" />
              <div className="flex-1">
                <Text strong className="block">電話番号</Text>
                <div className="flex items-center space-x-2">
                  <Text className="text-gray-600">{store.phone}</Text>
                  <Button
                    type="link"
                    size="small"
                    onClick={handleCallStore}
                    className="p-0 h-auto"
                  >
                    電話をかける
                  </Button>
                </div>
              </div>
            </div>
          )}

          {/* Business Hours */}
          <div className="flex items-start space-x-3">
            <ClockCircleOutlined className="text-orange-500 text-lg mt-1" />
            <div className="flex-1">
              <Text strong className="block">営業時間</Text>
              <div className="text-gray-600">
                <div>平日: 11:00 - 22:00</div>
                <div>土日祝: 10:00 - 23:00</div>
                <Text type="secondary" className="text-sm">
                  ※ ラストオーダーは閉店30分前
                </Text>
              </div>
            </div>
          </div>

          {/* WiFi Information */}
          <div className="flex items-start space-x-3">
            <WifiOutlined className="text-purple-500 text-lg mt-1" />
            <div className="flex-1">
              <Text strong className="block">Wi-Fi</Text>
              <div className="text-gray-600">
                <div>ネットワーク: {store.name}_WiFi</div>
                <div>パスワード: welcome123</div>
                <Text type="secondary" className="text-sm">
                  ※ 無料でご利用いただけます
                </Text>
              </div>
            </div>
          </div>
        </div>

        <Divider />

        {/* Service Information */}
        <div>
          <Title level={4} className="mb-3">
            サービスのご案内
          </Title>
          <div className="space-y-3">
            <div className="bg-blue-50 p-3 rounded-lg">
              <Text strong className="text-blue-800 block mb-1">
                QR注文システム
              </Text>
              <Text className="text-blue-700 text-sm">
                お席からスマートフォンで簡単にご注文いただけます。
                追加注文も可能です。
              </Text>
            </div>

            <div className="bg-green-50 p-3 rounded-lg">
              <Text strong className="text-green-800 block mb-1">
                お支払い方法
              </Text>
              <Text className="text-green-700 text-sm">
                現金・クレジットカード・電子マネー・QRコード決済に対応しています。
                お会計はレジまでお越しください。
              </Text>
            </div>

            <div className="bg-orange-50 p-3 rounded-lg">
              <Text strong className="text-orange-800 block mb-1">
                アレルギー対応
              </Text>
              <Text className="text-orange-700 text-sm">
                アレルギーをお持ちの方は、注文時の特別な要望欄にご記入いただくか、
                スタッフまでお声がけください。
              </Text>
            </div>
          </div>
        </div>

        <Divider />

        {/* Contact for Help */}
        <div className="text-center bg-gray-50 p-4 rounded-lg">
          <Text strong className="block mb-2">
            ご不明な点がございましたら
          </Text>
          <Text type="secondary" className="text-sm">
            お気軽にスタッフまでお声がけください。
            <br />
            快適なお食事をお楽しみいただけるよう、
            <br />
            スタッフ一同心よりお待ちしております。
          </Text>
        </div>
      </div>

      <style jsx global>{`
        .store-info-modal .ant-modal-body {
          max-height: 70vh;
          overflow-y: auto;
        }
        
        @media (max-width: 768px) {
          .store-info-modal .ant-modal-content {
            margin: 16px;
          }
          
          .store-info-modal .ant-typography h2 {
            font-size: 20px;
          }
          
          .store-info-modal .ant-typography h4 {
            font-size: 16px;
          }
        }
      `}</style>
    </Modal>
  );
};