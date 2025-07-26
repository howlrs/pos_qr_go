import { useMemo } from 'react';

import { useAuth } from '@/store/auth/authStore';
import { permissionUtils } from '@/lib/auth/permissions';

export const usePermissions = () => {
  const { user } = useAuth();

  const permissions = useMemo(() => {
    if (!user) return [];
    return user.permissions;
  }, [user]);

  const hasPermission = useMemo(() => {
    return (permission: string): boolean => {
      return permissionUtils.hasPermission(permissions, permission);
    };
  }, [permissions]);

  const hasAllPermissions = useMemo(() => {
    return (requiredPermissions: string[]): boolean => {
      return permissionUtils.hasAllPermissions(permissions, requiredPermissions);
    };
  }, [permissions]);

  const hasAnyPermission = useMemo(() => {
    return (requiredPermissions: string[]): boolean => {
      return permissionUtils.hasAnyPermission(permissions, requiredPermissions);
    };
  }, [permissions]);

  const isAdmin = useMemo(() => {
    return user?.role === 'admin';
  }, [user]);

  const isStore = useMemo(() => {
    return user?.role === 'store';
  }, [user]);

  return {
    permissions,
    hasPermission,
    hasAllPermissions,
    hasAnyPermission,
    isAdmin,
    isStore,
    user,
  };
};

export default usePermissions;