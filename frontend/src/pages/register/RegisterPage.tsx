import { Link, Navigate } from 'react-router-dom'
import { AuthLayout } from '../../features/auth/ui/AuthLayout'
import { RegisterForm } from '../../features/auth/ui/RegisterForm'
import { useAuthStore } from '../../features/auth/model/authStore'

export function RegisterPage() {
  const authenticated = useAuthStore((state) => Boolean(state.tokens))
  if (authenticated) return <Navigate to="/dashboard" replace />

  return (
    <AuthLayout title="Начните учиться" subtitle="Создайте бесплатный аккаунт студента">
      <RegisterForm />
      <p className="auth-switch">Уже есть аккаунт? <Link to="/login">Войти</Link></p>
    </AuthLayout>
  )
}
