import { IExpertResponse } from '@/types/apis/go-api';
import { CreateExpertData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorExpert {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createExpert(data: CreateExpertData): Promise<IExpertResponse> {
    const res = await this.api.post<IExpertResponse>('/expert', data);
    return res.data;
  }

  async getExpertsList(): Promise<IExpertResponse[]> {
    const res = await this.api.get<IExpertResponse[]>('/expert');
    return res.data;
  }
}
