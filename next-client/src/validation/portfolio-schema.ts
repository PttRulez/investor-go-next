import { z } from 'zod';

export const CreatePortfolioSchema = z.object({
  name: z.string(),
  compound: z.boolean(),
});

export const UpdatePortfolioSchema = CreatePortfolioSchema.extend({
  id: z.number(),
});

export type CreatePortfolioData = z.infer<typeof CreatePortfolioSchema>;
export type UpdatePortfolioData = z.infer<typeof UpdatePortfolioSchema>;
