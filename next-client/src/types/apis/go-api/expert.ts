import { IOpinionResponse } from './opinion';

export interface IExpertResponse {
	id: number;
	avatarUrl: string | null;
	name: string;
	opinions?: Array<IOpinionResponse>;
}