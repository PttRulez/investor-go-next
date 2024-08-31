import { MoexBoard, MoexEngine, MoexMarket, SecurityType } from '@/types/enums';
import { BaseSecurityResponse } from './security';

export enum MoexSecurityGroup {
  stock_index = 'stock_index', // Индексы
  stock_shares = 'stock_shares', // Акции
  stock_bonds = 'stock_bonds', // Облигации
  currency_selt = 'currency_selt', // Валюта
  futures_forts = 'futures_forts', // Фьючерсы
  futures_options = 'futures_options', // Опционы
  stock_dr = 'stock_dr', // Депозитарные расписки
  stock_foreign_shares = 'stock_foreign_shares', // Иностранные ц.б.
  stock_eurobond = 'stock_eurobond', // Еврооблигации
  stock_ppif = 'stock_ppif', // Паи ПИФов
  stock_etf = 'stock_etf', // Биржевые фонды
  currency_metal = 'currency_metal', // Драгоценные металлы
  stock_qnv = 'stock_qnv', // Квал. инвесторы
  stock_gcc = 'stock_gcc', // Клиринговые сертификаты участия
  stock_deposit = 'stock_deposit', // Депозиты с ЦК
  currency_futures = 'currency_futures', // Валютный фьючерс
  currency_indices = 'currency_indices', // Валютные фиксинги
  stock_mortgage = 'stock_mortgage', // Ипотечный сертификат
}

export interface IMoexSecurtiyResponse extends BaseSecurityResponse {
  board: MoexBoard;
  engine: MoexEngine;
  id: number;
  market: MoexMarket;
  name: string;
  shortName: string;
  securityType: SecurityType;
  ticker: string;
}

export type IMoexShareResponse = IMoexSecurtiyResponse;
export type IMoexBondResponse = IMoexSecurtiyResponse;
