import { Exchange } from '@/types/enums';
import { IMoexSecurtiyResponse } from './moex-securities';

export interface BaseSecurityResponse {
  exchange: Exchange;
}
export type SecurityResponse = BaseSecurityResponse & IMoexSecurtiyResponse;
