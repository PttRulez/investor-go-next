// https://next-auth.js.org/configuration/nextjs#advanced-usage
import { withAuth, NextRequestWithAuth } from 'next-auth/middleware';
import { NextResponse } from 'next/server';

export default withAuth(
  // 'withAuth' augments your 'Request' with the user's token.
  function middleware(request: NextRequestWithAuth) {
    if (request.nextUrl.pathname.startsWith('/extra') && request.nextauth.token?.role !== 'admin') {
      return NextResponse.rewrite(new URL('/denied', request.url));
    }

    if (
      request.nextUrl.pathname.startsWith('/client') &&
      request.nextauth.token?.role !== 'admin' &&
      request.nextauth.token?.role !== 'manager'
    ) {
      return NextResponse.rewrite(new URL('/denied', request.url));
    }
  },
  {
    pages: {
      signIn: '/login',
    },
    callbacks: {
      authorized: ({ req, token }) => {
        return !!token;
      },
    },
  },
);

// https://nextjs.org/docs/pages/building-your-application/routing/middleware#matcher
export const config = { matcher: ['/((?!register|api|login|$).*)'] };
