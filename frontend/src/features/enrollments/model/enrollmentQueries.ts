import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '../../auth/model/authStore'
import { enrollInCourse, getMyEnrollments, type EnrollmentListParams } from '../api/enrollmentApi'
import type { EnrollmentDto } from '../api/types'
import { progressKeys } from '../../progress/model/progressQueries'
import axios from 'axios'

export const enrollmentKeys = {
  all: ['enrollments'] as const,
  mine: (params: EnrollmentListParams) => [...enrollmentKeys.all, 'mine', params] as const,
}

export const allMyEnrollmentsParams = { limit: 100, offset: 0 } as const

export function useMyEnrollments(enabled = true) {
  const authenticated = useAuthStore((state) => Boolean(state.tokens))
  return useQuery({
    queryKey: enrollmentKeys.mine(allMyEnrollmentsParams),
    queryFn: () => getMyEnrollments(allMyEnrollmentsParams),
    enabled: authenticated && enabled,
  })
}

export function useEnrollInCourse() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: enrollInCourse,
    onSuccess: (enrollment) => {
      queryClient.setQueryData<EnrollmentDto[]>(enrollmentKeys.mine(allMyEnrollmentsParams), (current = []) => {
        if (current.some((item) => item.CourseID === enrollment.CourseID)) return current
        return [...current, enrollment]
      })
      void queryClient.invalidateQueries({ queryKey: enrollmentKeys.all })
      void queryClient.invalidateQueries({ queryKey: progressKeys.course(enrollment.CourseID) })
    },
    onError: (error) => {
      // A concurrent/remounted request may lose an already-enrolled race; refetch server truth.
      if (axios.isAxiosError(error) && error.response?.status === 409) {
        void queryClient.invalidateQueries({ queryKey: enrollmentKeys.all })
      }
    },
  })
}
