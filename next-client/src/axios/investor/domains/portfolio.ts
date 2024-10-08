import { IPortfolioResponse } from '@/types/apis/go-api';
import { CreatePortfolioData, UpdatePortfolioData } from '@/validation';
import { CreateExpenseData } from '@/validation/expense-schema';
import { AxiosInstance } from 'axios';

export class InvestorPortfolio {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  // Portfolio
  allPortfolios(): Promise<IPortfolioResponse[]> {
    return this.api.get('/portfolio').then(res => {
      return res.data ?? [];
    });
  }

  createPortfolio(data: CreatePortfolioData): Promise<IPortfolioResponse> {
    return this.api.post('/portfolio', data).then(res => {
      return res.data ?? null;
    });
  }

  deletePortfolio(id: number) {
    return this.api.delete(`/portfolio/${id}`).then(res => {
      return res.data ?? null;
    });
  }

  getPortfolio(id: string): Promise<IPortfolioResponse> {
    return this.api.get(`/portfolio/${id}`).then(res => {
      return res.data ?? null;
    });
  }

  updatePortfolio(data: UpdatePortfolioData) {
    return this.api.put(`/portfolio`, data);
  }

  // expense
  createExpense(data: CreateExpenseData): Promise<void> {
    return this.api.post('/expense', data).then(res => {
      return res.data ?? null;
    });
  }

  deleteExpense(id: number) {
    return this.api.delete(`/expense/${id}`).then(res => {
      return res.data ?? null;
    });
  }
}
