import { IDealResponse } from './deal';
import { ITransactionResponse } from './transaction';
import { SecurityResponse } from './security';
import { Exchange, SecurityType } from '@/types/enums';

export interface IPortfolioResponse {
  bondDeals: IDealResponse[];
  bondPositions: IPositionResponse[];
  cash: number;
  cashouts: ITransactionResponse[];
  cashoutsSum: number;
  compound: boolean;
  deals: IDealResponse[];
  deposits: ITransactionResponse[];
  depositsSum: number;
  id: number;
  name: string;
  positions: IPositionResponse[];
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
  exchange: Exchange;
  securityType: SecurityType;
  // security: SecurityResponse;
  targetPrice: number | null;
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
