import { DefaultSession, DefaultUser } from 'next-auth';
import { JWT, DefaultJWT } from 'next-auth/jwt';
import type { Role } from '@/types/backend';

declare module 'next-auth' {
  interface Session {
    user: {
      id: string;
      role: string;
      token: string;
    } & DefaultSession;
  }

  interface User extends DefaultUser {
    role: Role;
    id: number;
    token: string;
  }
}

declare module 'next-auth/jwt' {
  interface JWT extends DefaultJWT {
    role: string;
    token: string;
  }
}
