import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';

import { api, API_ENDPOINTS } from '@/lib/api';

// Seat types
export interface Seat {
  id: string;
  storeId: string;
  number: string;
  name: string;
  description?: string;
  capacity: number;
  status: SeatStatus;
  qrCode: string;
  qrCodeUrl: string;
  position?: SeatPosition;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export type SeatStatus = 'available' | 'occupied' | 'reserved' | 'cleaning' | 'maintenance';

export interface SeatPosition {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface CreateSeatRequest {
  number: string;
  name: string;
  description?: string;
  capacity: number;
  position?: SeatPosition;
}

export interface UpdateSeatRequest {
  number?: string;
  name?: string;
  description?: string;
  capacity?: number;
  status?: SeatStatus;
  position?: SeatPosition;
  isActive?: boolean;
}

export interface SeatsListResponse {
  seats: Seat[];
  total: number;
  page: number;
  limit: number;
}

export interface SeatResponse {
  seat: Seat;
}

export interface QRCodeResponse {
  qrCode: string;
  qrCodeUrl: string;
  sessionUrl: string;
}

// Query keys
const QUERY_KEYS = {
  SEATS: 'seats',
  SEAT: 'seat',
  SEAT_QR: 'seat-qr',
} as const;

// List seats hook
export const useSeats = (params?: {
  page?: number;
  limit?: number;
  search?: string;
  status?: SeatStatus;
  isActive?: boolean;
}) => {
  return useQuery({
    queryKey: [QUERY_KEYS.SEATS, params],
    queryFn: async (): Promise<SeatsListResponse> => {
      const response = await api.get<SeatsListResponse>(API_ENDPOINTS.STORE.SEATS, {
        params,
      });
      return response.data.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes for real-time updates
  });
};

// Get single seat hook
export const useSeat = (seatId: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.SEAT, seatId],
    queryFn: async (): Promise<Seat> => {
      const response = await api.get<SeatResponse>(
        `${API_ENDPOINTS.STORE.SEATS}/${seatId}`
      );
      return response.data.data.seat;
    },
    enabled: !!seatId,
    staleTime: 2 * 60 * 1000,
  });
};

// Get seat QR code hook
export const useSeatQR = (seatId: string) => {
  return useQuery({
    queryKey: [QUERY_KEYS.SEAT_QR, seatId],
    queryFn: async (): Promise<QRCodeResponse> => {
      const response = await api.get<QRCodeResponse>(
        `${API_ENDPOINTS.STORE.SEATS}/${seatId}/qr`
      );
      return response.data.data;
    },
    enabled: !!seatId,
    staleTime: 10 * 60 * 1000, // 10 minutes for QR codes
  });
};

// Create seat mutation
export const useCreateSeat = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateSeatRequest): Promise<Seat> => {
      const response = await api.post<SeatResponse>(
        API_ENDPOINTS.STORE.SEATS,
        data
      );
      return response.data.data.seat;
    },
    onSuccess: (newSeat) => {
      // Invalidate seats list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.SEATS] });
      
      // Add new seat to cache
      queryClient.setQueryData([QUERY_KEYS.SEAT, newSeat.id], newSeat);
      
      message.success('座席が正常に作成されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '座席の作成に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Update seat mutation
export const useUpdateSeat = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      seatId,
      data,
    }: {
      seatId: string;
      data: UpdateSeatRequest;
    }): Promise<Seat> => {
      const response = await api.put<SeatResponse>(
        `${API_ENDPOINTS.STORE.SEATS}/${seatId}`,
        data
      );
      return response.data.data.seat;
    },
    onSuccess: (updatedSeat) => {
      // Update seat in cache
      queryClient.setQueryData([QUERY_KEYS.SEAT, updatedSeat.id], updatedSeat);
      
      // Invalidate seats list to reflect changes
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.SEATS] });
      
      message.success('座席情報が正常に更新されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '座席の更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Delete seat mutation
export const useDeleteSeat = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (seatId: string): Promise<void> => {
      await api.delete(`${API_ENDPOINTS.STORE.SEATS}/${seatId}`);
    },
    onSuccess: (_, seatId) => {
      // Remove seat from cache
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.SEAT, seatId] });
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.SEAT_QR, seatId] });
      
      // Invalidate seats list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.SEATS] });
      
      message.success('座席が正常に削除されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || '座席の削除に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Update seat status mutation
export const useUpdateSeatStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      seatId,
      status,
    }: {
      seatId: string;
      status: SeatStatus;
    }): Promise<Seat> => {
      const response = await api.patch<SeatResponse>(
        `${API_ENDPOINTS.STORE.SEATS}/${seatId}/status`,
        { status }
      );
      return response.data.data.seat;
    },
    onSuccess: (updatedSeat) => {
      // Update seat in cache
      queryClient.setQueryData([QUERY_KEYS.SEAT, updatedSeat.id], updatedSeat);
      
      // Invalidate seats list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.SEATS] });
      
      message.success('座席ステータスが更新されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || 'ステータスの更新に失敗しました';
      message.error(errorMessage);
    },
  });
};

// Regenerate QR code mutation
export const useRegenerateQR = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (seatId: string): Promise<QRCodeResponse> => {
      const response = await api.post<QRCodeResponse>(
        `${API_ENDPOINTS.STORE.SEATS}/${seatId}/qr/regenerate`
      );
      return response.data.data;
    },
    onSuccess: (qrData, seatId) => {
      // Update QR code in cache
      queryClient.setQueryData([QUERY_KEYS.SEAT_QR, seatId], qrData);
      
      // Update seat with new QR code
      queryClient.setQueryData([QUERY_KEYS.SEAT, seatId], (oldSeat: Seat | undefined) => {
        if (oldSeat) {
          return {
            ...oldSeat,
            qrCode: qrData.qrCode,
            qrCodeUrl: qrData.qrCodeUrl,
          };
        }
        return oldSeat;
      });
      
      message.success('QRコードが再生成されました');
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.message || 'QRコードの再生成に失敗しました';
      message.error(errorMessage);
    },
  });
};

export default {
  useSeats,
  useSeat,
  useSeatQR,
  useCreateSeat,
  useUpdateSeat,
  useDeleteSeat,
  useUpdateSeatStatus,
  useRegenerateQR,
};