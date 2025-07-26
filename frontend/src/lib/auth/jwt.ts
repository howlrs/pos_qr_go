import { jwtDecode } from 'jwt-decode';

import { AuthUser } from '@/store/auth/authStore';

// JWT payload interface
export interface JWTPayload {
  sub: string; // user ID
  email: string;
  name: string;
  role: 'admin' | 'store';
  storeId?: string;
  permissions: string[];
  iat: number; // issued at
  exp: number; // expires at
}

// JWT utilities
export const jwt = {
  /**
   * Decode JWT token and extract payload
   */
  decode: (token: string): JWTPayload | null => {
    try {
      return jwtDecode<JWTPayload>(token);
    } catch (error) {
      console.error('Failed to decode JWT token:', error);
      return null;
    }
  },

  /**
   * Check if JWT token is expired
   */
  isExpired: (token: string): boolean => {
    try {
      const payload = jwtDecode<JWTPayload>(token);
      const currentTime = Date.now() / 1000;
      return payload.exp < currentTime;
    } catch (error) {
      return true;
    }
  },

  /**
   * Check if JWT token is valid (not expired and properly formatted)
   */
  isValid: (token: string): boolean => {
    try {
      const payload = jwtDecode<JWTPayload>(token);
      const currentTime = Date.now() / 1000;
      
      // Check if token has required fields and is not expired
      return !!(
        payload.sub &&
        payload.email &&
        payload.role &&
        payload.exp > currentTime
      );
    } catch (error) {
      return false;
    }
  },

  /**
   * Extract user information from JWT token
   */
  extractUser: (token: string): AuthUser | null => {
    try {
      const payload = jwtDecode<JWTPayload>(token);
      
      return {
        id: payload.sub,
        email: payload.email,
        name: payload.name,
        role: payload.role,
        storeId: payload.storeId,
        permissions: payload.permissions,
      };
    } catch (error) {
      console.error('Failed to extract user from JWT token:', error);
      return null;
    }
  },

  /**
   * Get token expiration time in milliseconds
   */
  getExpirationTime: (token: string): number | null => {
    try {
      const payload = jwtDecode<JWTPayload>(token);
      return payload.exp * 1000; // Convert to milliseconds
    } catch (error) {
      return null;
    }
  },

  /**
   * Get time until token expires in milliseconds
   */
  getTimeUntilExpiration: (token: string): number | null => {
    try {
      const expirationTime = jwt.getExpirationTime(token);
      if (!expirationTime) return null;
      
      const currentTime = Date.now();
      const timeUntilExpiration = expirationTime - currentTime;
      
      return Math.max(0, timeUntilExpiration);
    } catch (error) {
      return null;
    }
  },

  /**
   * Check if token will expire within specified minutes
   */
  willExpireSoon: (token: string, minutes: number = 5): boolean => {
    try {
      const timeUntilExpiration = jwt.getTimeUntilExpiration(token);
      if (!timeUntilExpiration) return true;
      
      const thresholdMs = minutes * 60 * 1000;
      return timeUntilExpiration <= thresholdMs;
    } catch (error) {
      return true;
    }
  },
};

export default jwt;