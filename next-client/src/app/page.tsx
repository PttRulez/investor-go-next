'use client';

import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import { MoexSearchHandler } from '@/components/ui/StocksSearch/types';
import Box from '@mui/material/Box';
import { redirect } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { getSecurityTypeFromMoexSecType } from '@/utils/helpers';
import { SecurityType } from '@/types/enums';

export default function Home() {
  const { data: session } = useSession({
    required: true,
    onUnauthenticated() {
      redirect('/api/auth/signin?callb ackUrl=/');
    },
  });

  const router = useRouter();
  console.log('session', session);
  if (!session?.user) return;

  const changeHandler: MoexSearchHandler = async (e, value, reason) => {
    if (value) {
      const type = getSecurityTypeFromMoexSecType(value.type);
      const url = type === SecurityType.SHARE ? 'shares' : 'bonds';
      router.push(`/${url}/moex/${value.ticker}`);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
      }}
    >
      <MoexSearch onChange={changeHandler} />
    </Box>
  );
}
