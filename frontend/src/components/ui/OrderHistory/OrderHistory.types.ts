import type { Order } from '@/types/models';

export interface OrderHistoryProps {
  orders?: Order[];
  loading?: boolean;
}