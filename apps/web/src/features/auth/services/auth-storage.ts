import { User, StoredUserCredentials } from '../types';

const PROFILE_KEY = 'user_profile';
const REGISTERED_USERS_KEY = 'noosphere_registered_users';

export function getRegisteredUsers(): StoredUserCredentials[] {
  if (typeof window === 'undefined') return [];
  const stored = localStorage.getItem(REGISTERED_USERS_KEY);
  return stored ? JSON.parse(stored) : [];
}

export function saveRegisteredUser(user: StoredUserCredentials): void {
  if (typeof window === 'undefined') return;
  const current = getRegisteredUsers();
  localStorage.setItem(REGISTERED_USERS_KEY, JSON.stringify([...current, user]));
}

export function getCurrentUser(): User | null {
  if (typeof window === 'undefined') return null;
  const stored = localStorage.getItem(PROFILE_KEY);
  return stored ? JSON.parse(stored) : null;
}

export function setCurrentUser(user: User): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(PROFILE_KEY, JSON.stringify(user));
}

export function logoutUser(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(PROFILE_KEY);
}
