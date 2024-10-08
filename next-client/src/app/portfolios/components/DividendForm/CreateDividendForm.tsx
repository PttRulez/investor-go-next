'use client';

import investorService from '@/axios/investor/investor.service';
import { Exchange } from '@/types/enums';
import { CreateDividendData, CreateDividendSchema } from '@/validation';
import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button/Button';
import {
  DefaultFormBox,
  FormDatePicker,
  FormSelect,
  FormText,
} from '@pttrulez';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs, { Dayjs } from 'dayjs';
import { useForm } from 'react-hook-form';

type Props = {
  afterSuccessfulSubmit: () => void;
  portfolioId: number;
  tickerList: SelectOption[];
};

const CreateDividendForm = ({
  afterSuccessfulSubmit,
  portfolioId,
  tickerList,
}: Props) => {
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateDividendData>({
      defaultValues: {
        date: dayjs().format('YYYY-MM-DD'),
        exchange: Exchange.MOEX,
        portfolioId,
      },
      resolver: zodResolver(CreateDividendSchema),
    });
  const watchAll = watch();
  const client = useQueryClient();

  // запрос на создание выплаты дивиденда
  const createDividend = useMutation({
    mutationFn: (formData: CreateDividendData) =>
      investorService.dividend.createDividend(formData),
    onSuccess: d => {
      afterSuccessfulSubmit();
      client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
    },
  });

  const onSubmit = (data: CreateDividendData) => {
    createDividend.mutate(data);
  };
  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
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
        label="Тикер"
        options={tickerList}
        value={watchAll.ticker}
        variant="outlined"
      />
      <FormText
        control={control}
        decimal
        error={!!formState.errors.paymentPerShare}
        handleClear={() => resetField('paymentPerShare')}
        helperText={formState.errors.paymentPerShare?.message}
        label={'Размер дивиденда на одну акцию'}
        name={'paymentPerShare'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('paymentPerShare', parseFloat(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.paymentPerShare}
      />
      <FormText
        control={control}
        error={!!formState.errors.sharesCount}
        handleClear={() => resetField('sharesCount')}
        helperText={formState.errors.sharesCount?.message}
        label={'Кол-во акций'}
        name={'sharesCount'}
        onChange={(e: any) => {
          setValue('sharesCount', parseInt(e.target.value));
        }}
        type="number"
        value={watchAll.sharesCount}
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

export default CreateDividendForm;
