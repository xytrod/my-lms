import type { UserRole } from '../../../shared/types/auth'

export interface LoginRequestDto {
  email: string
  password: string
}

export interface RegisterRequestDto extends LoginRequestDto {
  username: string
  first_name: string
  last_name: string
}

export interface AuthResponseDto {
  user_id: string
  role: UserRole
  access_token: string
  refresh_token: string
}

export interface CurrentUserDto {
  user_id: string
  role: UserRole
}
