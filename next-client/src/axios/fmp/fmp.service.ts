import { fmpApiV3 } from './fmp';
import {
  historicalDaily,
  requestPeriodType,
  searchStock,
} from '@/types/fmpTypes/fmp.type';
import { AxiosResponse } from 'axios';

export const fmpService = {
  search(searchParam: string): Promise<AxiosResponse<searchStock[]> | void> {
    return fmpApiV3
      .get(`/search?query=${searchParam}&limit=30`)
      .then(res => {
        return res.data;
      })
      .catch(err => console.log('[fmpService.search err]', err));
  },

  historicalDailyByTicker(ticker: string | undefined):
    | Promise<AxiosResponse<{
        symbol: string;
        historical: historicalDaily[];
      }> | void>
    | Promise<Awaited<{ symbol: string; historical: any[] }>> {
    if (!ticker) {
      return Promise.resolve({
        symbol: 'No ticker',
        historical: [],
      });
    }
    return fmpApiV3
      .get(`/historical-price-full/${ticker}`)
      .then(res => {
        return res.data;
      })
      .catch(err =>
        console.log('[fmpService.historicalPriceByTicker err]', err),
      );
  },

  keyMetricsTTM(ticker: string) {
    return fmpApiV3.get(`key-metrics-ttm/${ticker}`);
  },

  incomeStatementAnnual(ticker: string) {
    return fmpApiV3.get(`income-statement/${ticker}?limit=20`);
  },
  balanceSheetStatementAnnual(ticker: string) {
    return fmpApiV3.get(`balance-sheet-statement/${ticker}?limit=20`);
  },
  balanceSheet(
    ticker: string | undefined,
    period: requestPeriodType = 'year',
  ) {},
};
