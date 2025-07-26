// Order related types
export interface OrderSession {
  id: string;
  seatId: string;
  storeId: string;
  status: 'active' | 'completed' | 'expired';
  createdAt: string;
  expiresAt: string;
  seat: {
    id: string;
    number: string;
    name: string;
  };
  store: {
    id: string;
    name: string;
    address: string;
    phone: string;
  };
}

export interface MenuItem {
  id: string;
  name: string;
  description: string;
  price: number;
  categoryId: string;
  category: MenuCategory;
  imageUrl?: string;
  isAvailable: boolean;
  allergens?: string[];
  nutritionInfo?: {
    calories?: number;
    protein?: number;
    carbs?: number;
    fat?: number;
  };
}

export interface MenuCategory {
  id: string;
  name: string;
  description?: string;
  displayOrder: number;
  isActive: boolean;
}

export interface CartItem {
  id: string;
  menuItemId: string;
  menuItem: MenuItem;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  specialInstructions?: string;
  addedAt: string;
}

export interface Cart {
  id: string;
  sessionId: string;
  items: CartItem[];
  totalItems: number;
  totalAmount: number;
  updatedAt: string;
}

export interface Order {
  id: string;
  sessionId: string;
  orderNumber: string;
  status: 'pending' | 'confirmed' | 'preparing' | 'ready' | 'served' | 'cancelled';
  items: OrderItem[];
  totalAmount: number;
  specialInstructions?: string;
  placedAt: string;
  estimatedReadyAt?: string;
  completedAt?: string;
}

export interface OrderItem {
  id: string;
  orderId: string;
  menuItemId: string;
  menuItem: MenuItem;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  specialInstructions?: string;
}

// API Request/Response types
export interface AddToCartRequest {
  menuItemId: string;
  quantity: number;
  specialInstructions?: string;
}

export interface UpdateCartItemRequest {
  quantity: number;
  specialInstructions?: string;
}

export interface PlaceOrderRequest {
  specialInstructions?: string;
}

export interface OrderSessionResponse {
  session: OrderSession;
  menu: {
    categories: MenuCategory[];
    items: MenuItem[];
  };
  cart: Cart;
}

export interface MenuResponse {
  categories: MenuCategory[];
  items: MenuItem[];
}

export interface CartResponse {
  cart: Cart;
}

export interface PlaceOrderResponse {
  order: Order;
  estimatedWaitTime: number; // in minutes
}

export interface OrderHistoryResponse {
  orders: Order[];
  totalCount: number;
}

// Order status display helpers
export const ORDER_STATUS_LABELS = {
  pending: '注文受付中',
  confirmed: '注文確認済み',
  preparing: '調理中',
  ready: '準備完了',
  served: '提供済み',
  cancelled: 'キャンセル',
} as const;

export const ORDER_STATUS_COLORS = {
  pending: 'orange',
  confirmed: 'blue',
  preparing: 'purple',
  ready: 'green',
  served: 'gray',
  cancelled: 'red',
} as const;

export type OrderStatus = keyof typeof ORDER_STATUS_LABELS;