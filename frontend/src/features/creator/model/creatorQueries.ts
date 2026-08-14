import axios from 'axios'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { CourseDto } from '../../courses/api/types'
import { courseKeys } from '../../courses/model/courseQueries'
import { archiveCourse, createCourse, createLesson, deleteLesson, getCreatedCourse, getCreatedCourseLessons, getCreatedCourses, publishCourse, updateCourse, updateLesson } from '../api/creatorApi'

export const creatorKeys = {
  all: ['creator-courses'] as const,
  list: () => [...creatorKeys.all, 'list'] as const,
  detail: (id: string) => [...creatorKeys.all, 'detail', id] as const,
  lessons: (id: string) => [...creatorKeys.detail(id), 'lessons'] as const,
}

export function useCreatedCourses() { return useQuery({ queryKey: creatorKeys.list(), queryFn: getCreatedCourses }) }
export function useCreatedCourse(id: string) { return useQuery({ queryKey: creatorKeys.detail(id), queryFn: () => getCreatedCourse(id), enabled: Boolean(id) }) }
export function useCreatedCourseLessons(id: string) { return useQuery({ queryKey: creatorKeys.lessons(id), queryFn: () => getCreatedCourseLessons(id), enabled: Boolean(id) }) }

function updateCourseCaches(client: ReturnType<typeof useQueryClient>, course: CourseDto) {
  client.setQueryData(creatorKeys.detail(course.id), course)
  client.setQueryData<CourseDto[]>(creatorKeys.list(), (items = []) => items.map((item) => item.id === course.id ? course : item))
}

function useStateMutation(fn: (id: string) => Promise<CourseDto>) {
  const client = useQueryClient()
  return useMutation({ mutationFn: fn, onSuccess: async (course) => {
    updateCourseCaches(client, course)
    await Promise.all([
      client.invalidateQueries({ queryKey: creatorKeys.list() }),
      client.invalidateQueries({ queryKey: creatorKeys.detail(course.id) }),
      client.invalidateQueries({ queryKey: creatorKeys.lessons(course.id) }),
      client.invalidateQueries({ queryKey: courseKeys.all }),
    ])
  } })
}

export function useCreateCourse() { const client = useQueryClient(); return useMutation({ mutationFn: createCourse, onSuccess: (course) => { updateCourseCaches(client, course) } }) }
export function useUpdateCourse() { const client = useQueryClient(); return useMutation({ mutationFn: ({ id, body }: { id: string; body: Parameters<typeof updateCourse>[1] }) => updateCourse(id, body), onSuccess: (course) => {
  updateCourseCaches(client, course)
  if (course.state === 'published') void client.invalidateQueries({ queryKey: courseKeys.all })
} }) }
export function usePublishCourse() { return useStateMutation(publishCourse) }
export function useArchiveCourse() { return useStateMutation(archiveCourse) }

function useLessonMutation<TVariables>(courseId: string, mutationFn: (variables: TVariables) => Promise<unknown>) {
  const client = useQueryClient()
  return useMutation({ mutationFn, onSuccess: () => client.invalidateQueries({ queryKey: creatorKeys.lessons(courseId) }), onError: (error) => {
    if (axios.isAxiosError(error) && error.response?.status === 409) {
      void client.invalidateQueries({ queryKey: creatorKeys.detail(courseId) })
    }
  } })
}
export function useCreateLesson(courseId: string) { return useLessonMutation(courseId, (body: Parameters<typeof createLesson>[1]) => createLesson(courseId, body)) }
export function useUpdateLesson(courseId: string) { return useLessonMutation(courseId, ({ id, body }: { id: string; body: Parameters<typeof updateLesson>[1] }) => updateLesson(id, body)) }
export function useDeleteLesson(courseId: string) { return useLessonMutation(courseId, deleteLesson) }
