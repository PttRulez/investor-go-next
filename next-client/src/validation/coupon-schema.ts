import { Exchange } from '@/types/enums';
import { z } from 'zod';

export const CreateCouponSchema = z.object({
  bondsCount: z.number({
    errorMap: _ => ({
      message: 'Введите количество облигаций',
    }),
  }),
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  paymentPeriod: z.string(),
  portfolioId: z.number(),
  taxPaid: z.number({
    errorMap: _ => ({
      message: 'Введите размер уплаченного налога',
    }),
  }),
  ticker: z.string({
    errorMap: _ => ({
      message: 'Выберите облигацию',
    }),
  }),
  totalPayment: z.number({
    errorMap: _ => ({
      message: 'Введите сумму выплаты, полученную в рублях',
    }),
  }),
});

export const UpdateCouponSchema = CreateCouponSchema.partial().extend({
  id: z.number(),
});

export type CreateCouponData = z.infer<typeof CreateCouponSchema>;
export type UpdateCouponData = z.infer<typeof UpdateCouponSchema>;
