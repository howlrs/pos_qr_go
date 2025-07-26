import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { CartDrawer } from '../index';
import type { Cart } from '@/types/models';

// Mock API hooks
jest.mock('@/hooks/api/useOrders', () => ({
  useUpdateCartItem: jest.fn(() => ({
    mutateAsync: jest.fn(),
    isPending: false,
  })),
  useRemoveFromCart: jest.fn(() => ({
    mutateAsync: jest.fn(),
    isPending: false,
  })),
  useClearCart: jest.fn(() => ({
    mutateAsync: jest.fn(),
    isPending: false,
  })),
}));

const mockCart: Cart = {
  id: 'cart1',
  sessionId: 'session1',
  items: [
    {
      id: 'item1',
      menuItemId: 'menu1',
      menuItem: {
        id: 'menu1',
        name: 'テストメニュー1',
        description: 'テスト用メニュー1',
        price: 1000,
        categoryId: 'cat1',
        category: {
          id: 'cat1',
          name: 'テストカテゴリ',
          displayOrder: 1,
          isActive: true,
        },
        isAvailable: true,
      },
      quantity: 2,
      unitPrice: 1000,
      totalPrice: 2000,
      addedAt: '2025-01-01T00:00:00Z',
    },
    {
      id: 'item2',
      menuItemId: 'menu2',
      menuItem: {
        id: 'menu2',
        name: 'テストメニュー2',
        description: 'テスト用メニュー2',
        price: 1500,
        categoryId: 'cat1',
        category: {
          id: 'cat1',
          name: 'テストカテゴリ',
          displayOrder: 1,
          isActive: true,
        },
        isAvailable: true,
      },
      quantity: 1,
      unitPrice: 1500,
      totalPrice: 1500,
      specialInstructions: '辛さ控えめで',
      addedAt: '2025-01-01T00:00:00Z',
    },
  ],
  totalItems: 3,
  totalAmount: 3500,
  updatedAt: '2025-01-01T00:00:00Z',
};

const emptyCart: Cart = {
  id: 'cart1',
  sessionId: 'session1',
  items: [],
  totalItems: 0,
  totalAmount: 0,
  updatedAt: '2025-01-01T00:00:00Z',
};

const TestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('CartDrawer', () => {
  const mockOnOrderConfirm = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('空のカート', () => {
    it('空のカート状態が正しく表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={emptyCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      expect(screen.getByText('カートが空です')).toBeInTheDocument();
      expect(screen.getByText('メニューから商品を選んでください')).toBeInTheDocument();
    });

    it('カートがundefinedの場合も空の状態が表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      expect(screen.getByText('カートが空です')).toBeInTheDocument();
    });
  });

  describe('カート内容表示', () => {
    it('カート内の商品が正しく表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      expect(screen.getByText('テストメニュー1')).toBeInTheDocument();
      expect(screen.getByText('テストメニュー2')).toBeInTheDocument();
      expect(screen.getByText('¥1,000 × 2')).toBeInTheDocument();
      expect(screen.getByText('¥1,500 × 1')).toBeInTheDocument();
    });

    it('特別な要望が表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      expect(screen.getByText('備考: 辛さ控えめで')).toBeInTheDocument();
    });

    it('合計金額が正しく表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      expect(screen.getByText('3点')).toBeInTheDocument();
      expect(screen.getByText('¥3,500')).toBeInTheDocument();
    });
  });

  describe('数量変更', () => {
    it('数量を変更できる', async () => {
      const user = userEvent.setup();
      const { useUpdateCartItem } = require('@/hooks/api/useOrders');
      const mockMutateAsync = jest.fn();
      useUpdateCartItem.mockReturnValue({
        mutateAsync: mockMutateAsync,
        isPending: false,
      });

      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      const quantityInput = screen.getAllByDisplayValue('2')[0];
      await user.clear(quantityInput);
      await user.type(quantityInput, '3');

      expect(mockMutateAsync).toHaveBeenCalledWith({ quantity: 3 });
    });
  });

  describe('商品削除', () => {
    it('商品を削除できる', async () => {
      const user = userEvent.setup();
      const { useRemoveFromCart } = require('@/hooks/api/useOrders');
      const mockMutateAsync = jest.fn();
      useRemoveFromCart.mockReturnValue({
        mutateAsync: mockMutateAsync,
        isPending: false,
      });

      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      const deleteButtons = screen.getAllByText('削除');
      await user.click(deleteButtons[0]);

      expect(mockMutateAsync).toHaveBeenCalled();
    });
  });

  describe('カートクリア', () => {
    it('カートを空にできる', async () => {
      const user = userEvent.setup();
      const { useClearCart } = require('@/hooks/api/useOrders');
      const mockMutateAsync = jest.fn();
      useClearCart.mockReturnValue({
        mutateAsync: mockMutateAsync,
        isPending: false,
      });

      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      const clearButton = screen.getByText('カートを空にする');
      await user.click(clearButton);

      expect(mockMutateAsync).toHaveBeenCalled();
    });
  });

  describe('注文確認', () => {
    it('注文内容確認ボタンをクリックできる', async () => {
      const user = userEvent.setup();
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      const confirmButton = screen.getByText(/注文内容を確認する/);
      await user.click(confirmButton);

      expect(mockOnOrderConfirm).toHaveBeenCalled();
    });

    it('カートが空の場合は注文確認ボタンが無効', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={emptyCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      // 空のカートの場合、注文確認ボタンは表示されない
      expect(screen.queryByText(/注文内容を確認する/)).not.toBeInTheDocument();
    });
  });

  describe('ローディング状態', () => {
    it('ローディング中の表示が正しく動作する', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
            loading={true}
          />
        </TestWrapper>
      );

      // ローディング状態でも基本的な表示は維持される
      expect(screen.getByText('テストメニュー1')).toBeInTheDocument();
    });
  });

  describe('価格フォーマット', () => {
    it('日本円形式で価格が表示される', () => {
      render(
        <TestWrapper>
          <CartDrawer
            sessionId="session1"
            cart={mockCart}
            onOrderConfirm={mockOnOrderConfirm}
          />
        </TestWrapper>
      );

      // 日本円形式（¥記号付き、3桁区切り）で表示されることを確認
      expect(screen.getByText('¥2,000')).toBeInTheDocument();
      expect(screen.getByText('¥1,500')).toBeInTheDocument();
      expect(screen.getByText('¥3,500')).toBeInTheDocument();
    });
  });
});