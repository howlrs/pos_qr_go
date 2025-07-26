// Store model types
export interface Store {
  id: string;
  name: string;
  description?: string;
  address: string;
  phone: string;
  email: string;
  isActive: boolean;
  settings: StoreSettings;
  createdAt: string;
  updatedAt: string;
}

export interface StoreSettings {
  timezone: string;
  currency: string;
  language: string;
  orderTimeout: number; // minutes
  maxSeats: number;
  features: StoreFeature[];
}

export type StoreFeature =
  | 'qr_ordering'
  | 'table_service'
  | 'takeaway'
  | 'delivery'
  | 'payment_integration';

// Store creation/update types
export interface CreateStoreRequest {
  name: string;
  description?: string;
  address: string;
  phone: string;
  email: string;
  settings?: Partial<StoreSettings>;
}

export interface UpdateStoreRequest {
  name?: string;
  description?: string;
  address?: string;
  phone?: string;
  email?: string;
  isActive?: boolean;
  settings?: Partial<StoreSettings>;
}

// Store response types
export interface StoreResponse {
  store: Store;
}

export interface StoresListResponse {
  stores: Store[];
  total: number;
  page: number;
  limit: number;
}

// Store statistics
export interface StoreStats {
  totalOrders: number;
  totalRevenue: number;
  activeSeats: number;
  averageOrderValue: number;
  ordersToday: number;
  revenueToday: number;
}