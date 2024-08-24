import { DealType, Exchange, SecurityType } from '@/types/enums';
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
  comission: z.number(),
  date: z.date(),
  exchange: z.nativeEnum(Exchange),
  portfolioId: z.number(),
  price: z.number({
    errorMap: _ => ({
      message: 'Введите стоимость сделки',
    }),
  }),
  securityId: z.number(),
  securityType: z.nativeEnum(SecurityType),
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
