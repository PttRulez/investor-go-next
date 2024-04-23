'use client';
import { FC } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { Box, TextField } from '@mui/material';
import { LoadingButton } from '@mui/lab';
import investorService from '@/axios/investor/investor.service';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { RegisterData } from '@/validation';

const Register: FC = () => {
  const { handleSubmit, register } = useForm<RegisterData>();
  const router = useRouter();

  const registerUser = useMutation(
    (formData: RegisterData) => investorService.auth.register(formData),
    {
      onSuccess: () => {
        router.push('/login');
      },
    },
  );

  const onSubmit: SubmitHandler<RegisterData> = async data => {
    registerUser.mutate(data);
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100%',
        gap: 1,
      }}
    >
      <Box
        component="form"
        onSubmit={handleSubmit(onSubmit)}
        sx={{
          display: 'flex',
          flexDirection: 'column',
          rowGap: 4,
          minWidth: 500,
          border: '1px solid',
          borderColor: 'grey.A400',
          borderRadius: 6,
          padding: 4,
        }}
      >
        <TextField
          label="Имя или никнейм"
          variant="standard"
          {...register('name')}
        />
        <TextField label="Email" variant="standard" {...register('email')} />
        <TextField
          label="Пароль"
          variant="standard"
          type="password"
          {...register('password')}
        />
        <LoadingButton
          variant="outlined"
          color="inherit"
          type="submit"
          loading={registerUser.isLoading}
        >
          Зарегистрироваться
        </LoadingButton>
      </Box>
    </Box>
  );
};

export default Register;
