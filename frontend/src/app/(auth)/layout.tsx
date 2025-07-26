'use client';

import { AuthGuard } from '@/lib/auth/guards';

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayoutPage({ children }: AuthLayoutProps) {
  return (
    <AuthGuard>
      {children}
    </AuthGuard>
  );
}