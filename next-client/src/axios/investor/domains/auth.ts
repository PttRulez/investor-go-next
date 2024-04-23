import { IUserResponse } from '@/types/apis/go-api';
import { LoginData, RegisterData } from '@/validation';
import { AxiosInstance, AxiosResponse } from 'axios';

export class InvestorAuth {
  private readonly api: AxiosInstance;

  constructor(api: AxiosInstance) {
    this.api = api;
  }

  async login(dto: LoginData): Promise<AxiosResponse<IUserResponse, any>> {
    return this.api.post('/login', { ...dto });
  }

  register(dto: RegisterData): Promise<AxiosResponse<IUserResponse, any>> {
    return this.api.post('/register', { ...dto });
  }
}
