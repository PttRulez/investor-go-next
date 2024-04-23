import { TransactionType } from '@/types/enums';

export interface ITransactionResponse {
  id: number;
  amount: number;
  type: TransactionType;
  date: Date;
}
