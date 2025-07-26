import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MenuCard } from '../index';
import type { MenuItem } from '@/types/models';

// Mock Ant Design Image component
jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  Image: ({ alt, src, className, preview }: any) => (
    <img alt={alt} src={src} className={className} data-preview={preview} />
  ),
}));

const mockMenuItem: MenuItem = {
  id: '1',
  name: 'テストメニュー',
  description: 'テスト用のメニューアイテムです',
  price: 1000,
  categoryId: 'cat1',
  category: {
    id: 'cat1',
    name: 'テストカテゴリ',
    displayOrder: 1,
    isActive: true,
  },
  imageUrl: 'https://example.com/test-image.jpg',
  isAvailable: true,
  allergens: ['卵', '小麦'],
  nutritionInfo: {
    calories: 300,
    protein: 15,
    carbs: 30,
    fat: 10,
  },
};

const mockUnavailableMenuItem: MenuItem = {
  ...mockMenuItem,
  id: '2',
  name: '売り切れメニュー',
  isAvailable: false,
};

describe('MenuCard', () => {
  const mockOnQuantityChange = jest.fn();
  const mockOnAddToCart = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('基本表示', () => {
    it('メニュー情報が正しく表示される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('テストメニュー')).toBeInTheDocument();
      expect(screen.getByText('テスト用のメニューアイテムです')).toBeInTheDocument();
      expect(screen.getByText('¥1,000')).toBeInTheDocument();
    });

    it('画像が正しく表示される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const image = screen.getByAltText('テストメニュー');
      expect(image).toBeInTheDocument();
      expect(image).toHaveAttribute('src', 'https://example.com/test-image.jpg');
    });

    it('アレルゲン情報が表示される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('アレルゲン:')).toBeInTheDocument();
      expect(screen.getByText('卵')).toBeInTheDocument();
      expect(screen.getByText('小麦')).toBeInTheDocument();
    });

    it('栄養情報が表示される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('カロリー: 300kcal')).toBeInTheDocument();
    });
  });

  describe('数量操作', () => {
    it('数量を増加できる', async () => {
      const user = userEvent.setup();
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={0}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const plusButton = screen.getByRole('button', { name: /plus/i });
      await user.click(plusButton);

      expect(mockOnQuantityChange).toHaveBeenCalledWith(1);
    });

    it('数量を減少できる', async () => {
      const user = userEvent.setup();
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={2}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const minusButton = screen.getByRole('button', { name: /minus/i });
      await user.click(minusButton);

      expect(mockOnQuantityChange).toHaveBeenCalledWith(1);
    });

    it('数量を直接入力できる', async () => {
      const user = userEvent.setup();
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={1}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const input = screen.getByDisplayValue('1');
      await user.clear(input);
      await user.type(input, '3');

      expect(mockOnQuantityChange).toHaveBeenCalledWith(3);
    });

    it('数量が0以下にならない', async () => {
      const user = userEvent.setup();
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={0}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const minusButton = screen.getByRole('button', { name: /minus/i });
      await user.click(minusButton);

      expect(mockOnQuantityChange).toHaveBeenCalledWith(0);
    });
  });

  describe('カート追加', () => {
    it('カートに追加できる', async () => {
      const user = userEvent.setup();
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={2}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const addButton = screen.getByText(/カートに追加/);
      await user.click(addButton);

      expect(mockOnAddToCart).toHaveBeenCalledWith(2);
    });

    it('数量が0の時はカート追加ボタンが無効', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={0}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const addButton = screen.getByText(/カートに追加/);
      expect(addButton).toBeDisabled();
    });

    it('合計金額が正しく表示される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={3}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('カートに追加 (¥3,000)')).toBeInTheDocument();
    });
  });

  describe('売り切れ状態', () => {
    it('売り切れメニューが正しく表示される', () => {
      render(
        <MenuCard
          menuItem={mockUnavailableMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('売り切れ')).toBeInTheDocument();
      expect(screen.getByText('売り切れ', { selector: 'button' })).toBeDisabled();
    });

    it('売り切れメニューは数量操作ができない', () => {
      render(
        <MenuCard
          menuItem={mockUnavailableMenuItem}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      const buttons = screen.queryAllByRole('button');
      const enabledButtons = buttons.filter(button => !button.hasAttribute('disabled'));
      
      // 売り切れボタンのみが表示され、他のボタンは無効化されている
      expect(enabledButtons).toHaveLength(0);
    });
  });

  describe('ローディング状態', () => {
    it('ローディング中はボタンが無効化される', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={1}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
          loading={true}
        />
      );

      const addButton = screen.getByText(/カートに追加/);
      expect(addButton).toBeDisabled();
    });
  });

  describe('無効化状態', () => {
    it('無効化時は全ての操作ができない', () => {
      render(
        <MenuCard
          menuItem={mockMenuItem}
          quantity={1}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
          disabled={true}
        />
      );

      const input = screen.getByDisplayValue('1');
      const addButton = screen.getByText(/カートに追加/);
      
      expect(input).toBeDisabled();
      expect(addButton).toBeDisabled();
    });
  });

  describe('画像なしの場合', () => {
    it('画像がない場合のプレースホルダーが表示される', () => {
      const menuItemWithoutImage = { ...mockMenuItem, imageUrl: undefined };
      render(
        <MenuCard
          menuItem={menuItemWithoutImage}
          onQuantityChange={mockOnQuantityChange}
          onAddToCart={mockOnAddToCart}
        />
      );

      expect(screen.getByText('画像なし')).toBeInTheDocument();
    });
  });
});