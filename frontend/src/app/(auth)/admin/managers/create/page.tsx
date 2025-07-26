'use client';

import { useRouter } from 'next/navigation';
import {
  Typography,
  Form,
  Input,
  Select,
  Space,
  Row,
  Col,
  Divider,
  Checkbox,
} from 'antd';
import { ArrowLeftOutlined, UserOutlined, MailOutlined, LockOutlined } from '@ant-design/icons';

import { Button, Card } from '@/components';
import { useCreateManager, CreateManagerRequest } from '@/hooks/api/useManagers';
import { PERMISSIONS, PERMISSION_GROUPS, permissionUtils } from '@/lib/auth/permissions';

const { Title, Text } = Typography;
const { Password } = Input;
const { Option } = Select;

interface ManagerFormData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
  permissions: string[];
}

export default function CreateManagerPage() {
  const router = useRouter();
  const [form] = Form.useForm<ManagerFormData>();
  const createManagerMutation = useCreateManager();

  const handleSubmit = async (values: ManagerFormData) => {
    try {
      const managerData: CreateManagerRequest = {
        name: values.name,
        email: values.email,
        password: values.password,
        permissions: values.permissions,
      };

      const newManager = await createManagerMutation.mutateAsync(managerData);
      router.push(`/admin/managers/${newManager.id}`);
    } catch (error) {
      console.error('Manager creation error:', error);
    }
  };

  const handleCancel = () => {
    router.back();
  };

  const handlePermissionPreset = (preset: 'admin_all' | 'admin_readonly') => {
    const permissions = preset === 'admin_all' 
      ? PERMISSION_GROUPS.ADMIN_ALL 
      : [PERMISSIONS.ADMIN.VIEW_ANALYTICS, PERMISSIONS.COMMON.VIEW_PROFILE];
    
    form.setFieldsValue({ permissions: [...permissions] });
  };

  // Group permissions by category
  const adminPermissions = [
    PERMISSIONS.ADMIN.MANAGE_STORES,
    PERMISSIONS.ADMIN.MANAGE_MANAGERS,
    PERMISSIONS.ADMIN.VIEW_ANALYTICS,
    PERMISSIONS.ADMIN.SYSTEM_SETTINGS,
  ];

  const commonPermissions = [
    PERMISSIONS.COMMON.VIEW_PROFILE,
    PERMISSIONS.COMMON.EDIT_PROFILE,
  ];

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
          <Title level={2}>新しい管理者を作成</Title>
          <Text type="secondary">
            新しい管理者の基本情報と権限を設定してください
          </Text>
        </div>

        {/* Form */}
        <Card>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{
              permissions: [PERMISSIONS.COMMON.VIEW_PROFILE],
            }}
          >
            {/* Basic Information */}
            <Title level={4}>基本情報</Title>
            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="name"
                  label="管理者名"
                  rules={[
                    { required: true, message: '管理者名を入力してください' },
                    { min: 2, message: '管理者名は2文字以上で入力してください' },
                  ]}
                >
                  <Input
                    prefix={<UserOutlined />}
                    placeholder="山田太郎"
                  />
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
                  <Input
                    prefix={<MailOutlined />}
                    placeholder="admin@example.com"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col xs={24} md={12}>
                <Form.Item
                  name="password"
                  label="パスワード"
                  rules={[
                    { required: true, message: 'パスワードを入力してください' },
                    { min: 8, message: 'パスワードは8文字以上で入力してください' },
                    {
                      pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
                      message: 'パスワードは大文字、小文字、数字を含む必要があります',
                    },
                  ]}
                >
                  <Password
                    prefix={<LockOutlined />}
                    placeholder="パスワード"
                  />
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  name="confirmPassword"
                  label="パスワード確認"
                  dependencies={['password']}
                  rules={[
                    { required: true, message: 'パスワードを再入力してください' },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue('password') === value) {
                          return Promise.resolve();
                        }
                        return Promise.reject(new Error('パスワードが一致しません'));
                      },
                    }),
                  ]}
                >
                  <Password
                    prefix={<LockOutlined />}
                    placeholder="パスワード確認"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Divider />

            {/* Permissions */}
            <Title level={4}>権限設定</Title>
            
            {/* Permission Presets */}
            <div className="mb-4">
              <Text strong>権限プリセット:</Text>
              <div className="mt-2">
                <Space>
                  <Button
                    size="small"
                    onClick={() => handlePermissionPreset('admin_all')}
                  >
                    全権限
                  </Button>
                  <Button
                    size="small"
                    onClick={() => handlePermissionPreset('admin_readonly')}
                  >
                    閲覧のみ
                  </Button>
                </Space>
              </div>
            </div>

            <Form.Item
              name="permissions"
              label="個別権限"
              rules={[
                { required: true, message: '少なくとも1つの権限を選択してください' },
              ]}
            >
              <Checkbox.Group className="w-full">
                <div>
                  <Text strong className="block mb-2">管理者権限</Text>
                  <Row gutter={[16, 8]}>
                    {adminPermissions.map((permission) => (
                      <Col xs={24} sm={12} key={permission}>
                        <Checkbox value={permission}>
                          {permissionUtils.getPermissionDescription(permission)}
                        </Checkbox>
                      </Col>
                    ))}
                  </Row>
                </div>
                
                <Divider />
                
                <div>
                  <Text strong className="block mb-2">共通権限</Text>
                  <Row gutter={[16, 8]}>
                    {commonPermissions.map((permission) => (
                      <Col xs={24} sm={12} key={permission}>
                        <Checkbox value={permission}>
                          {permissionUtils.getPermissionDescription(permission)}
                        </Checkbox>
                      </Col>
                    ))}
                  </Row>
                </div>
              </Checkbox.Group>
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
                  loading={createManagerMutation.isPending}
                >
                  管理者を作成
                </Button>
              </Col>
            </Row>
          </Form>
        </Card>
      </Space>
    </div>
  );
}