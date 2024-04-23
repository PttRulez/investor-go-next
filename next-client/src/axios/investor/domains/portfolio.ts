import { IPortfolioResponse } from '@/types/apis/go-api';
import { CreatePortfolioData, UpdatePortfolioData } from '@/validation';
import { AxiosInstance } from 'axios';

export class InvestorPortfolio {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  allPortfolios(): Promise<IPortfolioResponse[]> {
    return this.api
      .get('/portfolio')
      .then(res => {
        return res.data ?? [];
      })
      .catch(err => {
        console.log('[investorApi.allPortfolios ERR]:', err);
      });
  }

  createPortfolio(data: CreatePortfolioData): Promise<IPortfolioResponse> {
    return this.api
      .post('/portfolio', data)
      .then(res => {
        return res.data ?? null;
      })
      .catch(err => console.log('[investorApi.getPortfolio ERR]:', err));
  }

  deletePortfolio(id: number) {
    return this.api
      .delete(`/portfolio/${id}`)
      .then(res => {
        return res.data ?? null;
      })
      .catch(err => console.log('[investorApi.getPortfolio ERR]:', err));
  }

  getPortfolio(id: string): Promise<IPortfolioResponse> {
    return this.api
      .get(`/portfolio/${id}`)
      .then(res => {
        return res.data ?? null;
      })
      .catch(err => console.log('[investorApi.getPortfolio ERR]:', err));
  }

  updatePortfolio(data: UpdatePortfolioData) {
    return this.api.put(`/portfolio`, data);
  }
}
