import type { CourseLevel } from '../../courses/api/types'

export interface CreateCourseRequest { title: string; description: string; level: CourseLevel }
export type UpdateCourseRequest = Partial<CreateCourseRequest>
export interface CreateLessonRequest { title: string; content: string; position: number; is_preview: boolean }
export type UpdateLessonRequest = Partial<CreateLessonRequest>

