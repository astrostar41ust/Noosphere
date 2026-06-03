'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ChatWindow } from '@/features/chat';
import { useCurrentUser } from '@/features/auth';
import { Loader2 } from 'lucide-react';

const TEST_SESSION_ID = 'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d';

export default function ChatPage() {
  const { data: currentUser, isLoading } = useCurrentUser();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !currentUser) {
      router.replace('/login');
    }
  }, [currentUser, isLoading, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-slate-950 text-slate-100 font-sans">
        <Loader2 className="animate-spin text-indigo-500" size={28} />
        <p className="text-xs text-slate-400 mt-3 tracking-wide">Verifying secure node session...</p>
      </div>
    );
  }

  if (!currentUser) {
    return null; // Will redirect in useEffect
  }

  return <ChatWindow sessionId={TEST_SESSION_ID} />;
}