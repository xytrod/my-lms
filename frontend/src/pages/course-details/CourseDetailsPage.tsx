import axios from 'axios'
import { Link, useParams } from 'react-router-dom'
import { useCourse, useCourseLessons } from '../../features/courses/model/courseQueries'
import { LessonList } from '../../features/courses/ui/LessonList'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'
import { EnrollmentAction } from '../../features/enrollments/ui/EnrollmentAction'
import { courseLevelLabels, getCourseStateLabel } from '../../features/courses/model/courseLabels'

export function CourseDetailsPage() {
  const { courseId = '' } = useParams()
  const course = useCourse(courseId)
  const lessons = useCourseLessons(courseId)
  const notFound = axios.isAxiosError(course.error) && course.error.response?.status === 404

  if (course.isPending) return <main className="page-container"><section className="panel-state"><Spinner /><p>Загружаем курс…</p></section></main>
  if (course.isError) return <main className="page-container"><section className="panel-state panel-state--error"><h1>{notFound ? 'Курс не найден' : 'Не удалось загрузить курс'}</h1><p>{notFound ? 'Возможно, курс ещё не опубликован или был удалён.' : getApiErrorMessage(course.error)}</p><Link to="/courses">Вернуться в каталог</Link></section></main>

  return <main className="page-container course-details">
    <Link className="back-link" to="/courses">← Все курсы</Link>
    <section className="course-details__hero">
      <div><div className="course-card__topline"><span className="course-level">{courseLevelLabels[course.data.level]}</span><span className="course-state">{getCourseStateLabel(course.data.state)}</span></div><h1>{course.data.title}</h1><p>{course.data.description}</p></div>
      <dl className="course-meta"><div><dt>Преподаватель</dt><dd>Автор курса</dd></div><div><dt>Опубликован</dt><dd>{new Intl.DateTimeFormat('ru-RU', { dateStyle: 'long' }).format(new Date(course.data.created_at))}</dd></div></dl>
    </section>
    <EnrollmentAction courseId={courseId} />
    <section className="lessons-section">
      <header><p className="eyebrow">Программа</p><h2>Уроки курса</h2></header>
      {lessons.isPending && <div className="inline-state"><Spinner /><span>Загружаем уроки…</span></div>}
      {lessons.isError && <div className="inline-error"><p>{getApiErrorMessage(lessons.error)}</p><button className="button button-secondary" onClick={() => lessons.refetch()}>Повторить</button></div>}
      {lessons.isSuccess && lessons.data.length === 0 && <div className="empty-lessons"><h3>Уроков пока нет</h3><p>Программа этого курса ещё не наполнена.</p></div>}
      {lessons.isSuccess && lessons.data.length > 0 && <LessonList lessons={lessons.data} courseId={courseId} />}
    </section>
  </main>
}
