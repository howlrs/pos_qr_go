import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';

import { api, API_ENDPOINTS } from '@/lib/api';
import { AuthUser } from '@/types';

// Manager types
export interface Manager extends AuthUser {
  isActive: boolean;
  lastLoginAt?: string;
  createdBy: string;
}

export interface CreateManagerRequest {
  name: string;
  email: string;
  password: string;
  permissions: string[];
}

export interface UpdateManagerRequest {
  name?: string;
  email?: string;
  password?: string;
  permissions?: string[];
  isActive?: boolean;
}

export interface ManagersListResponse {
  managers: Manager[];
  total: number;
  page: number;
  limit: number;
}

export interface ManagerResponse {
  manager: Manager;
}

// Query keys
const QUERY_KEYS = {
  MANAGERS: 'managers',
  MANAGER: 'manager',
} as const;

// List managers hook
export const useManagers = (params?: {
  page?: number;
  limit?: number;
  search?: string;
  isActive?: boolean;
}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.MANAGERS, params],
    queryFn: async (): Promise<ManagersListResponse> => {
      const response = await api.get<ManagersListResponse>(API_ENDPOINTS.ADMIN.MANAGERS, {
        params,
      });
      return response.data.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Get single manager hook
export const useManager = (managerId: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.MANAGER, managerId],
    queryFn: async (): Promise<Manager> => {
      const response = await api.get<ManagerResponse>(
        `${API_ENDPOINTS.ADMIN.MANAGERS}/${managerId}`
      );
      return response.data.data.manager;
    },
    enabled: !!managerId,
    staleTime: 5 * 60 * 1000,
  });
};

// Create manager mutation
export const useCreateManager = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateManagerRequest): Promise<Manager> => {
      const response = await api.post<ManagerResponse>(
        API_ENDPOINTS.ADMIN.MANAGERS,
        data
      );
      return response.data.data.manager;
    },
    onSuccess: (newManager) => {
      // Invalidate managers list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MANAGERS] });
      
      // Add new manager to cache
      queryClient.setQueryData([QUERY_KEYS.MANAGER, newManager.id], newManager);
      
      message.success('管理者が正常に作成されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '管理者の作成に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Update manager mutation
export const useUpdateManager = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      managerId,
      data,
    }: {
      managerId: string;
      data: UpdateManagerRequest;
    }): Promise<Manager> => {
      const response = await api.put<ManagerResponse>(
        `${API_ENDPOINTS.ADMIN.MANAGERS}/${managerId}`,
        data
      );
      return response.data.data.manager;
    },
    onSuccess: (updatedManager) => {
      // Update manager in cache
      queryClient.setQueryData([QUERY_KEYS.MANAGER, updatedManager.id], updatedManager);
      
      // Invalidate managers list to reflect changes
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MANAGERS] });
      
      message.success('管理者情報が正常に更新されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '管理者の更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Delete manager mutation
export const useDeleteManager = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (managerId: string): Promise<void> => {
      await api.delete(`${API_ENDPOINTS.ADMIN.MANAGERS}/${managerId}`);
    },
    onSuccess: (_, managerId) => {
      // Remove manager from cache
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.MANAGER, managerId] });
      
      // Invalidate managers list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MANAGERS] });
      
      message.success('管理者が正常に削除されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '管理者の削除に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Toggle manager status mutation
export const useToggleManagerStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      managerId,
      isActive,
    }: {
      managerId: string;
      isActive: boolean;
    }): Promise<Manager> => {
      const response = await api.patch<ManagerResponse>(
        `${API_ENDPOINTS.ADMIN.MANAGERS}/${managerId}/status`,
        { isActive }
      );
      return response.data.data.manager;
    },
    onSuccess: (updatedManager) => {
      // Update manager in cache
      queryClient.setQueryData([QUERY_KEYS.MANAGER, updatedManager.id], updatedManager);
      
      // Invalidate managers list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MANAGERS] });
      
      const statusText = updatedManager.isActive ? '有効' : '無効';
      message.success(`管理者が${statusText}に設定されました`);
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || 'ステータスの更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

export default {
  useManagers,
  useManager,
  useCreateManager,
  useUpdateManager,
  useDeleteManager,
  useToggleManagerStatus,
};