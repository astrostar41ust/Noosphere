'use client';

import React from 'react';
import { ChatWindow } from '@/features/chat';

const TEST_SESSION_ID = 'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d';

export default function ChatPage() {
  return <ChatWindow sessionId={TEST_SESSION_ID} />;
}