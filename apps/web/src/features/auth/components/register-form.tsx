'use client';

import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import Link from 'next/link';
import { Loader2, Lock, Mail, User, AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useRegister } from '../hooks/use-auth';
import { registerFormSchema, RegisterFormData } from '../schemas';

export function RegisterForm() {
  const { mutate: registerUser, isPending, error } = useRegister();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerFormSchema),
    defaultValues: { username: '', email: '', password: '', confirmPassword: '' },
  });

  const onSubmit = (data: RegisterFormData) => {
    registerUser(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {error && (
        <div className="flex gap-3 p-4 rounded-xl bg-red-950/40 border border-red-900/50 text-red-400 text-xs">
          <AlertTriangle size={16} className="shrink-0" />
          <p>{error.message}</p>
        </div>
      )}

      <div className="space-y-1.5">
        <label className="text-xs font-semibold text-slate-400 tracking-wider flex items-center gap-1.5">
          <User size={12} /> Username Node
        </label>
        <Input
          {...register('username')}
          type="text"
          placeholder="noosphere_node"
          disabled={isPending}
          className="bg-slate-950/60 border-slate-800 focus-visible:ring-indigo-500 rounded-xl px-4 py-5 text-sm"
        />
        {errors.username && (
          <p className="text-[11px] text-red-400 font-medium pl-1">{errors.username.message}</p>
        )}
      </div>

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
          <Lock size={12} /> Passphrase
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

      <div className="space-y-1.5">
        <label className="text-xs font-semibold text-slate-400 tracking-wider flex items-center gap-1.5">
          <Lock size={12} /> Verify Passphrase
        </label>
        <Input
          {...register('confirmPassword')}
          type="password"
          placeholder="••••••••"
          disabled={isPending}
          className="bg-slate-950/60 border-slate-800 focus-visible:ring-indigo-500 rounded-xl px-4 py-5 text-sm"
        />
        {errors.confirmPassword && (
          <p className="text-[11px] text-red-400 font-medium pl-1">{errors.confirmPassword.message}</p>
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
            Provisioning Client Node...
          </>
        ) : (
          'Register Authorization Node'
        )}
      </Button>

      <p className="text-center text-xs text-slate-400 mt-4">
        Already have credentials?{' '}
        <Link href="/login" className="text-indigo-400 hover:text-indigo-300 font-medium transition-colors">
          Access Console
        </Link>
      </p>
    </form>
  );
}
