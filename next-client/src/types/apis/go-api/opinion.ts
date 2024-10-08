import { CreateOpinionData } from '@/validation';

export type IOpinionResponse = CreateOpinionData & {
  expert?: IExpertResponse;
  id: number;
};

export interface IExpertResponse {
  id: number;
  avatarUrl: string | null;
  name: string;
  opinions?: Array<IOpinionResponse>;
}
