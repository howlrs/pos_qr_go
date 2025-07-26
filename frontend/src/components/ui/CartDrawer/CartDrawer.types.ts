import type { Cart } from '@/types/models';

export interface CartDrawerProps {
  sessionId: string;
  cart?: Cart;
  onOrderConfirm?: () => void;
  loading?: boolean;
}