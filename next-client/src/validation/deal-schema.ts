import { DealType, Exchange, SecurityType } from '@/types/enums';
import { z } from 'zod';

const DealSchema = z.object({
  amount: z
    .number({
      errorMap: _ => ({
        message: 'Введите кол-во бумаг',
      }),
    })
    .int()
    .positive(),
  comission: z.number({
    errorMap: _ => ({
      message: 'Введите комиссию',
    }),
  }),
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  nkd: z.number().optional(),
  portfolioId: z.number(),
  price: z.number({
    errorMap: _ => ({
      message: 'Введите стоимость сделки',
    }),
  }),
  shortName: z.string({
    errorMap: _ => ({
      message: 'Выберите бумагу',
    }),
  }),
  ticker: z.string({
    errorMap: _ => ({
      message: 'Выберите бумагу',
    }),
  }),
  securityType: z.nativeEnum(SecurityType),
  type: z.nativeEnum(DealType),
});

export const CreateDealSchema = DealSchema.refine(
  d => {
    if (d.securityType === SecurityType.BOND && d.nkd == undefined)
      return false;
    return true;
  },
  {
    message: 'Введите НКД',
    path: ['nkd'],
  },
);

export const UpdateDealSchema = DealSchema.partial().extend({
  id: z.number(),
});

export type CreateDealData = z.infer<typeof CreateDealSchema>;
export type UpdateDealData = z.infer<typeof UpdateDealSchema>;
