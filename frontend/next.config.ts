import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // Ant Design optimization
  transpilePackages: ['antd'],

  // API proxy for development
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.NEXT_PUBLIC_API_URL}/:path*`,
      },
    ];
  },

  // Environment-specific configuration
  env: {
    CUSTOM_KEY: process.env.CUSTOM_KEY,
  },

  // Image optimization
  images: {
    domains: ['localhost'],
    unoptimized: process.env.NODE_ENV === 'development',
  },

  // Experimental features
  experimental: {
    optimizePackageImports: ['antd'],
  },
};

export default nextConfig;
