import { SubmitHandler, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import { DefaultFormBox, FormText } from '@pttrulez/mui-based-ui';
import { Button } from '@mui/material';
import { CreateExpertData, CreateExpertSchema } from '@/validation';

interface ExpertFormProps {
  afterSuccessfulSubmit: () => void;
}

const ExpertForm = ({ afterSuccessfulSubmit }: ExpertFormProps) => {
  const {
    clearErrors,
    control,
    formState,
    handleSubmit,
    resetField,
    setValue,
    watch,
  } = useForm<CreateExpertData>({
    defaultValues: {
      avatarUrl: undefined,
      name: undefined,
    },
    resolver: zodResolver(CreateExpertSchema),
  });

  const watchAll = watch();
  const client = useQueryClient();

  const createExpert = useMutation(
    (formData: CreateExpertData) =>
      investorService.expert.createExpert(formData),
    {
      onSuccess: expert => {
        afterSuccessfulSubmit();
        client.invalidateQueries(['allExperts']);
      },
    },
  );

  const onSubmit: SubmitHandler<CreateExpertData> = data => {
    createExpert.mutate(data);
  };

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <FormText
        control={control}
        error={!!formState.errors.name}
        handleClear={() => setValue('name', '')}
        helperText={formState.errors.name?.message}
        label={'Имя эксперта'}
        name={'name'}
        value={watchAll.name}
      />
      <FormText
        control={control}
        error={!!formState.errors.avatarUrl}
        handleClear={() => setValue('avatarUrl', undefined)}
        helperText={formState.errors.avatarUrl?.message}
        label={'Ссылка на аватарку'}
        name={'avatarUrl'}
        value={watchAll.avatarUrl}
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

export default ExpertForm;
