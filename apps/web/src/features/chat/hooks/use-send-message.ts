import { useMutation, useQueryClient } from '@tanstack/react-query';
import { sendChatMessage } from '../services/chat-api';
import { Message } from '../types';

export function useSendMessage(sessionId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (content: string) => sendChatMessage(sessionId, content),
    onMutate: async (newUserText) => {
      await queryClient.cancelQueries({ queryKey: ['chatHistory', sessionId] });
      const previousHistory = queryClient.getQueryData<Message[]>(['chatHistory', sessionId]);

      const optimisticUserMsg: Message = {
        id: crypto.randomUUID(),
        session_id: sessionId,
        role: 'user',
        content: newUserText,
        created_at: new Date().toISOString(),
      };

      queryClient.setQueryData<Message[]>(
        ['chatHistory', sessionId],
        (old) => [...(old || []), optimisticUserMsg]
      );

      return { previousHistory };
    },
    onSuccess: (data) => {
      queryClient.setQueryData<Message[]>(['chatHistory', sessionId], (old) => {
        const current = old ? [...old] : [];
        return [...current, data];
      });
    },
    onError: (err, newUserText, context) => {
      if (context?.previousHistory) {
        queryClient.setQueryData(['chatHistory', sessionId], context.previousHistory);
      }
    },
  });
}
