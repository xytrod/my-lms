import { apiClient } from '../../../shared/api/client'
import type { CourseDto, LessonDto } from '../../courses/api/types'
import type { CreateCourseRequest, CreateLessonRequest, UpdateCourseRequest, UpdateLessonRequest } from './types'

export async function getCreatedCourses(): Promise<CourseDto[]> {
  const { data } = await apiClient.get<CourseDto[]>('/courses/my', { params: { limit: 100, offset: 0 } })
  return data
}
export async function getCreatedCourse(id: string): Promise<CourseDto> { return (await apiClient.get<CourseDto>(`/courses/my/${id}`)).data }
export async function getCreatedCourseLessons(courseId: string): Promise<LessonDto[]> { return (await apiClient.get<LessonDto[]>(`/courses/my/${courseId}/lessons`)).data }
export async function createCourse(body: CreateCourseRequest): Promise<CourseDto> { return (await apiClient.post<CourseDto>('/courses', body)).data }
export async function updateCourse(id: string, body: UpdateCourseRequest): Promise<CourseDto> { return (await apiClient.patch<CourseDto>(`/courses/${id}`, body)).data }
export async function publishCourse(id: string): Promise<CourseDto> { return (await apiClient.post<CourseDto>(`/courses/${id}/publish`)).data }
export async function archiveCourse(id: string): Promise<CourseDto> { return (await apiClient.post<CourseDto>(`/courses/${id}/archive`)).data }
export async function createLesson(courseId: string, body: CreateLessonRequest): Promise<LessonDto> { return (await apiClient.post<LessonDto>(`/courses/${courseId}/lessons`, body)).data }
export async function updateLesson(id: string, body: UpdateLessonRequest): Promise<LessonDto> { return (await apiClient.patch<LessonDto>(`/courses/lessons/${id}`, body)).data }
export async function deleteLesson(id: string): Promise<void> { await apiClient.delete(`/courses/lessons/${id}`) }
