import React from 'react';
import { AuthLayout, RegisterForm } from '@/features/auth';

export default function RegisterPage() {
  return (
    <AuthLayout
      title="Create Noosphere Node Profile"
      subtitle="Register authorization keys on the cluster network"
    >
      <RegisterForm />
    </AuthLayout>
  );
}
