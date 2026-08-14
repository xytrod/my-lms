import { Link } from 'react-router-dom'
import { useCreatedCourses } from '../../features/creator/model/creatorQueries'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'
import { courseLevelLabels, getCourseStateLabel } from '../../features/courses/model/courseLabels'

export function CreatorCoursesPage() {
  const courses = useCreatedCourses()
  return <main className="page-container creator-page"><header className="creator-heading"><div><p className="eyebrow">Авторская зона</p><h1>Созданные мной</h1><p>Создавайте программы, наполняйте их уроками и управляйте публикацией.</p></div><Link className="button" to="/creator/courses/new">Создать курс</Link></header>
    {courses.isPending && <section className="panel-state"><Spinner /><p>Загружаем курсы…</p></section>}
    {courses.isError && <section className="panel-state panel-state--error"><h2>Не удалось загрузить курсы</h2><p>{getApiErrorMessage(courses.error)}</p><button className="button" onClick={() => courses.refetch()}>Повторить</button></section>}
    {courses.isSuccess && courses.data.length === 0 && <section className="panel-state"><h2>У вас пока нет созданных курсов</h2><p>Начните с первой учебной программы.</p><Link className="button" to="/creator/courses/new">Создать курс</Link></section>}
    {courses.isSuccess && courses.data.length > 0 && <section className="creator-grid">{courses.data.map((course) => <article className="creator-card" key={course.id}><div className="course-card__topline"><span className={`creator-state creator-state--${course.state}`}>{getCourseStateLabel(course.state)}</span><span className="course-level">{courseLevelLabels[course.level]}</span></div><h2>{course.title}</h2><p>{course.description}</p><small>Обновлён {new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(new Date(course.updated_at))}</small><Link to={`/creator/courses/${course.id}`}>Управлять →</Link></article>)}</section>}
  </main>
}
