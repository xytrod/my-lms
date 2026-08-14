import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuthStore } from '../model/authStore'
import { useCurrentUser } from '../model/useAuth'
import { Spinner } from '../../../shared/ui/Spinner'

export function ProtectedRoute() {
  const location = useLocation()
  const tokens = useAuthStore((state) => state.tokens)
  const clearSession = useAuthStore((state) => state.clearSession)
  const currentUser = useCurrentUser()

  if (!tokens) return <Navigate to="/login" replace state={{ from: location }} />
  if (currentUser.isPending) return <main className="center-state"><Spinner /><p>Проверяем сессию…</p></main>
  if (currentUser.isError) {
    clearSession()
    return <Navigate to="/login" replace state={{ from: location }} />
  }
  return <Outlet />
}
