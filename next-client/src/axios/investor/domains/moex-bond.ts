import { IMoexBondResponse } from '@/types/apis/go-api';
import { SecurityType } from '@/types/enums';
import { AxiosInstance } from 'axios';

export class InvestorMoexBond {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  getByTicker(ticker: string): Promise<IMoexBondResponse> {
    return this.api.get(`/moex/bond/${ticker}`).then(res => {
      res.data.securityType = SecurityType.BOND;
      return res.data;
    });
  }
}
