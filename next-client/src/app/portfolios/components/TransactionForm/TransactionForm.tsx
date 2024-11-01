import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { FormDatePicker, FormSelect, FormText } from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs from 'dayjs';
import { SubmitHandler, useForm } from 'react-hook-form';

import investorService from '@/axios/investor/investor.service';
import { TransactionType } from '@/types/enums';
import { CreateTransactionData, CreateTransactionSchema } from '@/validation';
import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button';
import { useEffect } from 'react';

type Props = {
  afterSuccessfulSubmit: () => void;
  portfolioId: number;
};

const TransactionForm = ({ afterSuccessfulSubmit, portfolioId }: Props) => {
  const client = useQueryClient();
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateTransactionData>({
      defaultValues: {
        amount: undefined,
        date: dayjs().format('YYYY-MM-DD'),
        portfolioId,
        type: undefined,
      },
      resolver: zodResolver(CreateTransactionSchema),
    });

  const watchAll = watch();

  const createTransaction = useMutation(
    (formData: CreateTransactionData) =>
      investorService.transaction.createTransaction(formData),
    {
      onSuccess: _ => {
        afterSuccessfulSubmit();
        client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
      },
    },
  );

  const onSubmit: SubmitHandler<CreateTransactionData> = data => {
    createTransaction.mutate(data);
  };

  return (
    <Box
      onSubmit={handleSubmit(onSubmit)}
      component="form"
      sx={{ padding: '30px', minWidth: '400px' }}
    >
      <Stack gap={'20px'}>
        <FormDatePicker
          control={control}
          name={'date'}
          handleClear={() => resetField('date')}
          onChange={newValue => {
            if (newValue) {
              setValue('date', newValue?.format('YYYY-MM-DD'));
            } else {
              resetField('date');
            }
          }}
          label={'Дата'}
          value={watchAll.date}
        />
        <FormSelect
          control={control}
          name={'type'}
          label="Кэш или деп ?"
          options={{
            [TransactionType.CASHOUT]: 'Кэшаут',
            [TransactionType.DEPOSIT]: 'Депозит',
          }}
          value={watchAll.type}
        />
        <FormText
          control={control}
          error={!!formState.errors.amount}
          handleClear={() => resetField('amount')}
          helperText={formState.errors.amount?.message}
          label={
            watchAll.type === TransactionType.CASHOUT
              ? 'Сумма кэшаута'
              : 'Сумма депозита'
          }
          name={'amount'}
          onChange={e => {
            if (e.target.value != '') {
              setValue('amount', Number(e.target.value));
            }
          }}
          type="number"
          value={watchAll.amount}
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
};

export default TransactionForm;
