import * as z from 'zod';

export const loginFormSchema = z.object({
  email: z.string().trim().email('Please enter a valid email address format'),
  password: z.string().min(6, 'Password credentials must exceed 6 characters'),
});

export type LoginFormData = z.infer<typeof loginFormSchema>;

export const registerFormSchema = z
  .object({
    username: z
      .string()
      .trim()
      .min(3, 'Username node name must exceed 3 characters')
      .max(20, 'Username node name exceeds 20 characters limit'),
    email: z.string().trim().email('Please enter a valid email address format'),
    password: z.string().min(6, 'Password credentials must exceed 6 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passphrase verification nodes do not match',
    path: ['confirmPassword'],
  });

export type RegisterFormData = z.infer<typeof registerFormSchema>;
