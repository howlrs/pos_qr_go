// Permission constants
export const PERMISSIONS = {
  // Admin permissions
  ADMIN: {
    MANAGE_STORES: 'admin:manage_stores',
    MANAGE_MANAGERS: 'admin:manage_managers',
    VIEW_ANALYTICS: 'admin:view_analytics',
    SYSTEM_SETTINGS: 'admin:system_settings',
  },
  
  // Store permissions
  STORE: {
    MANAGE_SEATS: 'store:manage_seats',
    MANAGE_MENU: 'store:manage_menu',
    VIEW_ORDERS: 'store:view_orders',
    MANAGE_ORDERS: 'store:manage_orders',
    VIEW_ANALYTICS: 'store:view_analytics',
  },
  
  // Common permissions
  COMMON: {
    VIEW_PROFILE: 'common:view_profile',
    EDIT_PROFILE: 'common:edit_profile',
  },
} as const;

// Permission groups
export const PERMISSION_GROUPS = {
  ADMIN_ALL: [
    PERMISSIONS.ADMIN.MANAGE_STORES,
    PERMISSIONS.ADMIN.MANAGE_MANAGERS,
    PERMISSIONS.ADMIN.VIEW_ANALYTICS,
    PERMISSIONS.ADMIN.SYSTEM_SETTINGS,
    PERMISSIONS.COMMON.VIEW_PROFILE,
    PERMISSIONS.COMMON.EDIT_PROFILE,
  ],
  
  STORE_ALL: [
    PERMISSIONS.STORE.MANAGE_SEATS,
    PERMISSIONS.STORE.MANAGE_MENU,
    PERMISSIONS.STORE.VIEW_ORDERS,
    PERMISSIONS.STORE.MANAGE_ORDERS,
    PERMISSIONS.STORE.VIEW_ANALYTICS,
    PERMISSIONS.COMMON.VIEW_PROFILE,
    PERMISSIONS.COMMON.EDIT_PROFILE,
  ],
  
  STORE_READONLY: [
    PERMISSIONS.STORE.VIEW_ORDERS,
    PERMISSIONS.STORE.VIEW_ANALYTICS,
    PERMISSIONS.COMMON.VIEW_PROFILE,
  ],
} as const;

// Permission utility functions
export const permissionUtils = {
  /**
   * Check if user has specific permission
   */
  hasPermission: (userPermissions: string[], permission: string): boolean => {
    return userPermissions.includes(permission);
  },

  /**
   * Check if user has all required permissions
   */
  hasAllPermissions: (userPermissions: string[], requiredPermissions: string[]): boolean => {
    return requiredPermissions.every(permission => 
      userPermissions.includes(permission)
    );
  },

  /**
   * Check if user has any of the required permissions
   */
  hasAnyPermission: (userPermissions: string[], requiredPermissions: string[]): boolean => {
    return requiredPermissions.some(permission => 
      userPermissions.includes(permission)
    );
  },

  /**
   * Get permissions by role
   */
  getPermissionsByRole: (role: 'admin' | 'store'): string[] => {
    switch (role) {
      case 'admin':
        return [...PERMISSION_GROUPS.ADMIN_ALL];
      case 'store':
        return [...PERMISSION_GROUPS.STORE_ALL];
      default:
        return [];
    }
  },

  /**
   * Check if permission is admin-only
   */
  isAdminPermission: (permission: string): boolean => {
    return permission.startsWith('admin:');
  },

  /**
   * Check if permission is store-only
   */
  isStorePermission: (permission: string): boolean => {
    return permission.startsWith('store:');
  },

  /**
   * Get permission description
   */
  getPermissionDescription: (permission: string): string => {
    const descriptions: Record<string, string> = {
      [PERMISSIONS.ADMIN.MANAGE_STORES]: '店舗管理',
      [PERMISSIONS.ADMIN.MANAGE_MANAGERS]: '管理者管理',
      [PERMISSIONS.ADMIN.VIEW_ANALYTICS]: '分析データ閲覧',
      [PERMISSIONS.ADMIN.SYSTEM_SETTINGS]: 'システム設定',
      [PERMISSIONS.STORE.MANAGE_SEATS]: '座席管理',
      [PERMISSIONS.STORE.MANAGE_MENU]: 'メニュー管理',
      [PERMISSIONS.STORE.VIEW_ORDERS]: '注文閲覧',
      [PERMISSIONS.STORE.MANAGE_ORDERS]: '注文管理',
      [PERMISSIONS.STORE.VIEW_ANALYTICS]: '分析データ閲覧',
      [PERMISSIONS.COMMON.VIEW_PROFILE]: 'プロフィール閲覧',
      [PERMISSIONS.COMMON.EDIT_PROFILE]: 'プロフィール編集',
    };

    return descriptions[permission] || permission;
  },
};

export default permissionUtils;