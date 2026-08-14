import { useQuery } from '@tanstack/react-query'
import { getCourse, getCourseLessons, getCourses, type CourseListParams } from '../api/courseApi'

export const courseKeys = {
  all: ['courses'] as const,
  list: (params: CourseListParams) => [...courseKeys.all, 'list', params] as const,
  detail: (courseId: string) => [...courseKeys.all, 'detail', courseId] as const,
  lessons: (courseId: string) => [...courseKeys.detail(courseId), 'lessons'] as const,
}

export function useCourses(params: CourseListParams = {}) {
  return useQuery({ queryKey: courseKeys.list(params), queryFn: () => getCourses(params) })
}

export function useCourse(courseId: string) {
  return useQuery({ queryKey: courseKeys.detail(courseId), queryFn: () => getCourse(courseId), enabled: Boolean(courseId) })
}

export function useCourseLessons(courseId: string, enabled = true) {
  return useQuery({ queryKey: courseKeys.lessons(courseId), queryFn: () => getCourseLessons(courseId), enabled: Boolean(courseId) && enabled })
}
