import type { PropsWithChildren } from 'react'
import { Link } from 'react-router-dom'

interface AuthLayoutProps extends PropsWithChildren {
  title: string
  subtitle: string
}

export function AuthLayout({ title, subtitle, children }: AuthLayoutProps) {
  return (
    <main className="auth-page">
      <section className="auth-brand">
        <Link className="logo" to="/">L<span>MS</span></Link>
        <div>
          <p className="eyebrow">Знания, которые двигают вперёд</p>
          <h1>Учитесь. Создавайте. <em>Растите.</em></h1>
          <p>Единое пространство для студентов и преподавателей.</p>
        </div>
        <div className="brand-orb" aria-hidden="true" />
      </section>
      <section className="auth-panel">
        <div className="auth-card">
          <header><h2>{title}</h2><p>{subtitle}</p></header>
          {children}
        </div>
      </section>
    </main>
  )
}
