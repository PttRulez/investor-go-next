import { SubmitHandler, useForm } from 'react-hook-form';
import dayjs from 'dayjs';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { FormText, FormDatePicker, FormSelect } from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import investorService from '@/axios/investor/investor.service';
import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button';
import { CreateTransactionData, CreateTransactionSchema } from '@/validation';
import { TransactionType } from '@/types/enums';

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
        date: dayjs().toDate(),
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
      onSuccess: cashout => {
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
      sx={{ padding: '30px' }}
    >
      <Stack gap={'20px'}>
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
          label={'Сумма кэшаута'}
          name={'amount'}
          type="number"
          value={watchAll.amount}
        />
        <FormDatePicker
          control={control}
          name={'date'}
          handleClear={() => resetField('date')}
          onChange={newValue => {
            if (newValue) {
              setValue('date', newValue?.toDate());
            } else {
              resetField('date');
            }
          }}
          label={'Дата кэшаута'}
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
      </Stack>
    </Box>
  );
};

export default TransactionForm;
