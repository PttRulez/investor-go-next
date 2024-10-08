import { Exchange } from '@/types/enums';
import { IMoexSecurtiyResponse } from './moex';

export interface BaseSecurityResponse {
  exchange: Exchange;
}
export type SecurityResponse = BaseSecurityResponse & IMoexSecurtiyResponse;
