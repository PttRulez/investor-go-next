'use client';

import { Button } from '@mui/material';
import {
  DefaultFormBox,
  FormText,
  FormDatePicker,
} from '@pttrulez/mui-based-ui';
import { Controller, SubmitHandler, useForm } from 'react-hook-form';
import { FC, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import investorService from '@/axios/investor/investor.service';
import dayjs, { Dayjs } from 'dayjs';
import { MoexSearchAutocompleteOption } from '@/components/ui/StocksSearch/types';
import { getSecurityTypeFromMoexSecType } from '@/utils/helpers';
import { zodResolver } from '@hookform/resolvers/zod';
import { DealType, Exchange } from '@/types/enums';
import { CreateDealData, CreateDealSchema } from '@/validation';

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
  const {
    clearErrors,
    control,
    formState,
    handleSubmit,
    resetField,
    setValue,
    watch,
  } = useForm<CreateDealData>({
    defaultValues: {
      date: dayjs().toDate(),
      portfolioId,
      type: dealType,
      securityType: undefined,
      exchange: undefined,
      ticker: undefined,
    },
    resolver: zodResolver(CreateDealSchema),
  });
  const watchAll = watch();
  const client = useQueryClient();

  const changeMoexSecurityHandler = (
    value: MoexSearchAutocompleteOption | null,
  ) => {
    if (value) {
      setValue('exchange', Exchange.MOEX);
      setValue('securityType', getSecurityTypeFromMoexSecType(value.type));
      setValue('ticker', value.ticker);
      clearErrors('ticker');
    } else {
      resetField('ticker');
    }
  };

  const createDeal = useMutation(
    (formData: CreateDealData) => investorService.deal.createDeal(formData),
    {
      onSuccess: deal => {
        afterSuccessfulSubmit();
        client.invalidateQueries({ queryKey: ['portfolio', deal.portfolioId] });
      },
    },
  );

  const onSubmit: SubmitHandler<CreateDealData> = data => {
    createDeal.mutate(data);
  };

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <Controller
        control={control}
        name="ticker"
        render={({ field }) => (
          <MoexSearch
            onChange={(e, value) => {
              // field.onChange(e);
              changeMoexSecurityHandler(value);
            }}
            error={!!formState.errors.ticker}
            helperText={formState.errors.ticker?.message}
          />
        )}
      />
      <FormText
        control={control}
        error={!!formState.errors.amount}
        handleClear={() => setValue('amount', 0)}
        helperText={formState.errors.amount?.message}
        label={'Кол-во бумаг'}
        name={'amount'}
        type="number"
        value={watchAll.amount}
      />
      <FormText
        control={control}
        decimal
        error={!!formState.errors.price}
        handleClear={() => setValue('price', 0)}
        helperText={formState.errors.price?.message}
        label={'Цена покупки'}
        name={'price'}
        type="number"
        value={watchAll.price}
      />
      <FormDatePicker
        control={control}
        name={'date'}
        handleClear={() => resetField('date')}
        onChange={(newValue: Dayjs | null) => {
          if (newValue) {
            setValue('date', newValue?.toDate());
          } else {
            resetField('date');
          }
        }}
        label={'Дата покупки'}
        value={watchAll.date}
      />
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
