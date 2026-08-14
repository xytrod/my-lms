import axios from 'axios'
import type { ApiErrorResponse } from './types'

export function getApiErrorMessage(error: unknown): string {
  if (!axios.isAxiosError<ApiErrorResponse | string>(error)) return 'Произошла непредвиденная ошибка'

  if (typeof error.response?.data === 'string') {
    switch (error.response.status) {
      case 401: return 'Требуется повторный вход в аккаунт.'
      case 403: return 'У вас нет доступа к этому действию.'
      case 404: return 'Запрошенные данные не найдены.'
      case 409: return 'Действие уже выполнено или конфликтует с текущим состоянием.'
      case 503: return 'Сервис временно недоступен. Попробуйте позже.'
      default: return 'Не удалось выполнить запрос. Попробуйте ещё раз.'
    }
  }
  const responseData = error.response?.data
  const details = typeof responseData === 'object' ? responseData.course_errors : undefined
  if (details?.fields?.length) return details.fields.map((field) => field.message).join('. ')
  if (details?.code === 'course_lessons_locked') return 'Структуру уроков можно изменять только пока курс находится в черновике.'
  if (details?.code === 'course_hasnt_lessons') return 'Для публикации курса добавьте хотя бы один урок.'
  return details?.message || 'Не удалось выполнить запрос. Попробуйте ещё раз.'
}
