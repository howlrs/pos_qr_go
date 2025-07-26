'use client';

import { AuthLayout } from '@/components/layouts';

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayoutPage({ children }: AuthLayoutProps) {
  return <AuthLayout>{children}</AuthLayout>;
}