// API client exports
export { default as apiClient, api } from './client';
export type { ApiResponse, ApiError } from './client';
export { getAuthToken, setAuthToken, clearAuthToken } from './client';

// API endpoints
export { API_ENDPOINTS, buildUrl } from './endpoints';
export type { ApiEndpoints } from './endpoints';

// Re-export axios types for convenience
export type {
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
} from 'axios';