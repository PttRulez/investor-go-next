import { ICouponResponse } from '@/types/apis/go-api';
import { CreateCouponData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorCoupon {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async createCoupon(data: CreateCouponData): Promise<ICouponResponse> {
    const res = await this.api.post<ICouponResponse>('/coupon', data);
    return res.data;
  }

  async deleteCoupon(id: number): Promise<ICouponResponse> {
    const res = await this.api.delete<ICouponResponse>(`/coupon/${id}`);
    return res.data;
  }
}
