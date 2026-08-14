import { apiClient } from '../../../shared/api/client'
import type { CourseDto, LessonDto } from './types'

export interface CourseListParams {
  search?: string
  limit?: number
  offset?: number
}

export async function getCourses(params: CourseListParams = {}): Promise<CourseDto[]> {
  const { data } = await apiClient.get<CourseDto[]>('/courses', { params })
  return data
}

export async function getCourse(courseId: string): Promise<CourseDto> {
  const { data } = await apiClient.get<CourseDto>(`/courses/${courseId}`)
  return data
}

export async function getCourseLessons(courseId: string): Promise<LessonDto[]> {
  const { data } = await apiClient.get<LessonDto[]>(`/courses/${courseId}/lessons`)
  return data
}
