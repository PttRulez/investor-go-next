import investorService from '@/axios/investor/investor.service';
import { UpdatePositionData } from '@/validation';
import { Button } from '@mui/material';
import { DefaultFormBox, FormText } from '@pttrulez/mui-based-ui';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { SubmitHandler, useForm } from 'react-hook-form';

type Props = {
  afterSuccessfulSubmit: () => void;
  id: number;
  position: UpdatePositionData;
  portfolioId: number;
};

const PositionForm = ({
  afterSuccessfulSubmit,
  id,
  position,
  portfolioId,
}: Props) => {
  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<UpdatePositionData & { opinions?: { id: number }[] }>({
      defaultValues: {
        comment: position.comment,
        targetPrice: position.targetPrice,
      },
    });
  const watchAll = watch();

  const queryClient = useQueryClient();

  const updatePosition = useMutation(
    (formData: UpdatePositionData) =>
      investorService.position.update(id, formData),
    {
      onSuccess: deal => {
        afterSuccessfulSubmit();
        queryClient.invalidateQueries({
          queryKey: ['portfolio', portfolioId],
        });
      },
    },
  );

  const onSubmit: SubmitHandler<UpdatePositionData> = data => {
    updatePosition.mutate(data);
  };

  return (
    <DefaultFormBox
      onSubmit={handleSubmit(onSubmit)}
      sx={{
        minWidth: '800px',
      }}
    >
      <FormText
        control={control}
        error={!!formState.errors.comment}
        handleClear={() => setValue('comment', '')}
        helperText={formState.errors.comment?.message}
        label={'Комент'}
        name={'comment'}
        value={watchAll.comment}
        multiline
      />
      <FormText
        control={control}
        error={!!formState.errors.targetPrice}
        handleClear={() => setValue('targetPrice', 0)}
        helperText={formState.errors.targetPrice?.message}
        label={'Целевая цена'}
        name={'targetPrice'}
        onChange={(e: any) => {
          if (e.target.value != '') {
            setValue('targetPrice', parseFloat(e.target.value));
          }
        }}
        type="number"
        value={watchAll.targetPrice}
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

export default PositionForm;
