import api from '@/lib/api-client';
import { Message } from '../types';

export async function fetchChatHistory(sessionId: string): Promise<Message[]> {
  try {
    const response = await api.get<Message[]>(`/api/v1/chat/session/${sessionId}/history`);
    return response.data;
  } catch (error) {
    throw new Error('Failed to synchronize chat timeline logs.');
  }
}

export async function sendChatMessage(sessionId: string, content: string): Promise<Message> {
  try {
    const response = await api.post<Message>('/api/v1/chat/message', {
      session_id: sessionId,
      role: 'user',
      content,
    });
    return response.data;
  } catch (error) {
    throw new Error('Inference engine processing dropped.');
  }
}
