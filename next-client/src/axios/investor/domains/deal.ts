import { IDealResponse } from '@/types/apis/go-api';
import { CreateDealData, UpdateDealData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorDeal {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createDeal(data: CreateDealData): Promise<IDealResponse> {
    const res = await this.api.post<IDealResponse>('/deal', data);
    return res.data;
  }

  async deleteDeal(id: number): Promise<IDealResponse> {
    const res = await this.api.delete<IDealResponse>(`/deal/${id}`);
    return res.data;
  }
}

// 2024-08-25T09:00:26.130Z
