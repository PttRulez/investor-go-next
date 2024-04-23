'use client';

import investorService from '@/axios/investor/investor.service';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { getSession, useSession } from 'next-auth/react';

const client = new QueryClient();

const ReactQueryProvider = ({ children }: { children: React.ReactNode }) => {
  const { data: session } = useSession();
  if (session?.user?.token) {
    investorService.setToken(session?.user.token);
  }
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
};

export default ReactQueryProvider;
