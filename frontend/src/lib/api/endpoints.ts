// API endpoint definitions
export const API_ENDPOINTS = {
  // Authentication
  AUTH: {
    ADMIN_LOGIN: '/auth/admin/login',
    STORE_LOGIN: '/auth/store/login',
    REFRESH: '/auth/refresh',
    LOGOUT: '/auth/logout',
  },

  // Admin endpoints
  ADMIN: {
    DASHBOARD: '/admin/dashboard',
    STORES: '/admin/stores',
    STORE_BY_ID: (id: string) => `/admin/stores/${id}`,
    MANAGERS: '/admin/managers',
    MANAGER_BY_ID: (id: string) => `/admin/managers/${id}`,
  },

  // Store endpoints
  STORE: {
    DASHBOARD: '/store/dashboard',
    PROFILE: '/store/profile',
    SEATS: '/store/seats',
    SEAT_BY_ID: (id: string) => `/store/seats/${id}`,
    SEAT_QR: (id: string) => `/store/seats/${id}/qr`,
    ORDERS: '/store/orders',
    ORDER_BY_ID: (id: string) => `/store/orders/${id}`,
    MENU: '/store/menu',
    MENU_BY_ID: (id: string) => `/store/menu/${id}`,
    CATEGORIES: '/store/menu/categories',
  },

  // Customer order endpoints
  ORDER: {
    SESSION: (sessionId: string) => `/order/session/${sessionId}`,
    MENU: (sessionId: string) => `/order/session/${sessionId}/menu`,
    CART: (sessionId: string) => `/order/session/${sessionId}/cart`,
    PLACE_ORDER: (sessionId: string) => `/order/session/${sessionId}/place`,
    HISTORY: (sessionId: string) => `/order/session/${sessionId}/history`,
  },

  // Common endpoints
  COMMON: {
    HEALTH: '/health',
    VERSION: '/version',
  },
} as const;

// Helper function to build URLs with query parameters
export const buildUrl = (
  endpoint: string,
  params?: Record<string, string | number | boolean>
): string => {
  if (!params) return endpoint;

  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    searchParams.append(key, String(value));
  });

  return `${endpoint}?${searchParams.toString()}`;
};

// Export endpoint types for TypeScript
export type ApiEndpoints = typeof API_ENDPOINTS;