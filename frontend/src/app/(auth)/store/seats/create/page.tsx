'use client';

import { useRouter } from 'next/navigation';
import {
  Typography,
  Form,
  Input,
  InputNumber,
  Space,
  Row,
  Col,
  Divider,
} from 'antd';
import { ArrowLeftOutlined, TableOutlined, UserOutlined } from '@ant-design/icons';

import { Button, Card } from '@/components';
import { useCreateSeat, CreateSeatRequest } from '@/hooks/api/useSeats';

const { Title, Text } = Typography;
const { TextArea } = Input;

interface SeatFormData {
  number: string;
  name: string;
  description?: string;
  capacity: number;
}

export default function CreateSeatPage() {
  const router = useRouter();
  const [form] = Form.useForm<SeatFormData>();
  const createSeatMutation = useCreateSeat();

  const handleSubmit = async (values: SeatFormData) => {
    try {
      const seatData: CreateSeatRequest = {
        number: values.number,
        name: values.name,
        description: values.description,
        capacity: values.capacity,
      };

      const newSeat = await createSeatMutation.mutateAsync(seatData);
      router.push(`/store/seats/${newSeat.id}`);
    } catch (error) {
      console.error('Seat creation error:', error);
    }
  };

  const handleCancel = () => {
    router.back();
  };

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
          <Title level={2}>新しい座席を作成</Title>
          <Text type="secondary">
            新しい座席の基本情報を入力してください
          </Text>
        </div>

        {/* Form */}
        <Card>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{
              capacity: 2,
            }}
          >
            {/* Basic Information */}
            <Title level={4}>基本情報</Title>
            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="number"
                  label="座席番号"
                  rules={[
                    { required: true, message: '座席番号を入力してください' },
                    { min: 1, message: '座席番号は1文字以上で入力してください' },
                  ]}
                >
                  <Input
                    prefix={<TableOutlined />}
                    placeholder="例: T01, A-1"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="name"
                  label="座席名"
                  rules={[
                    { required: true, message: '座席名を入力してください' },
                    { min: 2, message: '座席名は2文字以上で入力してください' },
                  ]}
                >
                  <Input placeholder="例: テーブル1, カウンター席A" />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="description"
              label="座席説明"
            >
              <TextArea
                rows={3}
                placeholder="座席の特徴や説明を入力してください（任意）"
              />
            </Form.Item>

            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="capacity"
                  label="定員"
                  rules={[
                    { required: true, message: '定員を入力してください' },
                    { type: 'number', min: 1, max: 20, message: '1-20名の範囲で入力してください' },
                  ]}
                >
                  <InputNumber
                    min={1}
                    max={20}
                    className="w-full"
                    placeholder="2"
                    prefix={<UserOutlined />}
                    suffix="名"
                  />
                </Form.Item>
              </Col>
            </Row>

            {/* Info */}
            <div className="bg-blue-50 p-4 rounded-lg mb-6">
              <Title level={5} className="!mb-2">
                📝 座席作成について
              </Title>
              <ul className="text-sm text-gray-600 space-y-1 mb-0">
                <li>• 座席作成後、自動的にQRコードが生成されます</li>
                <li>• QRコードは座席詳細ページで確認・印刷できます</li>
                <li>• 座席番号は店舗内で一意である必要があります</li>
                <li>• 作成後も座席情報は編集可能です</li>
              </ul>
            </div>

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
                  loading={createSeatMutation.isPending}
                >
                  座席を作成
                </Button>
              </Col>
            </Row>
          </Form>
        </Card>
      </Space>
    </div>
  );
}