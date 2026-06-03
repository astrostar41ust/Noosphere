import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Send } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const chatFormSchema = z.object({
  message: z.string().min(1, 'Message cannot be empty').max(2000, 'Message exceeds length limit'),
});

type ChatFormData = z.infer<typeof chatFormSchema>;

interface ChatInputProps {
  onSendMessage: (content: string) => void;
  disabled: boolean;
}

export function ChatInput({ onSendMessage, disabled }: ChatInputProps) {
  const { register, handleSubmit, reset, watch } = useForm<ChatFormData>({
    resolver: zodResolver(chatFormSchema),
    defaultValues: { message: '' },
  });

  const messageInputValue = watch('message');

  const onSubmit = (data: ChatFormData) => {
    const text = data.message.trim();
    if (!text) return;
    reset();
    onSendMessage(text);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="max-w-4xl mx-auto flex gap-3">
      <Input
        {...register('message')}
        type="text"
        placeholder="Query your local inference architecture..."
        disabled={disabled}
        className="flex-1 bg-slate-900 border-slate-800 focus-visible:ring-indigo-500 rounded-xl px-4 py-6 text-sm"
      />
      <Button
        type="submit"
        disabled={!messageInputValue?.trim() || disabled}
        className="bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl h-auto px-5 transition-all shadow-sm active:scale-95"
      >
        <Send size={15} />
      </Button>
    </form>
  );
}
