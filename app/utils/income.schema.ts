import { z } from 'zod'

export const IncomeSchema = z.object({
  vl: z.number().min(1, 'Salary must be at least 0'),
  type: z.string().min(1, 'Type is required'),
  day: z.number().optional(),
  month: z.string().min(1, 'Month is required'),
})
