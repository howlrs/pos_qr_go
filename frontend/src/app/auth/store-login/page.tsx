'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Form, Input, message, Typography, Space } from 'antd';
import { UserOutlined, LockOutlined, ShopOutlined } from '@ant-design/icons';

import { Button, Card } from '@/components';
import { useAuth } from '@/hooks';
import { StoreLoginRequest } from '@/types';

const { Title, Text } = Typography;

interface StoreLoginForm {
  email: string;
  password: string;
}

export default function StoreLoginPage() {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm<StoreLoginForm>();
  const router = useRouter();
  const { login } = useAuth();

  const handleSubmit = async (values: StoreLoginForm) => {
    try {
      setLoading(true);
      
      const loginData: StoreLoginRequest = {
        email: values.email,
        password: values.password,
        role: 'store',
      };

      await login(loginData);
      
      message.success('ログインに成功しました');
      router.push('/store/dashboard');
    } catch (error) {
      message.error('ログインに失敗しました。メールアドレスとパスワードを確認してください。');
      console.error('Store login error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <Space direction="vertical" size="large" className="w-full text-center">
          <div>
            <ShopOutlined style={{ fontSize: '48px', color: '#52c41a', marginBottom: '16px' }} />
            <Title level={2} className="!mb-2">店舗ログイン</Title>
            <Text type="secondary">POS QR システム店舗管理者用ログイン</Text>
          </div>

          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            size="large"
            className="w-full"
          >
            <Form.Item
              name="email"
              label="メールアドレス"
              rules={[
                { required: true, message: 'メールアドレスを入力してください' },
                { type: 'email', message: '正しいメールアドレスを入力してください' },
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="store@example.com"
                autoComplete="email"
              />
            </Form.Item>

            <Form.Item
              name="password"
              label="パスワード"
              rules={[
                { required: true, message: 'パスワードを入力してください' },
                { min: 6, message: 'パスワードは6文字以上で入力してください' },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="パスワード"
                autoComplete="current-password"
              />
            </Form.Item>

            <Form.Item className="!mb-0">
              <Button
                variant="primary"
                size="large"
                htmlType="submit"
                loading={loading}
                fullWidth
              >
                ログイン
              </Button>
            </Form.Item>
          </Form>

          <div className="text-center">
            <Text type="secondary">
              システム管理者の方は{' '}
              <Button
                variant="link"
                size="small"
                onClick={() => router.push('/auth/admin-login')}
              >
                こちら
              </Button>
            </Text>
          </div>
        </Space>
      </Card>
    </div>
  );
}