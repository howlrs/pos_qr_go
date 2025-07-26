'use client';

import dynamic from 'next/dynamic';
import { Spin } from 'antd';

// Loading component
const LoadingSpinner = () => (
  <div className="flex justify-center items-center p-8">
    <Spin size="large" />
  </div>
);

// Lazy load heavy components
export const LazyOrderConfirmation = dynamic(
  () => import('./OrderConfirmation').then(mod => ({ default: mod.OrderConfirmation })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

export const LazyOrderSuccess = dynamic(
  () => import('./OrderSuccess').then(mod => ({ default: mod.OrderSuccess })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

export const LazyOrderHistory = dynamic(
  () => import('./OrderHistory').then(mod => ({ default: mod.OrderHistory })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

export const LazyStoreInfo = dynamic(
  () => import('./StoreInfo').then(mod => ({ default: mod.StoreInfo })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

export const LazyCartDrawer = dynamic(
  () => import('./CartDrawer').then(mod => ({ default: mod.CartDrawer })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

// Admin components (loaded only when needed)
export const LazyAdminDashboard = dynamic(
  () => import('../../app/(auth)/admin/dashboard/page'),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

export const LazyStoreDashboard = dynamic(
  () => import('../../app/(auth)/store/dashboard/page'),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);

// Chart components (heavy libraries)
export const LazyChart = dynamic(
  () => import('antd').then(mod => ({ default: mod.Progress })),
  {
    loading: LoadingSpinner,
    ssr: false,
  }
);