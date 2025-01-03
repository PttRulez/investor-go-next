import axios from 'axios';
import type { NextAuthOptions } from 'next-auth';
import NextAuth from 'next-auth/next';
import CredentialsProvider from 'next-auth/providers/credentials';

const options: NextAuthOptions = {
  pages: {
    signIn: '/login',
    signOut: '/auth/signout',
  },
  providers: [
    CredentialsProvider({
      name: 'Credentials',
      credentials: {
        email: {
          label: 'Email:',
          type: 'email',
          placeholder: 'Введите email',
        },
        password: {
          label: 'password:',
          type: 'password',
          placeholder: 'Напечатайте пароль',
        },
      },
      async authorize(credentials) {
        if (!credentials) return null;

        try {
          const response = await axios.post(
            `${process.env.NEXT_PUBLIC_INVESTOR_API_INTERNAL_URL!}/login`,
            credentials,
          );
          if (!response) {
            return null;
          }
          return response.data;
        } catch (e: any) {
          throw new Error(e.response.data);
        }
      },
    }),
  ],
  callbacks: {
    // https://authjs.dev/guides/basics/role-based-access-control#persisting-the-role
    async jwt({ token, user, profile }) {
      if (user) {
        token.role = user.role;
        token.token = user.token;
      }
      return token;
    },
    // If you want to use the role in client components
    async session({ session, token, user }) {
      if (session?.user) {
        session.user.role = token.role;
        session.user.token = token.token;
      }
      return session;
    },
  },
  logger: {
    error(code, metadata) {
      console.error('[NEXT_AUTH LOGGER]:', code, metadata);
    },
  },
};

const handler = NextAuth(options);

export default handler;
