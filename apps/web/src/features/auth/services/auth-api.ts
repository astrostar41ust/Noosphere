import api, { setAccessToken } from '@/lib/api-client';
import { User, LoginCredentials, RegisterCredentials } from '../types';

export async function loginUser(credentials: LoginCredentials): Promise<User> {
  const response = await api.post<{ user: User; access_token: string }>('/api/v1/auth/login', credentials);
  const { user, access_token } = response.data;
  
  setAccessToken(access_token);
  return user;
}

export async function registerUser(credentials: RegisterCredentials): Promise<User> {
  const response = await api.post<{ user: User; access_token: string }>('/api/v1/auth/register', credentials);
  const { user, access_token } = response.data;
  
  setAccessToken(access_token);
  return user;
}

export async function refreshUserSession(): Promise<{ user: User; accessToken: string }> {
  const response = await api.post<{ user: User; access_token: string }>('/api/v1/auth/refresh');
  const { user, access_token } = response.data;
  
  setAccessToken(access_token);
  return { user, accessToken: access_token };
}

export async function logoutUserSession(): Promise<void> {
  await api.post('/api/v1/auth/logout');
  setAccessToken(null);
}
