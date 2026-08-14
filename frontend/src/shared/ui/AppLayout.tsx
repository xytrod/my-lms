import { NavLink, Outlet } from 'react-router-dom'
import { useAuthStore } from '../../features/auth/model/authStore'
import { useLogout } from '../../features/auth/model/useAuth'
import { Button } from './Button'
import { Spinner } from './Spinner'

const roleLabels = { student: 'Студент', teacher: 'Преподаватель', admin: 'Администратор' } as const
export function AppLayout() {
  const user = useAuthStore((state) => state.user); const logout = useLogout()
  return <div className="app-shell"><header className="site-header"><NavLink className="logo" to="/courses" aria-label="LMS — каталог курсов">L<span>MS</span></NavLink>
    <nav className="main-nav" aria-label="Основная навигация"><NavLink to="/courses">Каталог</NavLink>{user && <NavLink to="/my-courses">Мои курсы</NavLink>}{user && <NavLink to="/creator">Созданные мной</NavLink>}<NavLink to="/dashboard">Кабинет</NavLink></nav>
    <div className="account-nav">{user ? <><div className="user-summary"><span>{roleLabels[user.role]}</span></div><Button className="button-secondary header-button" onClick={() => logout.mutate()} disabled={logout.isPending}>{logout.isPending ? <Spinner /> : 'Выйти'}</Button></> : <><NavLink className="login-link" to="/login">Войти</NavLink><NavLink className="register-link" to="/register">Регистрация</NavLink></>}</div>
  </header><Outlet /></div>
}
