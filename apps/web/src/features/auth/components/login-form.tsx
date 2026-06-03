'use client';

import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Link from 'next/link';
import { Loader2, Lock, Mail, AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useLogin } from '../hooks/use-auth';
import { loginFormSchema, LoginFormData } from '../schemas';

export function LoginForm() {
  const { mutate: login, isPending, error } = useLogin();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: { email: '', password: '' },
  });

  const onSubmit = (data: LoginFormData) => {
    login(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      {error && (
        <div className="flex gap-3 p-4 rounded-xl bg-red-950/40 border border-red-900/50 text-red-400 text-xs">
          <AlertTriangle size={16} className="shrink-0" />
          <p>{error.message}</p>
        </div>
      )}

      <div className="space-y-1.5">
        <label className="text-xs font-semibold text-slate-400 tracking-wider flex items-center gap-1.5">
          <Mail size={12} /> Email Node
        </label>
        <Input
          {...register('email')}
          type="email"
          placeholder="developer@noosphere.network"
          disabled={isPending}
          className="bg-slate-950/60 border-slate-800 focus-visible:ring-indigo-500 rounded-xl px-4 py-5 text-sm"
        />
        {errors.email && (
          <p className="text-[11px] text-red-400 font-medium pl-1">{errors.email.message}</p>
        )}
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold text-slate-400 tracking-wider flex items-center gap-1.5">
          <Lock size={12} /> Authorization Pass
        </label>
        <Input
          {...register('password')}
          type="password"
          placeholder="••••••••"
          disabled={isPending}
          className="bg-slate-950/60 border-slate-800 focus-visible:ring-indigo-500 rounded-xl px-4 py-5 text-sm"
        />
        {errors.password && (
          <p className="text-[11px] text-red-400 font-medium pl-1">{errors.password.message}</p>
        )}
      </div>

      <Button
        type="submit"
        disabled={isPending}
        className="w-full bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl py-5 font-semibold text-sm transition-all shadow-md hover:shadow-indigo-500/10 active:scale-98 mt-2 h-auto flex items-center justify-center gap-2"
      >
        {isPending ? (
          <>
            <Loader2 className="animate-spin" size={16} />
            Authenticating Core Handshake...
          </>
        ) : (
          'Establish Connection'
        )}
      </Button>

      <p className="text-center text-xs text-slate-400 mt-4">
        Need permission credentials?{' '}
        <Link href="/register" className="text-indigo-400 hover:text-indigo-300 font-medium transition-colors">
          Create Core Account
        </Link>
      </p>
    </form>
  );
}
