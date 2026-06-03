import React from 'react';
import { Bot, User } from 'lucide-react';
import { Message } from '../types';

interface MessageBubbleProps {
  message: Message;
}

export function MessageBubble({ message }: MessageBubbleProps) {
  const isUser = message.role === 'user';

  return (
    <div
      className={`flex items-start gap-4 p-4 rounded-xl max-w-[85%] transition-all ${
        isUser
          ? 'ml-auto bg-indigo-600 text-white rounded-tr-none shadow-md shadow-indigo-600/10'
          : 'mr-auto bg-slate-900 border border-slate-800/80 rounded-tl-none'
      }`}
    >
      <div className={`p-1.5 rounded-md ${isUser ? 'bg-indigo-700' : 'bg-slate-800 text-slate-400'}`}>
        {isUser ? <User size={14} /> : <Bot size={14} />}
      </div>
      <div className="flex-1 space-y-1 overflow-x-auto">
        <p className="text-sm leading-relaxed whitespace-pre-wrap">{message.content}</p>
        <span className="block text-[9px] opacity-40 text-right">
          {new Date(message.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </span>
      </div>
    </div>
  );
}
