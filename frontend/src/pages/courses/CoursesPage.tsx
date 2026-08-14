import { useCourses } from '../../features/courses/model/courseQueries'
import { CourseCard } from '../../features/courses/ui/CourseCard'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'

export function CoursesPage() {
  const courses = useCourses({ limit: 100, offset: 0 })
  return <main className="page-container catalog-page">
    <header className="page-hero"><p className="eyebrow">Учитесь в своём темпе</p><h1>Найдите курс для следующего шага</h1><p>Практические знания от преподавателей платформы — от основ до продвинутого уровня.</p></header>
    {courses.isPending && <section className="panel-state"><Spinner /><p>Загружаем каталог…</p></section>}
    {courses.isError && <section className="panel-state panel-state--error"><h2>Каталог временно недоступен</h2><p>{getApiErrorMessage(courses.error)}</p><button className="button" onClick={() => courses.refetch()}>Попробовать снова</button></section>}
    {courses.isSuccess && courses.data.length === 0 && <section className="panel-state"><h2>Курсов пока нет</h2><p>Опубликованные курсы появятся здесь.</p></section>}
    {courses.isSuccess && courses.data.length > 0 && <section className="course-grid" aria-label="Доступные курсы">{courses.data.map((course) => <CourseCard key={course.id} course={course} />)}</section>}
  </main>
}
