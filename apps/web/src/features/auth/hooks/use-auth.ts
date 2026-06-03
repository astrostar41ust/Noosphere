'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { loginUser, registerUser } from '../services/auth-api';
import { getCurrentUser, logoutUser } from '../services/auth-storage';
import { User } from '../types';

export function useCurrentUser() {
  return useQuery<User | null>({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    staleTime: Infinity, // The auth session remains stable until active logout/login
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: loginUser,
    onSuccess: (user) => {
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
      queryClient.setQueryData(['currentUser'], user);
      router.push('/chat');
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const router = useRouter();

  return () => {
    logoutUser();
    queryClient.setQueryData(['currentUser'], null);
    router.push('/login');
  };
}
