import { z } from 'zod';
import { Role } from '@/types/enums';

export const LoginSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

export type LoginData = z.infer<typeof LoginSchema>;

export const RegisterSchema = LoginSchema.extend({
  name: z.string(),
  role: z.nativeEnum(Role).optional(),
});

export type RegisterData = z.infer<typeof RegisterSchema>;
