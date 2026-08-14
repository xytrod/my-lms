import { z } from 'zod'

export const loginSchema = z.object({
  email: z.email('Введите корректный email').max(255),
  password: z.string().min(8, 'Минимум 8 символов').max(42, 'Максимум 42 символа'),
})

export const registerSchema = z.object({
  email: z.email('Введите корректный email').max(255),
  password: z.string().min(8, 'Минимум 8 символов').max(20, 'Максимум 20 символов'),
  username: z.string().min(7, 'Минимум 7 символов').max(25, 'Максимум 25 символов'),
  first_name: z.string().trim().min(2, 'Минимум 2 символа').max(100),
  last_name: z.string().trim().min(2, 'Минимум 2 символа').max(100),
})

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>
