import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';

import { api, API_ENDPOINTS } from '@/lib/api';
import {
  Store,
  StoresListResponse,
  StoreResponse,
  CreateStoreRequest,
  UpdateStoreRequest,
  StoreStats,
} from '@/types';

// Query keys
const QUERY_KEYS = {
  STORES: 'stores',
  STORE: 'store',
  STORE_STATS: 'store-stats',
} as const;

// List stores hook
export const useStores = (params?: {
  page?: number;
  limit?: number;
  search?: string;
  isActive?: boolean;
}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.STORES, params],
    queryFn: async (): Promise<StoresListResponse> => {
      const response = await api.get<StoresListResponse>(API_ENDPOINTS.ADMIN.STORES, {
        params,
      });
      return response.data.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Get single store hook
export const useStore = (storeId: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.STORE, storeId],
    queryFn: async (): Promise<Store> => {
      const response = await api.get<StoreResponse>(
        `${API_ENDPOINTS.ADMIN.STORES}/${storeId}`
      );
      return response.data.data.store;
    },
    enabled: !!storeId,
    staleTime: 5 * 60 * 1000,
  });
};

// Get store statistics hook
export const useStoreStats = (storeId: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.STORE_STATS, storeId],
    queryFn: async (): Promise<StoreStats> => {
      const response = await api.get<{ stats: StoreStats }>(
        `${API_ENDPOINTS.ADMIN.STORES}/${storeId}/stats`
      );
      return response.data.data.stats;
    },
    enabled: !!storeId,
    staleTime: 2 * 60 * 1000, // 2 minutes for stats
  });
};

// Create store mutation
export const useCreateStore = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateStoreRequest): Promise<Store> => {
      const response = await api.post<StoreResponse>(
        API_ENDPOINTS.ADMIN.STORES,
        data
      );
      return response.data.data.store;
    },
    onSuccess: (newStore) => {
      // Invalidate stores list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STORES] });
      
      // Add new store to cache
      queryClient.setQueryData([QUERY_KEYS.STORE, newStore.id], newStore);
      
      message.success('店舗が正常に作成されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '店舗の作成に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Update store mutation
export const useUpdateStore = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      storeId,
      data,
    }: {
      storeId: string;
      data: UpdateStoreRequest;
    }): Promise<Store> => {
      const response = await api.put<StoreResponse>(
        `${API_ENDPOINTS.ADMIN.STORES}/${storeId}`,
        data
      );
      return response.data.data.store;
    },
    onSuccess: (updatedStore) => {
      // Update store in cache
      queryClient.setQueryData([QUERY_KEYS.STORE, updatedStore.id], updatedStore);
      
      // Invalidate stores list to reflect changes
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STORES] });
      
      message.success('店舗情報が正常に更新されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '店舗の更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Delete store mutation
export const useDeleteStore = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (storeId: string): Promise<void> => {
      await api.delete(`${API_ENDPOINTS.ADMIN.STORES}/${storeId}`);
    },
    onSuccess: (_, storeId) => {
      // Remove store from cache
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.STORE, storeId] });
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.STORE_STATS, storeId] });
      
      // Invalidate stores list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STORES] });
      
      message.success('店舗が正常に削除されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '店舗の削除に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Toggle store status mutation
export const useToggleStoreStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      storeId,
      isActive,
    }: {
      storeId: string;
      isActive: boolean;
    }): Promise<Store> => {
      const response = await api.patch<StoreResponse>(
        `${API_ENDPOINTS.ADMIN.STORES}/${storeId}/status`,
        { isActive }
      );
      return response.data.data.store;
    },
    onSuccess: (updatedStore) => {
      // Update store in cache
      queryClient.setQueryData([QUERY_KEYS.STORE, updatedStore.id], updatedStore);
      
      // Invalidate stores list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.STORES] });
      
      const statusText = updatedStore.isActive ? '有効' : '無効';
      message.success(`店舗が${statusText}に設定されました`);
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || 'ステータスの更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

export default {
  useStores,
  useStore,
  useStoreStats,
  useCreateStore,
  useUpdateStore,
  useDeleteStore,
  useToggleStoreStatus,
};