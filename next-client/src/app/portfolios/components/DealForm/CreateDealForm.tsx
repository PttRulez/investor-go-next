'use client';

import investorService from '@/axios/investor/investor.service';
import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import { MoexSearchAutocompleteOption } from '@/components/ui/StocksSearch/types';
import {
  DealType,
  Exchange,
  MoexSecurityType,
  SecurityType,
} from '@/types/enums';
import { CreateDealData, CreateDealSchema } from '@/validation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Box, Button, Typography } from '@mui/material';
import {
  DefaultFormBox,
  FormDatePicker,
  FormText,
} from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs, { Dayjs } from 'dayjs';
import { FC, useEffect } from 'react';
import { Controller, SubmitHandler, useForm } from 'react-hook-form';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import ArrowCircleDownIcon from '@mui/icons-material/ArrowCircleDown';

interface DealFormProps {
  afterSuccessfulSubmit: () => void;
  dealType: DealType;
  portfolioId: number;
}

const CreateDealForm: FC<DealFormProps> = ({
  afterSuccessfulSubmit,
  dealType,
  portfolioId,
}) => {
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateDealData>({
      defaultValues: {
        comission: 0,
        date: dayjs().format('YYYY-MM-DD'),
        exchange: Exchange.MOEX,
        portfolioId,
        type: dealType,
      },
      resolver: zodResolver(CreateDealSchema),
    });
  const watchAll = watch();
  const client = useQueryClient();

  // запрос на создание сделки
  const createDeal = useMutation({
    mutationFn: (formData: CreateDealData) =>
      investorService.deal.createDeal(formData),
    onSuccess: _ => {
      afterSuccessfulSubmit();
      client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
    },
  });

  const onSubmit: SubmitHandler<CreateDealData> = async data => {
    createDeal.mutate(data);
  };

  const shareTypes: Array<MoexSecurityType> = [
    MoexSecurityType.common_share,
    MoexSecurityType.depositary_receipt,
    MoexSecurityType.preferred_share,
  ];

  const bondTypes: Array<MoexSecurityType> = [
    MoexSecurityType.corporate_bond,
    MoexSecurityType.exchange_bond,
    MoexSecurityType.ofz_bond,
  ];

  const onSecChange = async (
    _: React.SyntheticEvent,
    secInfo: MoexSearchAutocompleteOption | null,
  ) => {
    if (!secInfo) {
      resetField('ticker');
      resetField('securityType');
      return;
    }

    setValue('ticker', secInfo.ticker);
    setValue('shortName', secInfo.shortName);

    if (shareTypes.includes(secInfo.type)) {
      setValue('securityType', SecurityType.SHARE);
    } else if (bondTypes.includes(secInfo.type)) {
      setValue('securityType', SecurityType.BOND);
    } else {
      resetField('securityType');
    }
  };

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <Box
        sx={{
          alignItems: 'center',
          display: 'flex',
          gap: '20px',
          justifyContent: 'center',
        }}
      >
        <Typography variant="h6">
          Добавляем {dealType === DealType.SELL ? 'продажу' : 'покупку'} бумаги{' '}
        </Typography>
        {dealType === DealType.SELL ? (
          <ArrowCircleDownIcon sx={{ color: 'red' }} />
        ) : (
          <ArrowCircleUpIcon sx={{ color: 'green' }} />
        )}
      </Box>
      <FormDatePicker
        control={control}
        name={'date'}
        handleClear={() => resetField('date')}
        onChange={(newValue: Dayjs | null) => {
          if (newValue) {
            setValue('date', newValue?.format('YYYY-MM-DD'));
          } else {
            resetField('date');
          }
        }}
        label={'Дата'}
        value={watchAll.date}
      />
      <Controller
        control={control}
        name="ticker"
        render={({ field }) => (
          <MoexSearch
            onChange={onSecChange}
            error={!!formState.errors.ticker}
            helperText={formState.errors.ticker?.message}
          />
        )}
      />
      <FormText
        control={control}
        error={!!formState.errors.amount}
        handleClear={() => resetField('amount')}
        helperText={formState.errors.amount?.message}
        label={'Кол-во бумаг'}
        name={'amount'}
        onChange={(e: any) => {
          setValue('amount', parseInt(e.target.value));
        }}
        type="number"
        value={watchAll.amount}
      />

      <FormText
        control={control}
        decimal
        error={!!formState.errors.price}
        handleClear={() => resetField('price')}
        helperText={formState.errors.price?.message}
        label={dealType === DealType.BUY ? 'Цена покупки' : 'Цена продажи'}
        name={'price'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('price', Number(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.price}
      />
      <FormText
        control={control}
        decimal
        error={!!formState.errors.comission}
        handleClear={() => resetField('comission')}
        helperText={formState.errors.comission?.message}
        label={'Комиссия'}
        name={'comission'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('comission', Number(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.comission}
      />
      {watchAll.securityType === SecurityType.BOND && (
        <FormText
          control={control}
          error={!!formState.errors.nkd}
          handleClear={() => resetField('nkd')}
          helperText={formState.errors.nkd?.message}
          label={'НКД'}
          name={'nkd'}
          onChange={(e: any) => {
            if (e.target.value != '') {
              setValue('nkd', Number(e.target.value));
            }
          }}
          type="number"
          value={watchAll.nkd}
        />
      )}
      <Button
        variant="outlined"
        color="primary"
        type="submit"
        sx={{ color: 'grey.700' }}
      >
        Сохранить
      </Button>
    </DefaultFormBox>
  );
};

export default CreateDealForm;
