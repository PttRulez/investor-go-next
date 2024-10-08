import { Exchange } from '@/types/enums';
import { z } from 'zod';

export const CreateDividendSchema = z.object({
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  paymentPeriod: z.string(),
  paymentPerShare: z.number({
    errorMap: _ => ({
      message: 'Введите размер выплаты на одну акцию сделки',
    }),
  }),
  portfolioId: z.number(),
  ticker: z.string({
    errorMap: _ => ({
      message: 'Выберите бумагу',
    }),
  }),
  sharesCount: z.number({
    errorMap: _ => ({
      message: 'Введите количество акций',
    }),
  }),
});

export const UpdateDividendSchema = CreateDividendSchema.partial().extend({
  id: z.number(),
});

export type CreateDividendData = z.infer<typeof CreateDividendSchema>;
export type UpdateDividendData = z.infer<typeof UpdateDividendSchema>;
