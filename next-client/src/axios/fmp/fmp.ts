import axios, {AxiosInstance, CreateAxiosDefaults} from 'axios';
import {urls} from '@/constants/common';
export const fmpApiV3: AxiosInstance = axios.create({
  // baseURL: 'https://financialmodelingprep.com/api/v3',
  baseURL: urls.fmp3,
  params: {
    apikey: '91fce00e525c655128402dd927de30d6',
  },
} as CreateAxiosDefaults);

fmpApiV3.interceptors.request.use((config) => {
  config.headers['Accept'] = 'application/json';
  config.headers['Content-Type'] = 'application/json';
  return config;
});
