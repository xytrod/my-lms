import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getCurrentUser, login, logout, register } from '../api/authApi'
import type { AuthResponseDto } from '../api/types'
import { useAuthStore } from './authStore'
import type { LoginFormValues, RegisterFormValues } from './schemas'

const currentUserKey = ['auth', 'me'] as const

function applyAuthResponse(data: AuthResponseDto) {
  useAuthStore.getState().setSession(
    { userId: data.user_id, role: data.role },
    { accessToken: data.access_token, refreshToken: data.refresh_token },
  )
}

export function useLogin() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (values: LoginFormValues) => login(values),
    onSuccess: (data) => {
      applyAuthResponse(data)
      queryClient.setQueryData(currentUserKey, { user_id: data.user_id, role: data.role })
    },
  })
}

export function useRegister() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (values: RegisterFormValues) => register(values),
    onSuccess: (data) => {
      applyAuthResponse(data)
      queryClient.setQueryData(currentUserKey, { user_id: data.user_id, role: data.role })
    },
  })
}

export function useCurrentUser() {
  const hasTokens = useAuthStore((state) => Boolean(state.tokens))
  return useQuery({ queryKey: currentUserKey, queryFn: getCurrentUser, enabled: hasTokens })
}

export function useLogout() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async () => {
      const refreshToken = useAuthStore.getState().tokens?.refreshToken
      if (refreshToken) await logout(refreshToken)
    },
    onSettled: () => {
      useAuthStore.getState().clearSession()
      queryClient.clear()
    },
  })
}
