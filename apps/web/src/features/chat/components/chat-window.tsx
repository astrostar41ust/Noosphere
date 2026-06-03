import React, { useRef, useEffect } from 'react';
import { Bot, Loader2, RefreshCw, LogOut } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { MessageBubble } from './message-bubble';
import { ChatInput } from './chat-input';
import { useChatMessages } from '../hooks/use-chat-messages';
import { useSendMessage } from '../hooks/use-send-message';
import { useCurrentUser, useLogout } from '@/features/auth';

interface ChatWindowProps {
  sessionId: string;
}

export function ChatWindow({ sessionId }: ChatWindowProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const { data: messages = [], isLoading: isHistoryLoading, isError: isHistoryError, refetch } = useChatMessages(sessionId);
  const { mutate: sendMessage, isPending: isSending } = useSendMessage(sessionId);
  
  const { data: currentUser } = useCurrentUser();
  const logout = useLogout();

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isSending]);

  const handleSendMessage = (content: string) => {
    sendMessage(content);
  };

  return (
    <div className="flex flex-col h-screen bg-slate-950 text-slate-100 antialiased font-sans">
      {/* Top Header Section */}
      <header className="flex items-center justify-between px-6 py-4 border-b border-slate-800 bg-slate-900/40 backdrop-blur">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-indigo-600/10 text-indigo-400 border border-indigo-500/20">
            <Bot size={20} />
          </div>
          <div>
            <h1 className="text-sm font-semibold tracking-wide">Noosphere Core Interface</h1>
            <p className="text-[11px] text-slate-400">
              {currentUser ? `Active Node: ${currentUser.username}` : 'Radix UI Context Nodes Active'}
            </p>
          </div>
        </div>
        {currentUser && (
          <Button
            variant="outline"
            size="sm"
            onClick={logout}
            className="gap-2 border-slate-800 hover:bg-red-950/20 hover:text-red-400 hover:border-red-900/30 text-xs rounded-xl px-3 transition-all"
          >
            <LogOut size={12} /> Terminate Link
          </Button>
        )}
      </header>


      {/* Message Output Viewport Canvas */}
      <main className="flex-1 overflow-y-auto px-4 py-6 md:px-8 space-y-6 max-w-4xl mx-auto w-full">
        {isHistoryLoading && (
          <div className="flex flex-col items-center justify-center h-full space-y-3 py-20">
            <Loader2 className="animate-spin text-indigo-500" size={24} />
            <p className="text-xs text-slate-400">Synchronizing database log streams...</p>
          </div>
        )}

        {isHistoryError && (
          <div className="flex flex-col items-center justify-center h-full space-y-4 py-20 max-w-sm mx-auto text-center">
            <p className="text-xs text-red-400">Unable to link with Go cluster API.</p>
            <Button variant="outline" size="sm" onClick={() => refetch()} className="gap-2 border-slate-800">
              <RefreshCw size={12} /> Retry Handshake
            </Button>
          </div>
        )}

        {!isHistoryLoading && !isHistoryError && messages.length === 0 && (
          <div className="flex flex-col items-center justify-center h-full text-slate-500 py-20 text-center space-y-2">
            <Bot size={36} className="stroke-1 animate-pulse text-indigo-500/40" />
            <p className="text-xs">Noosphere system initialized. Post an instruction node to begin model calculations.</p>
          </div>
        )}

        {!isHistoryLoading && !isHistoryError && messages.map((msg) => (
          <MessageBubble key={msg.id} message={msg} />
        ))}

        {/* Neural Network Think Loader */}
        {isSending && (
          <div className="flex items-center gap-3 mr-auto bg-slate-900/40 border border-slate-800/60 p-4 rounded-xl rounded-tl-none max-w-[80%] text-slate-400">
            <Loader2 size={14} className="animate-spin text-indigo-400" />
            <p className="text-xs tracking-wide animate-pulse">Calculating local model weights...</p>
          </div>
        )}
        <div ref={messagesEndRef} />
      </main>

      {/* Message Input Control Tray */}
      <footer className="border-t border-slate-800 bg-slate-900/20 p-4 md:p-6 backdrop-blur">
        <ChatInput onSendMessage={handleSendMessage} disabled={isSending || isHistoryLoading} />
      </footer>
    </div>
  );
}
