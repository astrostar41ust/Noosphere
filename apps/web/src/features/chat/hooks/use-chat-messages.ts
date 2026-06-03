import { useQuery } from '@tanstack/react-query';
import { fetchChatHistory } from '../services/chat-api';
import { Message } from '../types';

export function useChatMessages(sessionId: string) {
  return useQuery<Message[]>({
    queryKey: ['chatHistory', sessionId],
    queryFn: () => fetchChatHistory(sessionId),
  });
}
