import axios, { AxiosError, AxiosInstance, CreateAxiosDefaults } from 'axios';
import { urls } from '@/constants/common';
import { getSession, signOut } from 'next-auth/react';

const investorAxiosInstance: AxiosInstance = axios.create({
  baseURL: urls.investor,
  withCredentials: true,
} as CreateAxiosDefaults);

investorAxiosInstance.interceptors.request.use(async config => {
  const s = await getSession();
  config.headers['Accept'] = 'application/json';
  config.headers['Content-Type'] = 'application/json';
  config.headers['Authorization'] = `Bearer ${s?.user?.token}`;
  return config;
});

investorAxiosInstance.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error: AxiosError) {
    if (error.response?.status === 401) {
      signOut();
    }
    throw error;
  },
);

export default investorAxiosInstance;
