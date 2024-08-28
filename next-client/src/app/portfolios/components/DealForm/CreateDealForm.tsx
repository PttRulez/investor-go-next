'use client';

import { Button, TextField } from '@mui/material';
import { DefaultFormBox, FormText, FormDatePicker } from '@pttrulez';
import {
  ChangeHandler,
  Controller,
  SubmitHandler,
  useForm,
} from 'react-hook-form';
import { ChangeEvent, FC, SyntheticEvent, useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import investorService from '@/axios/investor/investor.service';
import dayjs, { Dayjs } from 'dayjs';
import { MoexSearchAutocompleteOption } from '@/components/ui/StocksSearch/types';
import { getSecurityTypeFromMoexSecType } from '@/utils/helpers';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  DealType,
  Exchange,
  MoexSecurityType,
  SecurityType,
} from '@/types/enums';
import { CreateDealData, CreateDealSchema } from '@/validation';
import { MoexSecurityGroup } from '@/types/apis/go-api';

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
        date: dayjs().toDate(),
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
    onSuccess: deal => {
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
      resetField('secid');
      resetField('securityType');
      return;
    }

    setValue('secid', secInfo.ticker);
    if (shareTypes.includes(secInfo.type)) {
      setValue('securityType', SecurityType.SHARE);
    } else if (bondTypes.includes(secInfo.type)) {
      setValue('securityType', SecurityType.BOND);
    } else {
      console.log('securityType', secInfo.type, 'не обработан');
      resetField('securityType');
    }
  };

  useEffect(() => {
    console.log('formState.errors', formState.errors);
  }, [formState.errors]);

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <Controller
        control={control}
        name="secid"
        render={({ field }) => (
          <MoexSearch
            onChange={onSecChange}
            error={!!formState.errors.secid}
            helperText={formState.errors.secid?.message}
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
        label={'Цена покупки'}
        name={'price'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('price', parseFloat(e.target.value));
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
        error={!!formState.errors.price}
        handleClear={() => resetField('price')}
        helperText={formState.errors.price?.message}
        label={'Комиссия'}
        name={'comission'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('comission', parseFloat(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.comission}
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
