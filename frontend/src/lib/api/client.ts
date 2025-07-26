import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

import { env, isDevelopment } from '../config/env';

// API Response types
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  success: boolean;
}

export interface ApiError {
  message: string;
  code?: string;
  details?: Record<string, unknown>;
}

// Create axios instance
const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: env.API_URL,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor
  client.interceptors.request.use(
    (config) => {
      // Add auth token if available
      const token = getAuthToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      // Add request ID for debugging
      if (isDevelopment) {
        config.headers['X-Request-ID'] = generateRequestId();
      }

      // Log request in development
      if (isDevelopment && env.DEBUG) {
        // eslint-disable-next-line no-console
        console.log('🚀 API Request:', {
          method: config.method?.toUpperCase(),
          url: config.url,
          data: config.data,
        });
      }

      return config;
    },
    (error) => {
      if (isDevelopment) {
        // eslint-disable-next-line no-console
        console.error('❌ Request Error:', error);
      }
      return Promise.reject(error);
    }
  );

  // Response interceptor
  client.interceptors.response.use(
    (response: AxiosResponse) => {
      // Log response in development
      if (isDevelopment && env.DEBUG) {
        // eslint-disable-next-line no-console
        console.log('✅ API Response:', {
          status: response.status,
          url: response.config.url,
          data: response.data,
        });
      }

      return response;
    },
    (error) => {
      // Handle different error types
      if (error.response) {
        // Server responded with error status
        const { status, data } = error.response;
        
        if (status === 401) {
          // Unauthorized - clear auth token and redirect to login
          clearAuthToken();
          if (typeof window !== 'undefined') {
            window.location.href = '/auth/admin-login';
          }
        }

        if (isDevelopment) {
          // eslint-disable-next-line no-console
          console.error('❌ API Error Response:', {
            status,
            url: error.config?.url,
            data,
          });
        }

        // Transform error to consistent format
        const apiError: ApiError = {
          message: data?.message || `Request failed with status ${status}`,
          code: data?.code || `HTTP_${status}`,
          details: data?.details || {},
        };

        return Promise.reject(apiError);
      } else if (error.request) {
        // Network error
        const networkError: ApiError = {
          message: 'Network error - please check your connection',
          code: 'NETWORK_ERROR',
        };

        if (isDevelopment) {
          // eslint-disable-next-line no-console
          console.error('❌ Network Error:', error.request);
        }

        return Promise.reject(networkError);
      } else {
        // Other error
        const unknownError: ApiError = {
          message: error.message || 'An unexpected error occurred',
          code: 'UNKNOWN_ERROR',
        };

        if (isDevelopment) {
          // eslint-disable-next-line no-console
          console.error('❌ Unknown Error:', error);
        }

        return Promise.reject(unknownError);
      }
    }
  );

  return client;
};

// Auth token management
const AUTH_TOKEN_KEY = 'pos_qr_auth_token';

export const getAuthToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(AUTH_TOKEN_KEY);
};

export const setAuthToken = (token: string): void => {
  if (typeof window === 'undefined') return;
  localStorage.setItem(AUTH_TOKEN_KEY, token);
};

export const clearAuthToken = (): void => {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(AUTH_TOKEN_KEY);
};

// Request ID generation for debugging
const generateRequestId = (): string => {
  return Math.random().toString(36).substring(2, 15);
};

// Create and export the API client instance
export const apiClient = createApiClient();

// Convenience methods
export const api = {
  get: <T = unknown>(url: string, config?: AxiosRequestConfig) =>
    apiClient.get<ApiResponse<T>>(url, config),
  
  post: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.post<ApiResponse<T>>(url, data, config),
  
  put: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.put<ApiResponse<T>>(url, data, config),
  
  patch: <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.patch<ApiResponse<T>>(url, data, config),
  
  delete: <T = unknown>(url: string, config?: AxiosRequestConfig) =>
    apiClient.delete<ApiResponse<T>>(url, config),
};

export default apiClient;