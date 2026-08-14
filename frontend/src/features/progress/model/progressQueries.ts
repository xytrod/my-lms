import { useQuery } from '@tanstack/react-query'
import { getCourseProgress } from '../api/progressApi'

export const progressKeys = {
  all: ['progress'] as const,
  course: (courseId: string) => [...progressKeys.all, courseId] as const,
}

export function useCourseProgress(courseId: string, enabled = true) {
  return useQuery({
    queryKey: progressKeys.course(courseId),
    queryFn: () => getCourseProgress(courseId),
    enabled: Boolean(courseId) && enabled,
  })
}
