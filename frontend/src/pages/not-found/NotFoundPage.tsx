import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return <main className="center-state"><strong className="error-code">404</strong><h1>Страница не найдена</h1><Link to="/courses">Вернуться в каталог</Link></main>
}
