import { User, LoginCredentials, RegisterCredentials, StoredUserCredentials } from '../types';
import { getRegisteredUsers, saveRegisteredUser, setCurrentUser } from './auth-storage';

export async function loginUser(credentials: LoginCredentials): Promise<User> {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 1000));

  const users = getRegisteredUsers();
  const matchedUser = users.find((u) => u.email.toLowerCase() === credentials.email.toLowerCase());

  if (!matchedUser) {
    throw new Error('No user profile matching this email address was found.');
  }

  if (matchedUser.password !== credentials.password) {
    throw new Error('Authentication credentials validation rejected.');
  }

  const activeUser: User = {
    id: matchedUser.id,
    username: matchedUser.username,
    email: matchedUser.email,
  };

  setCurrentUser(activeUser);

  return activeUser;
}

export async function registerUser(credentials: RegisterCredentials): Promise<User> {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 1000));

  const users = getRegisteredUsers();
  const emailExists = users.some((u) => u.email.toLowerCase() === credentials.email.toLowerCase());
  const usernameExists = users.some((u) => u.username.toLowerCase() === credentials.username.toLowerCase());

  if (emailExists) {
    throw new Error('A user profile is already registered with this email address.');
  }

  if (usernameExists) {
    throw new Error('Username has already been claimed by another node.');
  }

  const newUser: User = {
    id: crypto.randomUUID(),
    username: credentials.username,
    email: credentials.email,
  };

  const fullCredentialEntry: StoredUserCredentials = {
    ...newUser,
    password: credentials.password,
    confirmPassword: credentials.confirmPassword,
  };

  saveRegisteredUser(fullCredentialEntry);
  setCurrentUser(newUser);

  return newUser;
}
