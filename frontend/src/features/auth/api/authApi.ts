import { apiClient } from '../../../shared/api/client'
import type { AuthResponseDto, CurrentUserDto, LoginRequestDto, RegisterRequestDto } from './types'

export async function login(payload: LoginRequestDto): Promise<AuthResponseDto> {
  const { data } = await apiClient.post<AuthResponseDto>('/auth/login', payload)
  return data
}

export async function register(payload: RegisterRequestDto): Promise<AuthResponseDto> {
  const { data } = await apiClient.post<AuthResponseDto>('/auth/register', payload)
  return data
}

export async function getCurrentUser(): Promise<CurrentUserDto> {
  const { data } = await apiClient.get<CurrentUserDto>('/auth/me')
  return data
}

export async function logout(refreshToken: string): Promise<void> {
  await apiClient.post('/auth/logout', { refresh_token: refreshToken })
}
