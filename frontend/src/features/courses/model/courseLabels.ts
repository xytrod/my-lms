import type { CourseLevel } from '../api/types'

export const courseLevelLabels: Record<CourseLevel, string> = {
  beginner: 'Начальный',
  intermediate: 'Средний',
  advanced: 'Продвинутый',
}

export const courseStateLabels: Record<string, string> = {
  Draft: 'Черновик',
  published: 'Опубликован',
  closed: 'Закрыт',
}

export function getCourseStateLabel(state: string): string {
  return courseStateLabels[state] ?? 'Неизвестное состояние'
}
