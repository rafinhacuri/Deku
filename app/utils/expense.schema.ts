import { z } from 'zod'

export const ExpenseSchema = z.object({
  vl: z.number().min(1, 'Value must be non-negative'),
  type: z.string().min(1, 'Type is required'),
  description: z.string().optional(),
  day: z.number().min(1).max(31).optional(),
  paymentMethod: z.string().min(1, 'Payment method is required'),
  month: z.string().min(1, 'Month is required'),
})
