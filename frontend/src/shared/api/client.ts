import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { useAuthStore } from '../../features/auth/model/authStore'
import type { AuthResponseDto } from '../../features/auth/api/types'
import { env } from '../config/env'

interface RetryableRequest extends InternalAxiosRequestConfig {
  _retry?: boolean
}

export const apiClient = axios.create({
  baseURL: env.apiUrl,
  headers: { 'Content-Type': 'application/json' },
})

const refreshClient = axios.create({
  baseURL: env.apiUrl,
  headers: { 'Content-Type': 'application/json' },
})

let refreshPromise: Promise<string> | null = null

apiClient.interceptors.request.use((config) => {
  const accessToken = useAuthStore.getState().tokens?.accessToken
  if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const request = error.config as RetryableRequest | undefined
    const isAuthRequest = request?.url?.startsWith('/auth/login') ||
      request?.url?.startsWith('/auth/register') ||
      request?.url?.startsWith('/auth/refresh')

    if (error.response?.status !== 401 || !request || request._retry || isAuthRequest) {
      return Promise.reject(error)
    }

    const refreshToken = useAuthStore.getState().tokens?.refreshToken
    if (!refreshToken) {
      useAuthStore.getState().clearSession()
      return Promise.reject(error)
    }

    request._retry = true
    refreshPromise ??= refreshClient
      .post<AuthResponseDto>('/auth/refresh', { refresh_token: refreshToken })
      .then(({ data }) => {
        useAuthStore.getState().setSession(
          { userId: data.user_id, role: data.role },
          { accessToken: data.access_token, refreshToken: data.refresh_token },
        )
        return data.access_token
      })
      .catch((refreshError: unknown) => {
        useAuthStore.getState().clearSession()
        throw refreshError
      })
      .finally(() => { refreshPromise = null })

    const accessToken = await refreshPromise
    request.headers.Authorization = `Bearer ${accessToken}`
    return apiClient(request)
  },
)
