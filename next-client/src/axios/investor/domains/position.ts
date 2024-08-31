import { IPositionResponse } from '@/types/apis/go-api';
import { UpdatePositionData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorPosition {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  getAllUserPosition(): Promise<IPositionResponse[]> {
    return this.api
      .get('/position')
      .then(res => {
        return res.data;
      })
      .catch(err => {
        console.log('[investorApi.position.getAllUserPosition ERR]:', err);
      });
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
