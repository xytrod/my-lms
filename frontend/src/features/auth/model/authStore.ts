import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { AuthUser, TokenPair } from '../../../shared/types/auth'

interface AuthState {
  user: AuthUser | null
  tokens: TokenPair | null
  setSession: (user: AuthUser, tokens: TokenPair) => void
  updateTokens: (tokens: TokenPair) => void
  clearSession: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      tokens: null,
      setSession: (user, tokens) => set({ user, tokens }),
      updateTokens: (tokens) => set({ tokens }),
      clearSession: () => set({ user: null, tokens: null }),
    }),
    {
      name: 'lms-auth',
      partialize: ({ user, tokens }) => ({ user, tokens }),
    },
  ),
)
