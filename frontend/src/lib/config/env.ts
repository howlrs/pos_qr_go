// Environment configuration with validation
function getEnvVar(key: string, defaultValue?: string): string {
  const value = process.env[key];
  if (!value && !defaultValue) {
    throw new Error(`Environment variable ${key} is required but not set`);
  }
  return value || (defaultValue as string);
}

function getBooleanEnvVar(key: string, defaultValue: boolean = false): boolean {
  const value = process.env[key];
  if (!value) return defaultValue;
  return value.toLowerCase() === 'true';
}

export const env = {
  // API Configuration
  API_URL: getEnvVar('NEXT_PUBLIC_API_URL', 'http://localhost:8080'),
  APP_NAME: getEnvVar('NEXT_PUBLIC_APP_NAME', 'POS QR System'),
  APP_VERSION: getEnvVar('NEXT_PUBLIC_APP_VERSION', '1.0.0'),
  ENVIRONMENT: getEnvVar('NEXT_PUBLIC_ENVIRONMENT', 'development'),
  
  // Debug Configuration
  DEBUG: getBooleanEnvVar('NEXT_PUBLIC_DEBUG', false),
  
  // Server-side only (JWT_SECRET should not be exposed to client)
  JWT_SECRET: process.env.JWT_SECRET || 'default-secret-change-in-production',
} as const;

// Environment helpers
export const isDevelopment = env.ENVIRONMENT === 'development';
export const isProduction = env.ENVIRONMENT === 'production';
export const isTest = env.ENVIRONMENT === 'test';

// API URL helpers
export const getApiUrl = (path: string = ''): string => {
  const baseUrl = env.API_URL.replace(/\/$/, ''); // Remove trailing slash
  const cleanPath = path.replace(/^\//, ''); // Remove leading slash
  return cleanPath ? `${baseUrl}/${cleanPath}` : baseUrl;
};

// Validation
if (isProduction && env.JWT_SECRET === 'default-secret-change-in-production') {
  // eslint-disable-next-line no-console
  console.warn('⚠️  Using default JWT_SECRET in production. Please set a secure JWT_SECRET.');
}

// Export types for TypeScript
export type Environment = typeof env.ENVIRONMENT;
export type EnvConfig = typeof env;
