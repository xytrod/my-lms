import { apiClient } from '../../../shared/api/client'
import type { EnrollmentDto } from './types'

export interface EnrollmentListParams {
  limit?: number
  offset?: number
}

export async function getMyEnrollments(params: EnrollmentListParams = {}): Promise<EnrollmentDto[]> {
  const { data } = await apiClient.get<EnrollmentDto[]>('/enrollments/my', { params })
  return data
}

export async function enrollInCourse(courseId: string): Promise<EnrollmentDto> {
  const { data } = await apiClient.post<EnrollmentDto>(`/enrollments/${courseId}`)
  return data
}

