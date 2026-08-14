export type CourseLevel = 'beginner' | 'intermediate' | 'advanced'

export interface CourseDto {
  id: string
  teacher_id: string
  title: string
  description: string
  level: CourseLevel
  state: string
  created_at: string
  updated_at: string
}

export interface LessonDto {
  id: string
  course_id: string
  title: string
  content: string
  position: number
  is_preview: boolean
  CreatedAt: string
  UpdatedAt: string
}
