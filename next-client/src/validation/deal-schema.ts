import { DealType } from '@/types/enums';
import { z } from 'zod';

export const CreateMoexShareDealSchema = z.object({
  amount: z
    .number({
      errorMap: _ => ({
        message: 'Введите кол-во бумаг',
      }),
    })
    .int()
    .positive(),
  date: z.date(),
  portfolioId: z.number(),
  price: z.number({
    errorMap: _ => ({
      message: 'Введите стоимость сделки',
    }),
  }),
  ticker: z.number({
    errorMap: _ => ({
      message: 'ВЫберите бумагу',
    }),
  }),
  type: z.nativeEnum(DealType),
});

export const UpdateDealSchema = CreateMoexShareDealSchema.partial().extend({
  id: z.number(),
});

export type CreateMoexShareDealData = z.infer<typeof CreateMoexShareDealSchema>;
export type UpdateDealData = z.infer<typeof UpdateDealSchema>;
