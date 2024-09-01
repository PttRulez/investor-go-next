import { z } from 'zod';

export const UpdatePositionSchema = z.object({
  comment: z.string().nullable(),
  targetPrice: z.number().nullable(),
});

export type UpdatePositionData = z.infer<typeof UpdatePositionSchema>;
