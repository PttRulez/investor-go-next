import { IOpinionResponse } from '@/types/apis/go-api';
import { CreateOpinionData, OpinionFilters } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorOpinion {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createOpinion(data: CreateOpinionData): Promise<IOpinionResponse> {
    const res = await this.api.post<IOpinionResponse>('/opinion', data);
    return res.data;
  }

  async getOpinionsList(filters: OpinionFilters): Promise<IOpinionResponse[]> {
    const res = await this.api.post<IOpinionResponse[]>(
      '/opinion/list',
      filters,
    );
    return res.data;
  }
}
