import { Link, Navigate } from 'react-router-dom'
import { AuthLayout } from '../../features/auth/ui/AuthLayout'
import { LoginForm } from '../../features/auth/ui/LoginForm'
import { useAuthStore } from '../../features/auth/model/authStore'

export function LoginPage() {
  const authenticated = useAuthStore((state) => Boolean(state.tokens))
  if (authenticated) return <Navigate to="/dashboard" replace />

  return (
    <AuthLayout title="С возвращением" subtitle="Войдите, чтобы продолжить обучение">
      <LoginForm />
      <p className="auth-switch">Нет аккаунта? <Link to="/register">Зарегистрироваться</Link></p>
    </AuthLayout>
  )
}
