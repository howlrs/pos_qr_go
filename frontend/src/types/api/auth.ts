// Authentication API types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AdminLoginRequest extends LoginCredentials {
  role: 'admin';
}

export interface StoreLoginRequest extends LoginCredentials {
  role: 'store';
}

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: 'admin' | 'store';
  };
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'store';
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface RefreshTokenResponse {
  token: string;
  refreshToken: string;
}
