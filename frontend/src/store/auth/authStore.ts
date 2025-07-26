import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

import { getAuthToken, setAuthToken, clearAuthToken } from '@/lib/api';
import { jwt } from '@/lib/auth/jwt';

// Auth user types
export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'store';
  storeId?: string;
  permissions: string[];
}

// Auth state interface
export interface AuthState {
  // State
  user: AuthUser | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  tokenExpiresAt: number | null;

  // Actions
  login: (token: string, user: AuthUser, refreshToken?: string) => void;
  logout: () => void;
  setUser: (user: AuthUser) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  updateTokens: (token: string, refreshToken?: string) => void;
  
  // Utilities
  hasPermission: (permission: string) => boolean;
  isAdmin: () => boolean;
  isStore: () => boolean;
  isTokenExpired: () => boolean;
  willTokenExpireSoon: (minutes?: number) => boolean;
  validateSession: () => boolean;
}

// Create auth store
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      tokenExpiresAt: null,

      // Actions
      login: (token: string, user: AuthUser, refreshToken?: string) => {
        setAuthToken(token);
        const expirationTime = jwt.getExpirationTime(token);
        
        set({
          token,
          user,
          refreshToken,
          isAuthenticated: true,
          error: null,
          tokenExpiresAt: expirationTime,
        });
      },

      logout: () => {
        clearAuthToken();
        set({
          user: null,
          token: null,
          refreshToken: null,
          isAuthenticated: false,
          error: null,
          tokenExpiresAt: null,
        });
      },

      setUser: (user: AuthUser) => {
        set({ user });
      },

      setLoading: (isLoading: boolean) => {
        set({ isLoading });
      },

      setError: (error: string | null) => {
        set({ error });
      },

      clearError: () => {
        set({ error: null });
      },

      updateTokens: (token: string, refreshToken?: string) => {
        setAuthToken(token);
        const expirationTime = jwt.getExpirationTime(token);
        
        set({
          token,
          refreshToken: refreshToken || get().refreshToken,
          tokenExpiresAt: expirationTime,
          error: null,
        });
      },

      // Utilities
      hasPermission: (permission: string): boolean => {
        const { user } = get();
        return user?.permissions.includes(permission) ?? false;
      },

      isAdmin: (): boolean => {
        const { user } = get();
        return user?.role === 'admin';
      },

      isStore: (): boolean => {
        const { user } = get();
        return user?.role === 'store';
      },

      isTokenExpired: (): boolean => {
        const { token } = get();
        if (!token) return true;
        return jwt.isExpired(token);
      },

      willTokenExpireSoon: (minutes: number = 5): boolean => {
        const { token } = get();
        if (!token) return true;
        return jwt.willExpireSoon(token, minutes);
      },

      validateSession: (): boolean => {
        const { token, user, isAuthenticated } = get();
        
        if (!isAuthenticated || !token || !user) {
          return false;
        }

        if (jwt.isExpired(token)) {
          // Token is expired, logout
          get().logout();
          return false;
        }

        return true;
      },
    }),
    {
      name: 'pos-qr-auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
        tokenExpiresAt: state.tokenExpiresAt,
      }),
      onRehydrateStorage: () => (state) => {
        // Validate session on rehydration
        if (state?.token) {
          if (jwt.isValid(state.token)) {
            setAuthToken(state.token);
          } else {
            // Token is invalid, clear auth state
            state.logout();
          }
        }
      },
    }
  )
);

// Auth store selectors
export const useAuth = () => useAuthStore();
export const useAuthUser = () => useAuthStore((state) => state.user);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useAuthLoading = () => useAuthStore((state) => state.isLoading);
export const useAuthError = () => useAuthStore((state) => state.error);