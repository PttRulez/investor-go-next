'use client';

import investorService from '@/axios/investor/investor.service';
import { Exchange } from '@/types/enums';
import {
  CreateCouponData,
  CreateCouponSchema,
} from '@/validation/coupon-schema';
import { zodResolver } from '@hookform/resolvers/zod';
import { Typography } from '@mui/material';
import Button from '@mui/material/Button/Button';
import {
  DefaultFormBox,
  FormDatePicker,
  FormSelect,
  FormText,
} from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs, { Dayjs } from 'dayjs';
import { useForm } from 'react-hook-form';

type Props = {
  afterSuccessfulSubmit: () => void;
  portfolioId: number;
  tickerList: SelectOption[] | SelectList;
};

const CreateCouponForm = ({
  afterSuccessfulSubmit,
  portfolioId,
  tickerList,
}: Props) => {
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateCouponData>({
      defaultValues: {
        date: dayjs().format('YYYY-MM-DD'),
        exchange: Exchange.MOEX,
        portfolioId,
      },
      resolver: zodResolver(CreateCouponSchema),
    });
  const watchAll = watch();
  const client = useQueryClient();

  // запрос на создание выплаты купона
  const createCoupon = useMutation({
    mutationFn: (formData: CreateCouponData) =>
      investorService.coupon.createCoupon(formData),
    onSuccess: _ => {
      afterSuccessfulSubmit();
      client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
    },
  });

  const onSubmit = (data: CreateCouponData) => {
    createCoupon.mutate(data);
  };
  return (
    <DefaultFormBox
      onSubmit={handleSubmit(onSubmit)}
      sx={{ minWidth: '600px' }}
    >
      <Typography variant="h6">Добавляем выплату купона</Typography>
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
          console.log(newValue?.toDate());
        }}
        label={'Дата выплаты'}
        value={watchAll.date}
      />
      <FormSelect
        control={control}
        name={'ticker'}
        label="Тикер (рег. номер)"
        options={tickerList}
        value={watchAll.ticker}
        variant="outlined"
      />
      <FormText
        control={control}
        decimal
        error={!!formState.errors.totalPayment}
        handleClear={() => resetField('totalPayment')}
        helperText={formState.errors.totalPayment?.message}
        label={'Общая сумма выплаты, пришедшая на счёт'}
        name={'totalPayment'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('totalPayment', Number(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.totalPayment}
      />
      <FormText
        control={control}
        decimal
        error={!!formState.errors.taxPaid}
        handleClear={() => resetField('taxPaid')}
        helperText={formState.errors.taxPaid?.message}
        label={'Налог уплаченный помимо суммы выплаты'}
        name={'taxPaid'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('taxPaid', Number(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.taxPaid}
      />
      <FormText
        control={control}
        error={!!formState.errors.bondsCount}
        handleClear={() => resetField('bondsCount')}
        helperText={formState.errors.bondsCount?.message}
        label={'Кол-во облигаций'}
        name={'bondsCount'}
        onChange={(e: any) => {
          setValue('bondsCount', parseInt(e.target.value));
        }}
        type="number"
        value={watchAll.bondsCount}
      />
      <FormText
        control={control}
        error={!!formState.errors.paymentPeriod}
        handleClear={() => setValue('paymentPeriod', '')}
        helperText={formState.errors.paymentPeriod?.message}
        label={'Период оплаты'}
        name={'paymentPeriod'}
        value={watchAll.paymentPeriod}
        multiline
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

export default CreateCouponForm;
