import { useQueries } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getCourse } from '../../features/courses/api/courseApi'
import { courseKeys } from '../../features/courses/model/courseQueries'
import { getCourseProgress } from '../../features/progress/api/progressApi'
import { progressKeys } from '../../features/progress/model/progressQueries'
import { useMyEnrollments } from '../../features/enrollments/model/enrollmentQueries'
import { EnrolledCourseCard } from '../../features/enrollments/ui/EnrolledCourseCard'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'

export function MyCoursesPage() {
  const enrollments = useMyEnrollments()
  const items = enrollments.data ?? []
  const courseQueries = useQueries({ queries: items.map((item) => ({ queryKey: courseKeys.detail(item.CourseID), queryFn: () => getCourse(item.CourseID) })) })
  const progressQueries = useQueries({ queries: items.map((item) => ({ queryKey: progressKeys.course(item.CourseID), queryFn: () => getCourseProgress(item.CourseID) })) })
  const relatedPending = courseQueries.some((query) => query.isPending) || progressQueries.some((query) => query.isPending)

  return <main className="page-container my-courses-page">
    <header className="page-hero"><p className="eyebrow">Ваше обучение</p><h1>Мои курсы</h1><p>Продолжайте обучение и следите за прогрессом по выбранным курсам.</p></header>
    {enrollments.isPending && <section className="panel-state"><Spinner /><p>Загружаем ваши курсы…</p></section>}
    {enrollments.isError && <section className="panel-state panel-state--error"><h2>Не удалось загрузить записи</h2><p>{getApiErrorMessage(enrollments.error)}</p><button className="button" onClick={() => enrollments.refetch()}>Попробовать снова</button></section>}
    {enrollments.isSuccess && items.length === 0 && <section className="panel-state"><h2>Вы ещё не записаны на курсы</h2><p>Откройте каталог и выберите подходящую программу.</p><Link className="button" to="/courses">Перейти в каталог</Link></section>}
    {enrollments.isSuccess && items.length > 0 && relatedPending && <div className="inline-state"><Spinner /><span>Собираем информацию о курсах…</span></div>}
    {enrollments.isSuccess && items.length > 0 && !relatedPending && <section className="my-courses-grid">{items.map((enrollment, index) => {
      const courseQuery = courseQueries[index]; const progressQuery = progressQueries[index]
      if (!courseQuery?.data) return <article className="my-course-card my-course-card--error" key={enrollment.ID}><h2>Курс недоступен</h2><p>Не удалось получить информацию об этом курсе.</p></article>
      return <EnrolledCourseCard key={enrollment.ID} enrollment={enrollment} course={courseQuery.data} progress={progressQuery?.data} progressUnavailable={progressQuery?.isError} />
    })}</section>}
  </main>
}
