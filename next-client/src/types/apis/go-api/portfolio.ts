import { SecurityType } from '@/types/enums';
import { IOpinionResponse } from './opinion';
import { DealType, Exchange } from '@/types/enums';
import { TransactionType } from '@/types/enums';

export interface ICouponResponse {
  bondsCount: number;
  couponAmount: number;
  date: Date;
  exchange: Exchange;
  id: number;
  paymentPeriod: string;
  portfolioId: number;
  ticker: string;
}

export interface IDealResponse {
  amount: number;
  date: Date;
  exchange: Exchange;
  id: number;
  portfolioId: number;
  price: number;
  securityId: number;
  securityType: SecurityType;
  ticker: string;
  type: DealType;
}

export interface IDividendResponse {
  date: Date;
  exchange: Exchange;
  id: number;
  paymentPeriod: string;
  paymentPerShare: number;
  portfolioId: number;
  ticker: string;
  sharesCount: number;
}

export interface IPortfolioResponse {
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
  profitability: number;
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
  opinions: IOpinionResponse[];
  opinionIds: number[];
  portfolioName: string;
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

export interface ITransactionResponse {
  id: number;
  amount: number;
  type: TransactionType;
  date: Date;
}
