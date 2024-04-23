import axios, { AxiosInstance, CreateAxiosDefaults } from 'axios';
import { urls } from '@/constants/common';

export const moexApi: AxiosInstance = axios.create({
  baseURL: urls.moex,
  params: {
    ['iss.meta']: 'off',
  },
} as CreateAxiosDefaults);

moexApi.interceptors.request.use(config => {
  config.headers['Accept'] = 'application/json';
  config.headers['Content-Type'] = 'application/json';
  return config;
});
