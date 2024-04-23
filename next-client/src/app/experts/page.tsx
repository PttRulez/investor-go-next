'use client';

import { Avatar, Button, Dialog } from '@mui/material';
import { useState } from 'react';
import ExpertForm from './components/ExpertForm/ExpertForm';
import { AdvancedTable, AdvancedTableColumn } from '@pttrulez/mui-based-ui';
import { useQuery } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import { useRouter } from 'next/navigation';
import { IExpertResponse } from '@/types/apis/go-api';

const columns: AdvancedTableColumn<IExpertResponse>[] = [
  {
    name: 'avatarUrl',
    label: '',
    format: (url: string) => {
      return <Avatar src={url ?? 'https://avatar.iran.liara.run/public'} />;
    },
  },
  {
    name: 'name',
    label: '',
  },
];

const ExpertsPage = () => {
  const [expertModalOpen, setExperModalOpen] = useState<boolean>(false);

  const { data: experts } = useQuery({
    queryKey: ['allExperts'],
    queryFn: () => investorService.expert.getExpertsList(),
  });

  const router = useRouter();

  return (
    <>
      <Button
        variant="outlined"
        sx={{
          color: 'grey.700',
          borderColor: 'grey.700',
          marginBottom: '50px',
        }}
        onClick={() => setExperModalOpen(true)}
      >
        + Эксперт
      </Button>
      <AdvancedTable
        columns={columns}
        rows={experts ?? []}
        rowClick={(row: IExpertResponse) => {
          router.push(`/experts/${row.id}`);
        }}
      />
      <Dialog open={expertModalOpen} onClose={() => setExperModalOpen(false)}>
        <ExpertForm afterSuccessfulSubmit={() => setExperModalOpen(false)} />
      </Dialog>
    </>
  );
};

export default ExpertsPage;
