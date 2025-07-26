import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

import { getAuthToken, setAuthToken, clearAuthToken } from '@/lib/api';

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
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  // Actions
  login: (token: string, user: AuthUser) => void;
  logout: () => void;
  setUser: (user: AuthUser) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  
  // Utilities
  hasPermission: (permission: string) => boolean;
  isAdmin: () => boolean;
  isStore: () => boolean;
}

// Create auth store
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Actions
      login: (token: string, user: AuthUser) => {
        setAuthToken(token);
        set({
          token,
          user,
          isAuthenticated: true,
          error: null,
        });
      },

      logout: () => {
        clearAuthToken();
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          error: null,
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
    }),
    {
      name: 'pos-qr-auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Sync token with localStorage on rehydration
        if (state?.token) {
          setAuthToken(state.token);
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