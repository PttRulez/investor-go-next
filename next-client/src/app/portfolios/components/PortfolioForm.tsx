'use client';
import { Box, Button, Stack } from '@mui/material';
import { useForm } from 'react-hook-form';
import { FormText, FormCheckBox } from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import { CreatePortfolioData, UpdatePortfolioData } from '@/validation';

interface PortfolioFormProps {
  portfolio: CreatePortfolioData | UpdatePortfolioData;
  afterSuccessfulSubmit: () => void;
}

export default function PortfolioForm({
  portfolio,
  afterSuccessfulSubmit,
}: PortfolioFormProps) {
  const { control, handleSubmit, setValue, watch } = useForm<
    CreatePortfolioData | UpdatePortfolioData
  >({
    defaultValues: portfolio,
  });
  const watchAll = watch();
  const client = useQueryClient();

  const createPortfolio = useMutation(
    (formData: CreatePortfolioData) =>
      investorService.portfolio.createPortfolio(formData),
    {
      onSuccess: () => {
        afterSuccessfulSubmit();
        client.invalidateQueries(['allPortfolios']);
      },
    },
  );

  const updatePortfolio = useMutation(
    (formData: UpdatePortfolioData) =>
      investorService.portfolio.updatePortfolio(formData),
    {
      onSuccess: () => {
        afterSuccessfulSubmit();
        client.invalidateQueries(['allPortfolios']);
      },
    },
  );

  const onSubmit = (formData: CreatePortfolioData | UpdatePortfolioData) => {
    if (!('id' in formData)) {
      createPortfolio.mutate(formData);
    } else {
      updatePortfolio.mutate({
        compound: formData.compound,
        name: formData.name,
        id: formData.id,
      });
    }
  };

  return (
    <Box onSubmit={handleSubmit(onSubmit)} sx={{ backgroundColor: 'white' }}>
      <Stack component="form" gap="20px">
        <FormText
          control={control}
          name="name"
          label="Название"
          value={watchAll.name}
          handleClear={() => setValue('name', '')}
        />
        <FormCheckBox
          // @ts-ignore
          control={control}
          name="compound"
          label="Составной"
          checked={watchAll.compound}
        />
        <Button
          variant="outlined"
          color="primary"
          type="submit"
          sx={{ color: 'grey.700' }}
        >
          Сохранить
        </Button>
      </Stack>
    </Box>
  );
}
