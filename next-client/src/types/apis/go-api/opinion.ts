import { CreateOpinionData } from '@/validation';
import { IExpertResponse } from './expert';

export type IOpinionResponse = CreateOpinionData & {
  expert?: IExpertResponse;
};
