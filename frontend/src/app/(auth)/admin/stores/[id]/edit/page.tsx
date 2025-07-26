'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {
  Typography,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Space,
  Row,
  Col,
  Divider,
  Spin,
} from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';

import { Button, Card } from '@/components';
import { useStore, useUpdateStore } from '@/hooks/api/useStores';
import { UpdateStoreRequest, StoreFeature } from '@/types';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { Option } = Select;

interface StoreFormData {
  name: string;
  description?: string;
  address: string;
  phone: string;
  email: string;
  isActive: boolean;
  timezone: string;
  currency: string;
  language: string;
  orderTimeout: number;
  maxSeats: number;
  features: StoreFeature[];
}

interface EditStorePageProps {
  params: {
    id: string;
  };
}

export default function EditStorePage({ params }: EditStorePageProps) {
  const router = useRouter();
  const storeId = params.id;
  const [form] = Form.useForm<StoreFormData>();

  // API hooks
  const {
    data: store,
    isLoading: storeLoading,
    error: storeError,
  } = useStore(storeId);

  const updateStoreMutation = useUpdateStore();

  // Initialize form with store data
  useEffect(() => {
    if (store) {
      form.setFieldsValue({
        name: store.name,
        description: store.description,
        address: store.address,
        phone: store.phone,
        email: store.email,
        isActive: store.isActive,
        timezone: store.settings.timezone,
        currency: store.settings.currency,
        language: store.settings.language,
        orderTimeout: store.settings.orderTimeout,
        maxSeats: store.settings.maxSeats,
        features: store.settings.features,
      });
    }
  }, [store, form]);

  const handleSubmit = async (values: StoreFormData) => {
    try {
      const updateData: UpdateStoreRequest = {
        name: values.name,
        description: values.description,
        address: values.address,
        phone: values.phone,
        email: values.email,
        isActive: values.isActive,
        settings: {
          timezone: values.timezone,
          currency: values.currency,
          language: values.language,
          orderTimeout: values.orderTimeout,
          maxSeats: values.maxSeats,
          features: values.features,
        },
      };

      await updateStoreMutation.mutateAsync({
        storeId,
        data: updateData,
      });

      router.push(`/admin/stores/${storeId}`);
    } catch (error) {
      console.error('Store update error:', error);
    }
  };

  const handleCancel = () => {
    router.push(`/admin/stores/${storeId}`);
  };

  if (storeLoading) {
    return (
      <div className="p-6 flex justify-center items-center min-h-96">
        <Spin size="large" />
      </div>
    );
  }

  if (storeError || !store) {
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
        <div>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={handleCancel}
            className="mb-4"
          >
            戻る
          </Button>
          <Title level={2}>店舗を編集</Title>
          <Text type="secondary">
            「{store.name}」の情報を編集します
          </Text>
        </div>

        {/* Form */}
        <Card>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
          >
            {/* Status */}
            <Row justify="end">
              <Col>
                <Form.Item
                  name="isActive"
                  valuePropName="checked"
                  label="店舗ステータス"
                >
                  <Switch
                    checkedChildren="有効"
                    unCheckedChildren="無効"
                  />
                </Form.Item>
              </Col>
            </Row>

            {/* Basic Information */}
            <Title level={4}>基本情報</Title>
            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="name"
                  label="店舗名"
                  rules={[
                    { required: true, message: '店舗名を入力してください' },
                    { min: 2, message: '店舗名は2文字以上で入力してください' },
                  ]}
                >
                  <Input placeholder="例: 渋谷店" />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="email"
                  label="メールアドレス"
                  rules={[
                    { required: true, message: 'メールアドレスを入力してください' },
                    { type: 'email', message: '正しいメールアドレスを入力してください' },
                  ]}
                >
                  <Input placeholder="store@example.com" />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="description"
              label="店舗説明"
            >
              <TextArea
                rows={3}
                placeholder="店舗の特徴や説明を入力してください（任意）"
              />
            </Form.Item>

            <Row gutter={16}>
              <Col xs={24} md={16}>
                <Form.Item
                  name="address"
                  label="住所"
                  rules={[
                    { required: true, message: '住所を入力してください' },
                  ]}
                >
                  <Input placeholder="東京都渋谷区..." />
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="phone"
                  label="電話番号"
                  rules={[
                    { required: true, message: '電話番号を入力してください' },
                    { pattern: /^[0-9-+()]+$/, message: '正しい電話番号を入力してください' },
                  ]}
                >
                  <Input placeholder="03-1234-5678" />
                </Form.Item>
              </Col>
            </Row>

            <Divider />

            {/* Settings */}
            <Title level={4}>店舗設定</Title>
            <Row gutter={16}>
              <Col xs={24} md={8}>
                <Form.Item
                  name="timezone"
                  label="タイムゾーン"
                  rules={[{ required: true, message: 'タイムゾーンを選択してください' }]}
                >
                  <Select>
                    <Option value="Asia/Tokyo">Asia/Tokyo (JST)</Option>
                    <Option value="UTC">UTC</Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="currency"
                  label="通貨"
                  rules={[{ required: true, message: '通貨を選択してください' }]}
                >
                  <Select>
                    <Option value="JPY">日本円 (JPY)</Option>
                    <Option value="USD">米ドル (USD)</Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col xs={24} md={8}>
                <Form.Item
                  name="language"
                  label="言語"
                  rules={[{ required: true, message: '言語を選択してください' }]}
                >
                  <Select>
                    <Option value="ja">日本語</Option>
                    <Option value="en">English</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="orderTimeout"
                  label="注文タイムアウト（分）"
                  rules={[
                    { required: true, message: '注文タイムアウトを入力してください' },
                    { type: 'number', min: 5, max: 120, message: '5-120分の範囲で入力してください' },
                  ]}
                >
                  <InputNumber
                    min={5}
                    max={120}
                    className="w-full"
                    placeholder="30"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="maxSeats"
                  label="最大座席数"
                  rules={[
                    { required: true, message: '最大座席数を入力してください' },
                    { type: 'number', min: 1, max: 500, message: '1-500席の範囲で入力してください' },
                  ]}
                >
                  <InputNumber
                    min={1}
                    max={500}
                    className="w-full"
                    placeholder="20"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="features"
              label="利用可能機能"
              rules={[{ required: true, message: '少なくとも1つの機能を選択してください' }]}
            >
              <Select
                mode="multiple"
                placeholder="機能を選択してください"
              >
                <Option value="qr_ordering">QR注文</Option>
                <Option value="table_service">テーブルサービス</Option>
                <Option value="takeaway">テイクアウト</Option>
                <Option value="delivery">デリバリー</Option>
                <Option value="payment_integration">決済連携</Option>
              </Select>
            </Form.Item>

            {/* Actions */}
            <Divider />
            <Row justify="end" gutter={16}>
              <Col>
                <Button onClick={handleCancel}>
                  キャンセル
                </Button>
              </Col>
              <Col>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={updateStoreMutation.isPending}
                >
                  変更を保存
                </Button>
              </Col>
            </Row>
          </Form>
        </Card>
      </Space>
    </div>
  );
}