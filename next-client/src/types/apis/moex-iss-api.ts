import { MoexBoard, MoexSecurityType } from '../enums';
import { MoexSecurityGroup } from './go-api';

export interface IMoexISSSearchResults {
  securities: {
    columns: [
      'id', // 0
      'ticker', // 1
      'shortname', // 2
      'regnumber', // 3
      'name', // 4
      'isin', // 5
      'is_traded', // 6
      'emitent_id', // 7
      'emitent_title', // 8
      'emitent_inn', // 9
      'emitent_okpo', // 10
      'gosreg', // 11
      'type', // 12
      'group', // 13
      'primary_boardid', // 14
      'marketprice_boardid', // 15
    ];
    data: [
      number,
      string,
      string,
      number | null,
      string,
      string | null,
      0 | 1,
      number,
      string | null,
      string | null,
      string | null,
      /** gosreg */
      string | null,
      MoexSecurityType,
      MoexSecurityGroup,
      MoexBoard,
      MoexBoard,
    ][];
  };
}

export interface IMoexISSPricesHistory {
  history: {
    columns: [
      'SHORTNAME', // 0
      'OPEN', // 1
      'HIGH', // 2
      'LOW', // 3
      'CLOSE', // 4
      'TRADEDATE', // 5
    ];
    data: [
      string, // SHORTNAME
      number, // OPEN
      number, // HIGH
      number, // LOW
      number, // CLOSE
      string, // TRADEDATE
    ][];
  };
}
