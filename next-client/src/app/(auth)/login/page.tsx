'use client';
import { FC, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { Box, TextField } from '@mui/material';
import { LoadingButton } from '@mui/lab';
import { signIn } from 'next-auth/react';
import { LoginData } from '@/validation';
import { useRouter } from 'next/navigation';

const Login: FC = () => {
  const { handleSubmit, register } = useForm<LoginData>();
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const onSubmit: SubmitHandler<LoginData> = async data => {
    setLoading(true);
    try {
      const loginData = await signIn('credentials', {
        email: data.email,
        password: data.password,
        redirect: false,
      });

      setLoading(false);
      if (loginData?.ok) {
        router.push('/');
      }
    } catch (e) {
      setLoading(false);
    }
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
        <TextField label="Логин" variant="standard" {...register('email')} />
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
          loading={loading}
        >
          Логин
        </LoadingButton>
      </Box>
    </Box>
  );
};

export default Login;
