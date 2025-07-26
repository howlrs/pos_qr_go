import type { Order } from '@/types/models';

export interface OrderSuccessProps {
  visible: boolean;
  order?: Order;
  estimatedWaitTime?: number;
  onClose: () => void;
  onContinueOrdering: () => void;
  onViewHistory: () => void;
}