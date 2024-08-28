import { SecurityType } from '@/types/enums';
import { IDealResponse } from './deal';
import { ITransactionResponse } from './transaction';

export interface IPortfolioResponse {
  bondDeals: IDealResponse[];
  bondPositions: IPositionResponse[];
  cash: number;
  cashouts: ITransactionResponse[];
  cashoutsSum: number;
  compound: boolean;
  deposits: ITransactionResponse[];
  depositsSum: number;
  id: number;
  name: string;
  profitability: number;
  shareDeals: IDealResponse[];
  sharePositions: IPositionResponse[];
  totalCost: number;
}

export interface IPositionResponse {
  amount: number;
  averagePrice: number;
  comment: string | null;
  id: number;
  currentPrice: number;
  currentCost: number;
  securityType: SecurityType;
  shortName: string;
  targetPrice: number | null;
  ticker: string;
}

export type IPortfolioListResponse = Pick<
  IPortfolioResponse,
  'id' | 'name' | 'compound'
>;

export interface IPortfolioPositionsResponse {
  allPositions: Array<IPositionResponse>;
  bondPositions: Array<IPositionResponse>;
  bondsTotal: number;
  // tradeSaldo?: number;
  sharePositions: Array<IPositionResponse>;
  sharesTotal: number;
}
