import { IOpinionResponse } from '@/types/apis/go-api';
import { CreateOpinionData, OpinionFilters } from '@/validation';
import { AxiosInstance, AxiosResponse } from 'axios';

export class InvestorOpinion {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async attachToPosition(i: AttachInfo): Promise<AxiosResponse<any, any>> {
    return this.api.get(
      `/opinion/${i.opinionId}/attach-position/${i.positionId}`,
    );
  }

  async createOpinion(data: CreateOpinionData): Promise<IOpinionResponse> {
    const res = await this.api.post<IOpinionResponse>('/opinion', data);
    return res.data;
  }

  async getOpinionsList(f: OpinionFilters): Promise<IOpinionResponse[]> {
    // let key: keyof OpinionFilters;
    // let params: OpinionFilters = {} as OpinionFilters;
    // for (key in f) {
    //   if (f[key] != undefined) {
    //     params[key] = f[key];
    //   }
    // }
    const res = await this.api.get<IOpinionResponse[]>('/opinion/list', {
      params: f,
    });
    return res.data;
  }
}

export type AttachInfo = {
  opinionId: number;
  positionId: number;
};
