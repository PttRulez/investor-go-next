import { IDealResponse } from '@/types/apis/go-api';
import { CreateMoexShareDealData, UpdateDealData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorDeal {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createDeal(data: CreateMoexShareDealData): Promise<IDealResponse> {
    const res = await this.api.post<IDealResponse>('/deal/moex-share', data);
    return res.data;
  }

  async deleteDeal(id: number): Promise<IDealResponse> {
    const res = await this.api.delete<IDealResponse>(`/deal/${id}`);
    return res.data;
  }

  updateDeal(data: UpdateDealData) {
    return this.api.post(`/deals/moex-share/${data.id}`, data);
  }
}
