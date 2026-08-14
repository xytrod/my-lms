import { apiClient } from '../../../shared/api/client'
import type { ProgressDto } from './types'

export async function getCourseProgress(courseId: string): Promise<ProgressDto> {
  const { data } = await apiClient.get<ProgressDto>(`/enrollments/${courseId}/progress`)
  return data
}

export async function completeLesson(courseId: string, lessonId: string): Promise<void> {
  await apiClient.post(`/enrollments/${courseId}/lessons/${lessonId}/complete`)
}

