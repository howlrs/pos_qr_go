// Authentication API types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    type: 'admin' | 'store';
  };
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  type: 'admin' | 'store';
  createdAt: string;
  updatedAt: string;
}
