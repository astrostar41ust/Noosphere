export interface User {
  id: string;
  username: string;
  email: string;
}

export type LoginCredentials = Pick<User, 'email'> & {
  password: string;
};

export type RegisterCredentials = Omit<User, 'id'> & {
  password: string;
  confirmPassword: string;
};

export type StoredUserCredentials = User & {
  password: string;
  confirmPassword: string;
};
