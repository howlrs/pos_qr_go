'use client';

import { StoreGuard } from '@/lib/auth/guards';
import { StoreLayout } from '@/components/layouts';

interface StoreLayoutPageProps {
  children: React.ReactNode;
}

export default function StoreLayoutPage({ children }: StoreLayoutPageProps) {
  return (
    <StoreGuard>
      <StoreLayout>
        {children}
      </StoreLayout>
    </StoreGuard>
  );
}