import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api';
import { API_ENDPOINTS } from '@/lib/api/endpoints';
import type {
  OrderSessionResponse,
  MenuResponse,
  CartResponse,
  PlaceOrderResponse,
  OrderHistoryResponse,
  AddToCartRequest,
  UpdateCartItemRequest,
  PlaceOrderRequest,
} from '@/types/models';

// Query keys
export const ORDER_QUERY_KEYS = {
  session: (sessionId: string) => ['order', 'session', sessionId],
  menu: (sessionId: string) => ['order', 'menu', sessionId],
  cart: (sessionId: string) => ['order', 'cart', sessionId],
  history: (sessionId: string) => ['order', 'history', sessionId],
} as const;

// Session management
export const useOrderSession = (sessionId: string) => {
  return useQuery({
    queryKey: ORDER_QUERY_KEYS.session(sessionId),
    queryFn: async (): Promise<OrderSessionResponse> => {
      const response = await apiClient.get(API_ENDPOINTS.ORDER.SESSION(sessionId));
      return response.data;
    },
    enabled: !!sessionId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 3,
  });
};

// Menu management
export const useOrderMenu = (sessionId: string) => {
  return useQuery({
    queryKey: ORDER_QUERY_KEYS.menu(sessionId),
    queryFn: async (): Promise<MenuResponse> => {
      const response = await apiClient.get(API_ENDPOINTS.ORDER.MENU(sessionId));
      return response.data;
    },
    enabled: !!sessionId,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

// Cart management
export const useCart = (sessionId: string) => {
  return useQuery({
    queryKey: ORDER_QUERY_KEYS.cart(sessionId),
    queryFn: async (): Promise<CartResponse> => {
      const response = await apiClient.get(API_ENDPOINTS.ORDER.CART(sessionId));
      return response.data;
    },
    enabled: !!sessionId,
    refetchInterval: 30 * 1000, // Refetch every 30 seconds
  });
};

// Add item to cart
export const useAddToCart = (sessionId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: AddToCartRequest): Promise<CartResponse> => {
      const response = await apiClient.post(API_ENDPOINTS.ORDER.CART(sessionId), data);
      return response.data;
    },
    onSuccess: () => {
      // Invalidate cart query to refetch updated cart
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.cart(sessionId),
      });
    },
  });
};

// Update cart item
export const useUpdateCartItem = (sessionId: string, itemId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: UpdateCartItemRequest): Promise<CartResponse> => {
      const response = await apiClient.put(
        `${API_ENDPOINTS.ORDER.CART(sessionId)}/items/${itemId}`,
        data
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.cart(sessionId),
      });
    },
  });
};

// Remove item from cart
export const useRemoveFromCart = (sessionId: string, itemId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (): Promise<CartResponse> => {
      const response = await apiClient.delete(
        `${API_ENDPOINTS.ORDER.CART(sessionId)}/items/${itemId}`
      );
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.cart(sessionId),
      });
    },
  });
};

// Clear cart
export const useClearCart = (sessionId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (): Promise<CartResponse> => {
      const response = await apiClient.delete(API_ENDPOINTS.ORDER.CART(sessionId));
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.cart(sessionId),
      });
    },
  });
};

// Place order
export const usePlaceOrder = (sessionId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: PlaceOrderRequest): Promise<PlaceOrderResponse> => {
      const response = await apiClient.post(API_ENDPOINTS.ORDER.PLACE_ORDER(sessionId), data);
      return response.data;
    },
    onSuccess: () => {
      // Clear cart and refresh history after successful order
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.cart(sessionId),
      });
      queryClient.invalidateQueries({
        queryKey: ORDER_QUERY_KEYS.history(sessionId),
      });
    },
  });
};

// Order history
export const useOrderHistory = (sessionId: string) => {
  return useQuery({
    queryKey: ORDER_QUERY_KEYS.history(sessionId),
    queryFn: async (): Promise<OrderHistoryResponse> => {
      const response = await apiClient.get(API_ENDPOINTS.ORDER.HISTORY(sessionId));
      return response.data;
    },
    enabled: !!sessionId,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

// Combined hook for order page
export const useOrderPage = (sessionId: string) => {
  const sessionQuery = useOrderSession(sessionId);
  const menuQuery = useOrderMenu(sessionId);
  const cartQuery = useCart(sessionId);

  return {
    session: sessionQuery,
    menu: menuQuery,
    cart: cartQuery,
    isLoading: sessionQuery.isLoading || menuQuery.isLoading || cartQuery.isLoading,
    isError: sessionQuery.isError || menuQuery.isError || cartQuery.isError,
    error: sessionQuery.error || menuQuery.error || cartQuery.error,
  };
};