// import {StockTypes} from "@/constants/enums";

import { MoexSecurityType } from '@/types/enums';

export const moexStockTypesRU: Record<MoexSecurityType, string> = {
  common_share: 'акция обыкновенная',
  corporate_bond: 'корпоративная облигация',
  depositary_receipt: 'депозитарная расписка',
  exchange_bond: 'облигация',
  exchange_ppif: 'биржевой ПИФ',
  futures: 'фьючерс',
  ofz_bond: 'ОФЗ',
  preferred_share: 'акция привелигированная',
  public_ppif: 'публичный ПИФ',
  stock_index: 'индекс',
  stock_index_if: 'iNAV облигаций',
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
