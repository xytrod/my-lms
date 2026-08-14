import { Link, Navigate } from 'react-router-dom'
import { useAuthStore } from '../../features/auth/model/authStore'

const roleLabels = { student: 'Студент', teacher: 'Преподаватель', admin: 'Администратор' } as const

export function DashboardPage() {
  const user = useAuthStore((state) => state.user)
  if (!user) return <Navigate to="/login" replace />
  return <main className="page-container dashboard"><section className="welcome-card">
    <span className="role-badge">{roleLabels[user.role]}</span><h1>Добро пожаловать в LMS</h1>
    <p>{user.role === 'student' ? 'Выбирайте новые программы или возвращайтесь к курсам, на которые уже записаны.' : 'В каталоге доступны опубликованные программы платформы.'}</p>
    <div className="dashboard-actions"><Link className="button dashboard-action" to="/courses">Открыть каталог</Link>{user.role === 'student' && <Link className="button button-secondary dashboard-action" to="/my-courses">Мои курсы</Link>}</div>
  </section></main>
}
