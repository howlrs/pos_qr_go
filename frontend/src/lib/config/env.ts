// Environment configuration
export const env = {
  API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  APP_NAME: process.env.NEXT_PUBLIC_APP_NAME || 'POS QR System',
  ENVIRONMENT: process.env.NEXT_PUBLIC_ENVIRONMENT || 'development',
  JWT_SECRET: process.env.JWT_SECRET || 'default-secret',
} as const;

export const isDevelopment = env.ENVIRONMENT === 'development';
export const isProduction = env.ENVIRONMENT === 'production';
