'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { loginUser, registerUser, refreshUserSession, logoutUserSession } from '../services/auth-api';
import { getCurrentUser, setCurrentUser, logoutUser } from '../services/auth-storage';
import { User } from '../types';

export function useCurrentUser() {
  return useQuery<User | null>({
    queryKey: ['currentUser'],
    queryFn: async () => {
      const cachedUser = getCurrentUser();
      if (!cachedUser) return null;

      try {
        const session = await refreshUserSession();
        setCurrentUser(session.user);
        return session.user;
      } catch (error) {
        logoutUser();
        return null;
      }
    },
    staleTime: Infinity,
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: loginUser,
    onSuccess: (user) => {
      setCurrentUser(user);
      queryClient.setQueryData(['currentUser'], user);
      router.push('/chat');
    },
  });
}

export function useRegister() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: registerUser,
    onSuccess: (user) => {
      setCurrentUser(user);
      queryClient.setQueryData(['currentUser'], user);
      router.push('/chat');
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const router = useRouter();

  const mutation = useMutation({
    mutationFn: logoutUserSession,
    onSuccess: () => {
      logoutUser();
      queryClient.setQueryData(['currentUser'], null);
      router.push('/login');
    },
    onError: () => {
      // Invalidate locally even if api call fails
      logoutUser();
      queryClient.setQueryData(['currentUser'], null);
      router.push('/login');
    },
  });

  return () => {
    mutation.mutate();
  };
}
