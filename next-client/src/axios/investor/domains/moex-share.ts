import { IMoexShareResponse } from '@/types/apis/go-api';
import { SecurityType } from '@/types/enums';
import { AxiosInstance } from 'axios';

export class InvestorMoexShare {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  getByTicker(ticker: string): Promise<IMoexShareResponse> {
    return this.api.get(`/moex/share/${ticker}`).then(res => {
      res.data.securityType = SecurityType.SHARE;
      return res.data;
    });
  }
}
