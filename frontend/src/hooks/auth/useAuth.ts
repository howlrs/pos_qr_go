import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';

import { useAuth as useAuthStore } from '@/store/auth/authStore';
import { api, API_ENDPOINTS } from '@/lib/api';
import { jwt } from '@/lib/auth/jwt';

// Login request types
export interface AdminLoginRequest {
  email: string;
  password: string;
}

export interface StoreLoginRequest {
  email: string;
  password: string;
  storeId?: string;
}

// Login response type
export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: 'admin' | 'store';
    storeId?: string;
    permissions: string[];
  };
}

// Auth hook
export const useAuth = () => {
  const authStore = useAuthStore();
  const queryClient = useQueryClient();

  // Admin login mutation
  const adminLogin = useMutation({
    mutationFn: async (credentials: AdminLoginRequest): Promise<LoginResponse> => {
      const response = await api.post<LoginResponse>(
        API_ENDPOINTS.AUTH.ADMIN_LOGIN,
        credentials
      );
      return response.data.data;
    },
    onSuccess: (data) => {
      const { token, user } = data;
      
      // Validate token
      if (!jwt.isValid(token)) {
        throw new Error('Invalid token received');
      }

      // Login user
      authStore.login(token, user);
      
      message.success('ログインしました');
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'ログインに失敗しました';
      authStore.setError(errorMessage);
      message.error(errorMessage);
    },
  });

  // Store login mutation
  const storeLogin = useMutation({
    mutationFn: async (credentials: StoreLoginRequest): Promise<LoginResponse> => {
      const response = await api.post<LoginResponse>(
        API_ENDPOINTS.AUTH.STORE_LOGIN,
        credentials
      );
      return response.data.data;
    },
    onSuccess: (data) => {
      const { token, user } = data;
      
      // Validate token
      if (!jwt.isValid(token)) {
        throw new Error('Invalid token received');
      }

      // Login user
      authStore.login(token, user);
      
      message.success('ログインしました');
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'ログインに失敗しました';
      authStore.setError(errorMessage);
      message.error(errorMessage);
    },
  });

  // Logout mutation
  const logout = useMutation({
    mutationFn: async () => {
      try {
        await api.post(API_ENDPOINTS.AUTH.LOGOUT);
      } catch (error) {
        // Continue with logout even if API call fails
        console.warn('Logout API call failed:', error);
      }
    },
    onSuccess: () => {
      // Clear auth state
      authStore.logout();
      
      // Clear all cached queries
      queryClient.clear();
      
      message.success('ログアウトしました');
    },
    onError: (error: any) => {
      // Still logout locally even if API call fails
      authStore.logout();
      queryClient.clear();
      
      console.error('Logout error:', error);
      message.warning('ログアウトしました（サーバーエラーが発生しました）');
    },
  });

  // Token refresh mutation
  const refreshToken = useMutation({
    mutationFn: async (): Promise<LoginResponse> => {
      const response = await api.post<LoginResponse>(API_ENDPOINTS.AUTH.REFRESH);
      return response.data.data;
    },
    onSuccess: (data) => {
      const { token, user } = data;
      
      // Validate token
      if (!jwt.isValid(token)) {
        throw new Error('Invalid refresh token received');
      }

      // Update auth state
      authStore.login(token, user);
    },
    onError: (error: any) => {
      console.error('Token refresh failed:', error);
      // Force logout on refresh failure
      authStore.logout();
      queryClient.clear();
    },
  });

  // Unified login function
  const login = async (credentials: (AdminLoginRequest & { role: 'admin' }) | (StoreLoginRequest & { role: 'store' })) => {
    if (credentials.role === 'admin') {
      return adminLogin.mutateAsync(credentials);
    } else {
      return storeLogin.mutateAsync(credentials);
    }
  };

  return {
    // State
    ...authStore,
    
    // Unified login
    login,
    
    // Mutations
    adminLogin: {
      mutate: adminLogin.mutate,
      mutateAsync: adminLogin.mutateAsync,
      isLoading: adminLogin.isPending,
      error: adminLogin.error,
    },
    
    storeLogin: {
      mutate: storeLogin.mutate,
      mutateAsync: storeLogin.mutateAsync,
      isLoading: storeLogin.isPending,
      error: storeLogin.error,
    },
    
    logout: {
      mutate: logout.mutate,
      mutateAsync: logout.mutateAsync,
      isLoading: logout.isPending,
    },
    
    refreshToken: {
      mutate: refreshToken.mutate,
      mutateAsync: refreshToken.mutateAsync,
      isLoading: refreshToken.isPending,
    },

    // Utilities
    checkTokenExpiration: () => {
      const { token } = authStore;
      if (!token) return false;
      
      if (jwt.willExpireSoon(token, 5)) {
        // Auto refresh if token expires in 5 minutes
        refreshToken.mutate();
        return true;
      }
      
      return false;
    },
  };
};

export default useAuth;