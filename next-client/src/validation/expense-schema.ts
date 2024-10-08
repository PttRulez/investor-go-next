import { z } from 'zod';

export const CreateExpenseSchema = z.object({
  amount: z.number({
    errorMap: _ => ({
      message: 'Введите размер затраты в рублях',
    }),
  }),
  date: z.string(),
  description: z.string(),
  portfolioId: z.number(),
});

export const UpdateExpenseSchema = CreateExpenseSchema.partial().extend({
  id: z.number(),
});

export type CreateExpenseData = z.infer<typeof CreateExpenseSchema>;
export type UpdateExpenseData = z.infer<typeof UpdateExpenseSchema>;
