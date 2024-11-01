import { IPortfolioResponse } from '@/types/apis/go-api/portfolio';

type Props = {
  portfolio: IPortfolioResponse;
};

type HistoryRow = {
  // deals
  amount: number;
  comission: number;
  price: number;
  nkd?: number;

  // common
  date: Date;
};

export default function HistoryTable({ portfolio }: Props) {}
