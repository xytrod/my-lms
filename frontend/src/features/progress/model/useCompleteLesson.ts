import { useMutation, useQueryClient } from '@tanstack/react-query'
import { enrollmentKeys } from '../../enrollments/model/enrollmentQueries'
import { completeLesson } from '../api/progressApi'
import { progressKeys } from './progressQueries'

interface CompleteLessonVariables {
  courseId: string
  lessonId: string
}

export function useCompleteLesson() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ courseId, lessonId }: CompleteLessonVariables) => completeLesson(courseId, lessonId),
    onSuccess: async (_data, { courseId }) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: progressKeys.course(courseId) }),
        queryClient.invalidateQueries({ queryKey: enrollmentKeys.all }),
      ])
    },
  })
}
