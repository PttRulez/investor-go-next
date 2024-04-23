import { IDealResponse } from './deal';
import { ITransactionResponse } from './transaction';
import { SecurityResponse } from './security';

export interface IPortfolioResponse {
  cash: number;
  cashoutsSum: number;
  compound: boolean;
  deals: IDealResponse[];
  depositsSum: number;
  id: number;
  name: string;
  positions: IPortfolioPositionsResponse;
  profitability: string;
  totalCost: number;
  transactions: ITransactionResponse[];
}

export interface IPositionResponse {
  amount: number;
  comment: string | null;
  id: number;
  currentPrice: number;
  security: SecurityResponse;
  targetPrice: number | null;
  tradeSaldo: number;
  total: number;
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
