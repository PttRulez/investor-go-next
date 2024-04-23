import { UpdatePositionData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorPosition {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  update(id: number, data: UpdatePositionData) {
    return this.api
      .patch(`/position/${id}`, data)
      .then(res => {
        return res.data;
      })
      .catch(err => {
        console.log('[investorApi.position.update ERR]:', err);
      });
  }
}
