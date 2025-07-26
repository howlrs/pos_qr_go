import type { Cart } from '@/types/models';

export interface OrderConfirmationProps {
  visible: boolean;
  cart?: Cart;
  specialInstructions?: string;
  onSpecialInstructionsChange?: (value: string) => void;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}