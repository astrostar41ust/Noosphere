'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useCurrentUser } from '@/features/auth';
import { Loader2 } from 'lucide-react';

export default function Home() {
  const { data: currentUser, isLoading } = useCurrentUser();
  const router = useRouter();

  useEffect(() => {
    if (isLoading) return;
    if (currentUser) {
      router.replace('/chat');
    } else {
      router.replace('/login');
    }
  }, [currentUser, isLoading, router]);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-slate-950 text-slate-100 font-sans">
      <Loader2 className="animate-spin text-indigo-500" size={28} />
      <p className="text-xs text-slate-400 mt-3 tracking-wide">Syncing Noosphere user profile state...</p>
    </div>
  );
}
