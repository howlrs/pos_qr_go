import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import {
  useOrderSession,
  useOrderMenu,
  useCart,
  useAddToCart,
  useOrderPage,
} from '../useOrders';
import { apiClient } from '@/lib/api';

// Mock API client
jest.mock('@/lib/api', () => ({
  apiClient: {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
  },
}));

const mockApiClient = apiClient as jest.Mocked<typeof apiClient>;

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        staleTime: 0,
        gcTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });

  return ({ children }: { children: ReactNode }) => 
    React.createElement(QueryClientProvider, { client: queryClient }, children);
};

describe('useOrders hooks', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('useOrderSession', () => {
    const mockSessionResponse = {
      data: {
        session: {
          id: 'session1',
          seatId: 'seat1',
          storeId: 'store1',
          status: 'active',
          createdAt: '2025-01-01T00:00:00Z',
          expiresAt: '2025-01-01T02:00:00Z',
          seat: {
            id: 'seat1',
            number: '1',
            name: 'テーブル1',
          },
          store: {
            id: 'store1',
            name: 'テストレストラン',
            address: '東京都渋谷区',
            phone: '03-1234-5678',
          },
        },
        menu: {
          categories: [],
          items: [],
        },
        cart: {
          id: 'cart1',
          sessionId: 'session1',
          items: [],
          totalItems: 0,
          totalAmount: 0,
          updatedAt: '2025-01-01T00:00:00Z',
        },
      },
    };

    it('セッション情報を正しく取得する', async () => {
      mockApiClient.get.mockResolvedValueOnce(mockSessionResponse);

      const { result } = renderHook(
        () => useOrderSession('session1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/order/session/session1');
      expect(result.current.data).toEqual(mockSessionResponse.data);
    });

    it('sessionIdが空の場合はクエリが無効化される', () => {
      const { result } = renderHook(
        () => useOrderSession(''),
        { wrapper: createWrapper() }
      );

      expect(result.current.fetchStatus).toBe('idle');
      expect(mockApiClient.get).not.toHaveBeenCalled();
    });

    it('エラーが発生した場合の処理', async () => {
      const mockError = new Error('Session not found');
      mockApiClient.get.mockRejectedValueOnce(mockError);

      const { result } = renderHook(
        () => useOrderSession('invalid-session'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(mockError);
    });
  });

  describe('useOrderMenu', () => {
    const mockMenuResponse = {
      data: {
        categories: [
          {
            id: 'cat1',
            name: 'メイン',
            displayOrder: 1,
            isActive: true,
          },
        ],
        items: [
          {
            id: 'item1',
            name: 'テストメニュー',
            description: 'テスト用メニュー',
            price: 1000,
            categoryId: 'cat1',
            isAvailable: true,
          },
        ],
      },
    };

    it('メニュー情報を正しく取得する', async () => {
      mockApiClient.get.mockResolvedValueOnce(mockMenuResponse);

      const { result } = renderHook(
        () => useOrderMenu('session1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/order/session/session1/menu');
      expect(result.current.data).toEqual(mockMenuResponse.data);
    });
  });

  describe('useCart', () => {
    const mockCartResponse = {
      data: {
        cart: {
          id: 'cart1',
          sessionId: 'session1',
          items: [
            {
              id: 'item1',
              menuItemId: 'menu1',
              quantity: 2,
              unitPrice: 1000,
              totalPrice: 2000,
            },
          ],
          totalItems: 2,
          totalAmount: 2000,
          updatedAt: '2025-01-01T00:00:00Z',
        },
      },
    };

    it('カート情報を正しく取得する', async () => {
      mockApiClient.get.mockResolvedValueOnce(mockCartResponse);

      const { result } = renderHook(
        () => useCart('session1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/order/session/session1/cart');
      expect(result.current.data).toEqual(mockCartResponse.data);
    });
  });

  describe('useAddToCart', () => {
    const mockAddToCartResponse = {
      data: {
        cart: {
          id: 'cart1',
          sessionId: 'session1',
          items: [
            {
              id: 'item1',
              menuItemId: 'menu1',
              quantity: 1,
              unitPrice: 1000,
              totalPrice: 1000,
            },
          ],
          totalItems: 1,
          totalAmount: 1000,
          updatedAt: '2025-01-01T00:00:00Z',
        },
      },
    };

    it('カートに商品を追加できる', async () => {
      mockApiClient.post.mockResolvedValueOnce(mockAddToCartResponse);

      const { result } = renderHook(
        () => useAddToCart('session1'),
        { wrapper: createWrapper() }
      );

      const addToCartData = {
        menuItemId: 'menu1',
        quantity: 1,
      };

      await result.current.mutateAsync(addToCartData);

      expect(mockApiClient.post).toHaveBeenCalledWith(
        '/order/session/session1/cart',
        addToCartData
      );
    });

    it('エラーが発生した場合の処理', async () => {
      const mockError = new Error('Failed to add to cart');
      mockApiClient.post.mockRejectedValueOnce(mockError);

      const { result } = renderHook(
        () => useAddToCart('session1'),
        { wrapper: createWrapper() }
      );

      const addToCartData = {
        menuItemId: 'menu1',
        quantity: 1,
      };

      await expect(result.current.mutateAsync(addToCartData)).rejects.toThrow(
        'Failed to add to cart'
      );
    });
  });

  describe('useOrderPage', () => {
    const mockSessionResponse = {
      data: {
        session: {
          id: 'session1',
          seatId: 'seat1',
          storeId: 'store1',
          status: 'active',
        },
      },
    };

    const mockMenuResponse = {
      data: {
        categories: [],
        items: [],
      },
    };

    const mockCartResponse = {
      data: {
        cart: {
          id: 'cart1',
          sessionId: 'session1',
          items: [],
          totalItems: 0,
          totalAmount: 0,
        },
      },
    };

    it('複数のクエリを統合して管理する', async () => {
      mockApiClient.get
        .mockResolvedValueOnce(mockSessionResponse)
        .mockResolvedValueOnce(mockMenuResponse)
        .mockResolvedValueOnce(mockCartResponse);

      const { result } = renderHook(
        () => useOrderPage('session1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.session.isSuccess).toBe(true);
        expect(result.current.menu.isSuccess).toBe(true);
        expect(result.current.cart.isSuccess).toBe(true);
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.isError).toBe(false);
    });

    it('いずれかのクエリでエラーが発生した場合', async () => {
      const mockError = new Error('Session error');
      mockApiClient.get
        .mockRejectedValueOnce(mockError)
        .mockResolvedValueOnce(mockMenuResponse)
        .mockResolvedValueOnce(mockCartResponse);

      const { result } = renderHook(
        () => useOrderPage('session1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.session.isError).toBe(true);
      });

      expect(result.current.isError).toBe(true);
      expect(result.current.error).toEqual(mockError);
    });

    it('ローディング状態が正しく管理される', () => {
      const { result } = renderHook(
        () => useOrderPage('session1'),
        { wrapper: createWrapper() }
      );

      // 初期状態ではローディング中
      expect(result.current.isLoading).toBe(true);
    });
  });

  describe('クエリキー', () => {
    it('正しいクエリキーが生成される', () => {
      const { ORDER_QUERY_KEYS } = require('../useOrders');

      expect(ORDER_QUERY_KEYS.session('session1')).toEqual(['order', 'session', 'session1']);
      expect(ORDER_QUERY_KEYS.menu('session1')).toEqual(['order', 'menu', 'session1']);
      expect(ORDER_QUERY_KEYS.cart('session1')).toEqual(['order', 'cart', 'session1']);
      expect(ORDER_QUERY_KEYS.history('session1')).toEqual(['order', 'history', 'session1']);
    });
  });
});