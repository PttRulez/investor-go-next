import { IDealResponse } from './deal';
import { ITransactionResponse } from './transaction';

export interface IPortfolioResponse {
  bondDeals: IDealResponse[];
  bondPositions: IMoexBondPositionResponse[];
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
  sharePositions: IMoexSharePositionResponse[];
  totalCost: number;
}

export interface IPositionResponse {
  amount: number;
  averagePrice: number;
  comment: string | null;
  id: number;
  currentPrice: number;
  currentCost: number;
  shortName: string;
  targetPrice: number | null;
}

// export interface IMoexBondPositionResponse extends IPositionResponse {
//   isin: string;
// }
export type IMoexBondPositionResponse = IPositionResponse & {
  isin: string;
};
export type IMoexSharePositionResponse = IPositionResponse & {
  ticker: string;
};
// export interface IMoexSharePositionResponse extends IPositionResponse {
//   ticker: string;
// }

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
