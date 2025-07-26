// Application routes configuration
export const routes = {
  // Public routes
  home: '/',

  // Auth routes
  auth: {
    adminLogin: '/auth/admin-login',
    storeLogin: '/auth/store-login',
  },

  // Admin routes
  admin: {
    dashboard: '/admin/dashboard',
    stores: '/admin/stores',
    storeDetail: (id: string) => `/admin/stores/${id}`,
    storeCreate: '/admin/stores/create',
    storeEdit: (id: string) => `/admin/stores/${id}/edit`,
    managers: '/admin/managers',
    managerDetail: (id: string) => `/admin/managers/${id}`,
    managerCreate: '/admin/managers/create',
  },

  // Store routes
  store: {
    dashboard: '/store/dashboard',
    seats: '/store/seats',
    seatDetail: (id: string) => `/store/seats/${id}`,
    seatCreate: '/store/seats/create',
    seatQR: (id: string) => `/store/seats/${id}/qr`,
    orders: '/store/orders',
    orderDetail: (id: string) => `/store/orders/${id}`,
    menu: '/store/menu',
    menuDetail: (id: string) => `/store/menu/${id}`,
    menuCreate: '/store/menu/create',
    menuEdit: (id: string) => `/store/menu/${id}/edit`,
    categories: '/store/menu/categories',
  },

  // Customer routes
  order: {
    session: (sessionId: string) => `/order/${sessionId}`,
    menu: (sessionId: string) => `/order/${sessionId}/menu`,
    cart: (sessionId: string) => `/order/${sessionId}/cart`,
    history: (sessionId: string) => `/order/${sessionId}/history`,
  },
} as const;
