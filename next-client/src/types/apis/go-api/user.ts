import { Role } from '@/types/enums';

export interface IUserResponse {
  id: number;
  name: string;
  email: string;
  role: Role;
}
