// User model types
export interface BaseUser {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface AdminUser extends BaseUser {
  role: 'admin';
  permissions: AdminPermission[];
}

export interface StoreUser extends BaseUser {
  role: 'store';
  storeId: string;
  permissions: StorePermission[];
}

export type User = AdminUser | StoreUser;

// Permission types
export type AdminPermission =
  | 'admin:stores:read'
  | 'admin:stores:write'
  | 'admin:stores:delete'
  | 'admin:managers:read'
  | 'admin:managers:write'
  | 'admin:managers:delete'
  | 'admin:dashboard:read';

export type StorePermission =
  | 'store:dashboard:read'
  | 'store:seats:read'
  | 'store:seats:write'
  | 'store:orders:read'
  | 'store:orders:write'
  | 'store:menu:read'
  | 'store:menu:write';

export type Permission = AdminPermission | StorePermission;

// User creation/update types
export interface CreateUserRequest {
  email: string;
  name: string;
  password: string;
  role: 'admin' | 'store';
  storeId?: string;
  permissions: Permission[];
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  permissions?: Permission[];
}

// User response types
export interface UserResponse {
  user: User;
  token?: string;
}

export interface UsersListResponse {
  users: User[];
  total: number;
  page: number;
  limit: number;
}