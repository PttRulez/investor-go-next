import { Exchange } from '@/types/enums';
import { z } from 'zod';

export const CreateCouponSchema = z.object({
  bondsCount: z.number({
    errorMap: _ => ({
      message: 'Введите количество облигаций',
    }),
  }),
  couponAmount: z.number({
    errorMap: _ => ({
      message: 'Введите размер купона в рублях',
    }),
  }),
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  paymentPeriod: z.string(),
  portfolioId: z.number(),
  ticker: z.string({
    errorMap: _ => ({
      message: 'Выберите облигацию',
    }),
  }),
});

export const UpdateCouponSchema = CreateCouponSchema.partial().extend({
  id: z.number(),
});

export type CreateCouponData = z.infer<typeof CreateCouponSchema>;
export type UpdateCouponData = z.infer<typeof UpdateCouponSchema>;
