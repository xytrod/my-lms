export type UserRole = 'student' | 'teacher' | 'admin'

export interface AuthUser {
  userId: string
  role: UserRole
}

export interface TokenPair {
  accessToken: string
  refreshToken: string
}
