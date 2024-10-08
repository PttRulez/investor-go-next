import { IDividendResponse } from '@/types/apis/go-api';
import { CreateDividendData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorDividend {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createDividend(data: CreateDividendData): Promise<IDividendResponse> {
    const res = await this.api.post<IDividendResponse>('/dividend', data);
    return res.data;
  }

  async deleteDividend(id: number): Promise<IDividendResponse> {
    const res = await this.api.delete<IDividendResponse>(`/dividend/${id}`);
    return res.data;
  }
}
