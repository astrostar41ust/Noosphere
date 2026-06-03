import React from 'react';
import { AuthLayout, LoginForm } from '@/features/auth';

export default function LoginPage() {
  return (
    <AuthLayout
      title="Noosphere Console Authorization"
      subtitle="Input credential nodes to establish link"
    >
      <LoginForm />
    </AuthLayout>
  );
}
