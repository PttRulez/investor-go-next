'use client';

import { Button } from '@mui/material';
import {
  DefaultFormBox,
  FormText,
  FormDatePicker,
} from '@pttrulez/mui-based-ui';
import {
  ChangeHandler,
  Controller,
  SubmitHandler,
  useForm,
} from 'react-hook-form';
import { FC, useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import investorService from '@/axios/investor/investor.service';
import dayjs, { Dayjs } from 'dayjs';
import { MoexSearchAutocompleteOption } from '@/components/ui/StocksSearch/types';
import { getSecurityTypeFromMoexSecType } from '@/utils/helpers';
import { zodResolver } from '@hookform/resolvers/zod';
import { DealType, Exchange, MoexSecurityType } from '@/types/enums';
import {
  CreateMoexShareDealData,
  CreateMoexShareDealSchema,
} from '@/validation';
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
  const {
    clearErrors,
    control,
    formState,
    handleSubmit,
    resetField,
    setValue,
    watch,
  } = useForm<CreateMoexShareDealData>({
    defaultValues: {
      date: dayjs().toDate(),
      portfolioId,
      type: dealType,
    },
    resolver: zodResolver(CreateMoexShareDealSchema),
  });
  const watchAll = watch();
  const client = useQueryClient();
  const [ticker, setTicker] = useState<string | null>(null);

  // запрос на создание сделки
  const createDeal = useMutation(
    (formData: CreateMoexShareDealData) =>
      investorService.deal.createDeal(formData),
    {
      onSuccess: deal => {
        afterSuccessfulSubmit();
        client.invalidateQueries({ queryKey: ['portfolio', deal.portfolioId] });
      },
    },
  );

  const onSubmit: SubmitHandler<CreateMoexShareDealData> = async data => {
    if (!ticker) return;
    const shareInfo = await investorService.moexShare.getByTicker(ticker);

    data.securityId = shareInfo.id;
    createDeal.mutate(data);
  };

  let shareGroups: Array<MoexSecurityGroup> = [
    MoexSecurityGroup.stock_shares,
    MoexSecurityGroup.stock_dr,
  ];

  let bondGroups: Array<MoexSecurityGroup> = [
    MoexSecurityGroup.stock_bonds,
    MoexSecurityGroup.stock_eurobond,
  ];

  const onSecChange = async (
    _: React.SyntheticEvent,
    secInfo: MoexSearchAutocompleteOption | null,
  ) => {
    if (!secInfo) {
      console.log('BOMZHIHA !value');
      return;
    }

    if (shareGroups.includes(secInfo.group)) {
      let a = await investorService.moexShare.getByTicker(secInfo.ticker);
      console.log('received share info from our api', a);
    } else if (bondGroups.includes(secInfo.group)) {
      let a = await investorService.moexBond.getByTicker(secInfo.ticker);
      console.log('received bond info from our api', a);
    }
    console.log('BOMZHIHA Щсрщлщ', secInfo);
    setTicker(secInfo?.ticker ?? null);
  };

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <Controller
        control={control}
        name="securityId"
        render={({ field }) => (
          <MoexSearch
            onChange={onSecChange}
            error={!!formState.errors.securityId}
            helperText={formState.errors.securityId?.message}
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
