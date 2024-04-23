import { z } from "zod";

export const UpdatePositionSchema = z.object({
  comment: z.string().nullable(),
  // opinions: z.array(z.number()).nullable(),
  targetPrice: z.number().nullable(),
});

export type UpdatePositionData = z.infer<typeof UpdatePositionSchema>;
