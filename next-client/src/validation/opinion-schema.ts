import { Exchange, OpinionType, SecurityType } from '@/types/enums';
import { z } from 'zod';

// Attach Opinion
export const AttachOpinionSchema = z.object({
  opinionId: z.number(),
  positionId: z.number(),
});
export type AttachOpinionData = z.infer<typeof AttachOpinionSchema>;

// Create Opinion
export const CreateOpinionSchema = z.object({
  date: z.string(),
  exchange: z.nativeEnum(Exchange),
  expertId: z.number(),
  text: z.string(),
  securityType: z.nativeEnum(SecurityType),
  securityId: z.number(),
  sourceLink: z.string().nullable(),
  targetPrice: z.number().nullable(),
  ticker: z.string(),
  type: z.nativeEnum(OpinionType),
});
export type CreateOpinionData = z.infer<typeof CreateOpinionSchema>;

// Update Opinion
export const UpdateOpinionSchema = CreateOpinionSchema.partial().extend({
  id: z.number(),
});
export type UpdateOpinionData = z.infer<typeof UpdateOpinionSchema>;

// Get Opinion List
export const OpinionFiltersSchema = z.object({
  exchange: z.nativeEnum(Exchange).optional(),
  expertId: z.number().optional(),
  securityType: z.nativeEnum(SecurityType).optional(),
  securityId: z.number().optional(),
  type: z.nativeEnum(OpinionType).optional(),
});

export type OpinionFilters = z.infer<typeof OpinionFiltersSchema>;
