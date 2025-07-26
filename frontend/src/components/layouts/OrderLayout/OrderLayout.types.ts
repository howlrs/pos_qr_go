import { ReactNode } from 'react';
import type { OrderSession } from '@/types/models';

export interface OrderLayoutProps {
  children: ReactNode;
  sessionId: string;
  session?: OrderSession;
  onCartClick?: () => void;
  onHistoryClick?: () => void;
  onStoreInfoClick?: () => void;
  showCartDrawer?: boolean;
  onCartDrawerClose?: () => void;
  cartContent?: ReactNode;
}