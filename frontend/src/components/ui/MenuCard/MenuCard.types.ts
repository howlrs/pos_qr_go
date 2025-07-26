import type { MenuItem } from '@/types/models';

export interface MenuCardProps {
  menuItem: MenuItem;
  quantity?: number;
  onQuantityChange?: (quantity: number) => void;
  onAddToCart?: (quantity: number) => void;
  loading?: boolean;
  disabled?: boolean;
}