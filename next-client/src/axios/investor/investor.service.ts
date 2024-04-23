import { getSession } from 'next-auth/react';
import investorAxiosInstance from './config';
import { InvestorExpert } from './domains/expert';
import {
  InvestorAuth,
  InvestorDeal,
  InvestorMoexBond,
  InvestorMoexShare,
  InvestorOpinion,
  InvestorPortfolio,
  InvestorTransaction,
} from './domains/index';
import { InvestorPosition } from './domains/position';

class InvestorService {
  private static instance: InvestorService;
  private readonly investorApi = investorAxiosInstance;

  public auth: InvestorAuth;
  public deal: InvestorDeal;
  public expert: InvestorExpert;
  public opinion: InvestorOpinion;
  public portfolio: InvestorPortfolio;
  public position: InvestorPosition;
  public moexBond: InvestorMoexBond;
  public moexShare: InvestorMoexShare;
  public transaction: InvestorTransaction;

  constructor() {
    this.auth = new InvestorAuth(this.investorApi);
    this.deal = new InvestorDeal(this.investorApi);
    this.expert = new InvestorExpert(this.investorApi);
    this.moexBond = new InvestorMoexBond(this.investorApi);
    this.moexShare = new InvestorMoexShare(this.investorApi);
    this.opinion = new InvestorOpinion(this.investorApi);
    this.portfolio = new InvestorPortfolio(this.investorApi);
    this.position = new InvestorPosition(this.investorApi);
    this.transaction = new InvestorTransaction(this.investorApi);
  }

  // Singleton
  public static get(): InvestorService {
    if (!InvestorService.instance) {
      InvestorService.instance = new InvestorService();
    }
    return InvestorService.instance;
  }

  public setToken(token: string): void {
    this.investorApi.interceptors.request.use(config => {
      config.headers['Authorization'] = `Bearer ${token}`;
      return config;
    });
  }
}

const investorService = InvestorService.get();

export default investorService;
