import { DealType, Exchange, SecurityType } from '@/types/enums';
import { z } from 'zod';

export const CreateDealSchema = z.object({
  amount: z
    .number({
      errorMap: _ => ({
        message: 'Введите кол-во бумаг',
      }),
    })
    .int()
    .positive(),
  date: z.date(),
  exchange: z.nativeEnum(Exchange),
  portfolioId: z.number(),
  price: z.number({
    errorMap: _ => ({
      message: 'Введите стоимость сделки',
    }),
  }),
  securityType: z.nativeEnum(SecurityType),
  ticker: z
    .string({
      errorMap: _ => ({
        message: 'Выберите инструмент',
      }),
    })
    .trim()
    .min(1, 'Выберите инструмент'),
  type: z.nativeEnum(DealType),
});

export const UpdateDealSchema = CreateDealSchema.partial().extend({
  id: z.number(),
});

export type CreateDealData = z.infer<typeof CreateDealSchema>;
export type UpdateDealData = z.infer<typeof UpdateDealSchema>;
