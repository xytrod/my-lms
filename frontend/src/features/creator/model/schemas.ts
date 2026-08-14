import { z } from 'zod'

export const courseFormSchema = z.object({
  title: z.string().trim().min(20, 'Минимум 20 символов').max(60, 'Максимум 60 символов'),
  description: z.string().trim().min(20, 'Минимум 20 символов').max(200, 'Максимум 200 символов'),
  level: z.enum(['beginner', 'intermediate', 'advanced']),
})
export type CourseFormValues = z.infer<typeof courseFormSchema>

export const lessonFormSchema = z.object({
  title: z.string().trim().min(10, 'Минимум 10 символов').max(25, 'Максимум 25 символов'),
  content: z.string().trim().min(1, 'Добавьте содержание урока'),
  position: z.number().int('Введите целое число').min(1, 'Минимальная позиция — 1'),
  is_preview: z.boolean(),
})
export type LessonFormValues = z.infer<typeof lessonFormSchema>
