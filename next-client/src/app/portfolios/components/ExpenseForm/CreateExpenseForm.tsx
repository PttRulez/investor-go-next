'use client';

import investorService from '@/axios/investor/investor.service';
import {
  CreateExpenseData,
  CreateExpenseSchema,
} from '@/validation/expense-schema';
import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button/Button';
import {
  DefaultFormBox,
  FormDatePicker,
  FormText,
} from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs, { Dayjs } from 'dayjs';
import { useForm } from 'react-hook-form';

type Props = {
  afterSuccessfulSubmit: () => void;
  portfolioId: number;
};

const CreateExpenseForm = ({ afterSuccessfulSubmit, portfolioId }: Props) => {
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateExpenseData>({
      defaultValues: {
        date: dayjs().format('YYYY-MM-DD'),
        portfolioId,
      },
      resolver: zodResolver(CreateExpenseSchema),
    });
  const watchAll = watch();
  const client = useQueryClient();

  // запрос на создание траты
  const createExpense = useMutation({
    mutationFn: (formData: CreateExpenseData) =>
      investorService.portfolio.createExpense(formData),
    onSuccess: c => {
      afterSuccessfulSubmit();
      client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
    },
  });

  const onSubmit = (data: CreateExpenseData) => {
    createExpense.mutate(data);
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
      <FormText
        control={control}
        decimal
        error={!!formState.errors.amount}
        handleClear={() => resetField('amount')}
        helperText={formState.errors.amount?.message}
        label={'Размер траты'}
        name={'amount'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('amount', Number(e.target.value));
          }
        }}
        type="number"
        inputProps={{
          step: 'any',
        }}
        value={watchAll.amount}
      />
      <FormText
        control={control}
        error={!!formState.errors.description}
        handleClear={() => resetField('description')}
        helperText={formState.errors.description?.message}
        label={'Описание'}
        name={'description'}
        value={watchAll.description}
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

export default CreateExpenseForm;
