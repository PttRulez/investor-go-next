// import {StockTypes} from "@/constants/enums";

import { MoexSecurityType } from '@/types/enums';

export const moexStockTypesRU: Record<MoexSecurityType, string> = {
  stock_index_if: 'iNAV облигаций',
  common_share: 'акция обыкновенная',
  preferred_share: 'акция привелигированная',
  exchange_bond: 'облигация',
  corporate_bond: 'корпоративная облигация',
  ofz_bond: 'ОФЗ',
  futures: 'фьючерс',
  public_ppif: 'публичный ПИФ',
  exchange_ppif: 'биржевой ПИФ',
  stock_index: 'индекс',
};

// export const moexStockTypeToGeneralType = {
//   common_share: StockTypes.Share,
//   preferred_share: StockTypes.Share,
//   exchange_bond: StockTypes.Bond,
//   corporate_bond: StockTypes.Bond,
//   ofz_bond: StockTypes.Bond,
// }

// export const moexStockGroups = {
//   stock_bonds: 'облигации',
//   stock_shares: 'акция',
//   stock_index: 'индекс',
//   stock_ppif: 'ПИФ',
//   futures_forts: 'фьючерс'
// }

// export const moexMarkets = {
//   stock_bonds: 'bonds',
//   stock_shares: 'shares',
//   stock_index: 'index',
// }
