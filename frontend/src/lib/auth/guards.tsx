'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Spin } from 'antd';

import { useAuth } from '@/store/auth/authStore';
import { jwt } from './jwt';
import { getAuthToken } from '@/lib/api';

// Auth guard props
interface AuthGuardProps {
  children: React.ReactNode;
  requiredRole?: 'admin' | 'store';
  requiredPermissions?: string[];
  fallback?: React.ReactNode;
  redirectTo?: string;
}

// Loading component
const AuthLoading: React.FC = () => (
  <div className="min-h-screen flex items-center justify-center">
    <Spin size="large" />
  </div>
);

// Auth guard component
export const AuthGuard: React.FC<AuthGuardProps> = ({
  children,
  requiredRole,
  requiredPermissions = [],
  fallback,
  redirectTo,
}) => {
  const router = useRouter();
  const {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    logout,
    setLoading,
  } = useAuth();

  useEffect(() => {
    const initializeAuth = async () => {
      setLoading(true);

      try {
        // Get token from storage
        const storedToken = getAuthToken();
        
        if (!storedToken) {
          // No token found, redirect to login
          handleUnauthenticated();
          return;
        }

        // Check if token is valid
        if (!jwt.isValid(storedToken)) {
          // Invalid or expired token, logout and redirect
          logout();
          handleUnauthenticated();
          return;
        }

        // If we have a valid token but no user in store, extract user from token
        if (!user && storedToken) {
          const extractedUser = jwt.extractUser(storedToken);
          if (extractedUser) {
            login(storedToken, extractedUser);
          } else {
            // Failed to extract user, logout and redirect
            logout();
            handleUnauthenticated();
            return;
          }
        }
      } catch (error) {
        console.error('Auth initialization error:', error);
        logout();
        handleUnauthenticated();
      } finally {
        setLoading(false);
      }
    };

    const handleUnauthenticated = () => {
      if (redirectTo) {
        router.push(redirectTo);
      } else {
        // Default redirect based on required role
        const defaultRedirect = requiredRole === 'admin' 
          ? '/auth/admin-login' 
          : '/auth/store-login';
        router.push(defaultRedirect);
      }
    };

    // Only initialize if not already authenticated
    if (!isAuthenticated && !isLoading) {
      initializeAuth();
    }
  }, [
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    logout,
    setLoading,
    router,
    requiredRole,
    redirectTo,
  ]);

  // Show loading while checking authentication
  if (isLoading) {
    return fallback || <AuthLoading />;
  }

  // Check if user is authenticated
  if (!isAuthenticated || !user) {
    return fallback || <AuthLoading />;
  }

  // Check role requirement
  if (requiredRole && user.role !== requiredRole) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">
            Access Denied
          </h1>
          <p className="text-gray-600">
            You don't have permission to access this page.
          </p>
        </div>
      </div>
    );
  }

  // Check permission requirements
  if (requiredPermissions.length > 0) {
    const hasAllPermissions = requiredPermissions.every(permission =>
      user.permissions.includes(permission)
    );

    if (!hasAllPermissions) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-red-600 mb-4">
              Insufficient Permissions
            </h1>
            <p className="text-gray-600">
              You don't have the required permissions to access this page.
            </p>
          </div>
        </div>
      );
    }
  }

  // All checks passed, render children
  return <>{children}</>;
};

// Role-specific guard components
export const AdminGuard: React.FC<Omit<AuthGuardProps, 'requiredRole'>> = (props) => (
  <AuthGuard {...props} requiredRole="admin" />
);

export const StoreGuard: React.FC<Omit<AuthGuardProps, 'requiredRole'>> = (props) => (
  <AuthGuard {...props} requiredRole="store" />
);

// Permission guard component
interface PermissionGuardProps {
  children: React.ReactNode;
  permissions: string[];
  fallback?: React.ReactNode;
}

export const PermissionGuard: React.FC<PermissionGuardProps> = ({
  children,
  permissions,
  fallback,
}) => {
  const { user } = useAuth();

  if (!user) {
    return fallback || null;
  }

  const hasAllPermissions = permissions.every(permission =>
    user.permissions.includes(permission)
  );

  if (!hasAllPermissions) {
    return fallback || null;
  }

  return <>{children}</>;
};

export default AuthGuard;