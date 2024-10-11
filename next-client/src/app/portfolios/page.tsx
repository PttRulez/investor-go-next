'use client';

import { SyntheticEvent, useState } from 'react';
import { Dialog, IconButton } from '@mui/material';
import { AddNewButton } from '@pttrulez/mui-based-ui';
import PortfolioForm from './components/PortfolioForm';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import { AdvancedTable, AdvancedTableColumn } from '@pttrulez/mui-based-ui';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import DeleteIcon from '@mui/icons-material/Delete';
import ModeEditIcon from '@mui/icons-material/ModeEdit';
import { redirect, useRouter } from 'next/navigation';
import Head from 'next/head';
import { signOut, useSession } from 'next-auth/react';
import { IPortfolioResponse } from '@/types/apis/go-api';
import { CreatePortfolioData, UpdatePortfolioData } from '@/validation';

const emptyPortfolio: CreatePortfolioData = {
  name: '',
  compound: false,
};

export default function PortfoliosPage() {
  const { data: session } = useSession({
    required: true,
    onUnauthenticated() {
      signOut();
      // redirect('/api/auth/signin?callbackUrl=/portfolios');
    },
  });
  const client = useQueryClient();
  const [portfolioToEdit, setPortfolioToEdit] = useState<
    CreatePortfolioData | UpdatePortfolioData | null
  >(null);

  const { data: portfolioList } = useQuery({
    queryKey: ['allPortfolios'],
    queryFn: () => investorService.portfolio.allPortfolios(),
  });
  const router = useRouter();

  const addNewPortfolio = () => {
    setPortfolioToEdit(emptyPortfolio);
  };

  const deletePortfolio = useMutation(
    (id: number) => investorService.portfolio.deletePortfolio(id),
    {
      onSuccess: () => {
        client.invalidateQueries(['allPortfolios']);
      },
    },
  );

  const columns: AdvancedTableColumn<IPortfolioResponse>[] = [
    {
      name: 'name',
      label: 'Название',
    },
    {
      name: 'compound',
      label: 'Составной',
      format: (compound: boolean) => {
        return compound ? (
          <CheckIcon sx={{ color: 'success.main' }} />
        ) : (
          <CloseIcon sx={{ color: 'error.main' }} />
        );
      },
      align: 'center',
    },
    {
      name: 'actions',
      label: '',
      render: (_, portfolio) => {
        return (
          <>
            <IconButton
              onClick={(e: SyntheticEvent) => {
                e.stopPropagation();
                setPortfolioToEdit({
                  compound: portfolio.compound,
                  id: portfolio.id,
                  name: portfolio.name,
                });
              }}
            >
              <ModeEditIcon />
            </IconButton>
            <IconButton
              onClick={(e: SyntheticEvent) => {
                e.stopPropagation();
                deletePortfolio.mutate(portfolio.id);
              }}
            >
              <DeleteIcon />
            </IconButton>
          </>
        );
      },
    },
  ];

  // if (!session?.user) return;

  return (
    <>
      <Head>
        <title>My page title</title>
      </Head>
      <AdvancedTable
        columns={columns}
        rows={portfolioList ?? []}
        pagination={false}
        rowClick={row => {
          router.push(`/portfolios/${row.id}`);
        }}
      />
      <AddNewButton onClick={addNewPortfolio} />
      <Dialog
        open={!!portfolioToEdit}
        onClose={() => setPortfolioToEdit(null)}
        sx={{
          '.MuiPaper-root.MuiDialog-paper': {
            maxWidth: '90%',
            maxHeight: '90%',
            display: 'flex',
            justifyContent: 'center',
            padding: '20px',
          },
        }}
      >
        {portfolioToEdit && (
          <PortfolioForm
            afterSuccessfulSubmit={() => setPortfolioToEdit(null)}
            portfolio={portfolioToEdit}
          />
        )}
      </Dialog>
    </>
  );
}
