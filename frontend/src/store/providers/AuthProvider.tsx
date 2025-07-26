'use client';

import { useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';

import { useAuthStore } from '@/store/auth/authStore';
import { api, API_ENDPOINTS } from '@/lib/api';
import { jwt } from '@/lib/auth/jwt';

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const authStore = useAuthStore();
  const queryClient = useQueryClient();
  const refreshIntervalRef = useRef<NodeJS.Timeout | null>(null);

  // Auto token refresh function
  const refreshTokenIfNeeded = async () => {
    const { token, refreshToken, isAuthenticated, willTokenExpireSoon, updateTokens, logout } = authStore;

    if (!isAuthenticated || !token || !refreshToken) {
      return;
    }

    // Check if token will expire soon (within 5 minutes)
    if (willTokenExpireSoon(5)) {
      try {
        console.log('Token will expire soon, refreshing...');
        
        const response = await api.post(API_ENDPOINTS.AUTH.REFRESH, {
          refreshToken,
        });

        const responseData = response.data.data as { token: string; refreshToken?: string };
        const { token: newToken, refreshToken: newRefreshToken } = responseData;

        if (jwt.isValid(newToken)) {
          updateTokens(newToken, newRefreshToken);
          console.log('Token refreshed successfully');
        } else {
          throw new Error('Invalid refresh token received');
        }
      } catch (error) {
        console.error('Token refresh failed:', error);
        logout();
        queryClient.clear();
      }
    }
  };

  // Session validation function
  const validateSession = () => {
    const { validateSession } = authStore;
    
    if (!validateSession()) {
      queryClient.clear();
    }
  };

  // Setup auto refresh interval
  useEffect(() => {
    const { isAuthenticated } = authStore;

    if (isAuthenticated) {
      // Validate session immediately
      validateSession();

      // Set up interval to check token expiration every minute
      refreshIntervalRef.current = setInterval(() => {
        refreshTokenIfNeeded();
      }, 60 * 1000); // Check every minute

      return () => {
        if (refreshIntervalRef.current) {
          clearInterval(refreshIntervalRef.current);
        }
      };
    } else {
      // Clear interval if not authenticated
      if (refreshIntervalRef.current) {
        clearInterval(refreshIntervalRef.current);
      }
    }
  }, [authStore.isAuthenticated]);

  // Validate session on window focus
  useEffect(() => {
    const handleFocus = () => {
      validateSession();
      refreshTokenIfNeeded();
    };

    window.addEventListener('focus', handleFocus);
    return () => window.removeEventListener('focus', handleFocus);
  }, []);

  // Validate session on page visibility change
  useEffect(() => {
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        validateSession();
        refreshTokenIfNeeded();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
  }, []);

  return <>{children}</>;
};

export default AuthProvider;