import type { Order } from '@/types/models';

export interface OrderStatusProps {
  order?: Order;
  onRefresh?: () => void;
  refreshing?: boolean;
}