import { Exchange } from '@/types/enums';
import { z } from 'zod';

export const CreateDividendSchema = z.object({
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  paymentPeriod: z.string(),
  portfolioId: z.number(),
  sharesCount: z.number({
    errorMap: _ => ({
      message: 'Введите количество акций',
    }),
  }),
  taxPaid: z.number({
    errorMap: _ => ({
      message: 'Введите размер уплаченного налога',
    }),
  }),
  ticker: z.string({
    errorMap: _ => ({
      message: 'Выберите бумагу',
    }),
  }),
  totalPayment: z.number({
    errorMap: _ => ({
      message: 'Введите сумму выплаты, полученную в рублях',
    }),
  }),
});

export const UpdateDividendSchema = CreateDividendSchema.partial().extend({
  id: z.number(),
});

export type CreateDividendData = z.infer<typeof CreateDividendSchema>;
export type UpdateDividendData = z.infer<typeof UpdateDividendSchema>;
