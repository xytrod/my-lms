import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useAuthStore } from '../../features/auth/model/authStore'
import { useCourse, useCourseLessons } from '../../features/courses/model/courseQueries'
import { useMyEnrollments } from '../../features/enrollments/model/enrollmentQueries'
import { useCourseProgress } from '../../features/progress/model/progressQueries'
import { useCompleteLesson } from '../../features/progress/model/useCompleteLesson'
import { ProgressDisplay } from '../../features/progress/ui/ProgressDisplay'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'

export function LessonPage() {
  const { courseId = '', lessonId = '' } = useParams()
  const user = useAuthStore((state) => state.user)
  const course = useCourse(courseId)
  const lessons = useCourseLessons(courseId)
  const enrollments = useMyEnrollments(Boolean(user))
  const enrollment = enrollments.data?.find((item) => item.CourseID === courseId)
  const progress = useCourseProgress(courseId, Boolean(enrollment))
  const completion = useCompleteLesson()
  const [completedThisSession, setCompletedThisSession] = useState(false)
  useEffect(() => setCompletedThisSession(false), [lessonId])
  const orderedLessons = [...(lessons.data ?? [])].sort((a, b) => a.position - b.position || a.id.localeCompare(b.id))
  const lessonIndex = orderedLessons.findIndex((item) => item.id === lessonId)
  const lesson = lessonIndex >= 0 ? orderedLessons[lessonIndex] : undefined
  const previous = lessonIndex > 0 ? orderedLessons[lessonIndex - 1] : undefined
  const next = lessonIndex >= 0 && lessonIndex < orderedLessons.length - 1 ? orderedLessons[lessonIndex + 1] : undefined

  if (course.isPending || lessons.isPending) return <main className="page-container"><section className="panel-state"><Spinner /><p>Загружаем урок…</p></section></main>
  if (course.isError || lessons.isError) return <main className="page-container"><section className="panel-state panel-state--error"><h1>Не удалось открыть урок</h1><p>{getApiErrorMessage(course.error ?? lessons.error)}</p><Link to={`/courses/${courseId}`}>Вернуться к курсу</Link></section></main>
  if (!lesson) return <main className="page-container"><section className="panel-state"><h1>Урок не найден</h1><p>Этот урок не входит в выбранный курс.</p><Link to={`/courses/${courseId}`}>Вернуться к курсу</Link></section></main>

  const canComplete = user?.role === 'student' && Boolean(enrollment)
  const courseCompleted = enrollment?.Status === 'completed' || progress.data?.percentage === 100

  return <main className="page-container lesson-page">
    <nav className="lesson-breadcrumb" aria-label="Навигация по курсу"><Link to={`/courses/${courseId}`}>← {course.data.title}</Link><span>Урок {lesson.position}</span></nav>
    <article className="lesson-viewer">
      <header><div className="course-card__topline"><span className="course-level">Урок {lesson.position}</span>{lesson.is_preview && <span className="preview-badge">Открытый урок</span>}</div><h1>{lesson.title}</h1></header>
      <div className="lesson-content">{lesson.content}</div>
    </article>

    {enrollment && <section className="lesson-progress" aria-label="Прогресс обучения">
      {progress.isPending && <div className="inline-state"><Spinner /><span>Обновляем прогресс…</span></div>}
      {progress.isSuccess && <ProgressDisplay progress={progress.data} />}
      {progress.isError && <div className="inline-error"><p>{getApiErrorMessage(progress.error)}</p><button className="button button-secondary" onClick={() => progress.refetch()}>Повторить</button></div>}
      {courseCompleted && <p className="course-complete-message">Курс завершён</p>}
      {!courseCompleted && canComplete && !completedThisSession && <><button className="button" disabled={completion.isPending} onClick={() => completion.mutate({ courseId, lessonId }, { onSuccess: () => setCompletedThisSession(true) })}>{completion.isPending ? <><Spinner /> Завершаем…</> : 'Завершить урок'}</button>{completion.isError && <p className="form-error">{getApiErrorMessage(completion.error)}</p>}</>}
      {completedThisSession && !courseCompleted && <p className="lesson-complete-message">✓ Урок отмечен завершённым</p>}
    </section>}
    {!user && <aside className="lesson-access-note"><p>Войдите и запишитесь на курс, чтобы отмечать уроки завершёнными.</p><Link to="/login" state={{ from: { pathname: `/courses/${courseId}/lessons/${lessonId}` } }}>Войти</Link></aside>}
    {user && enrollments.isPending && <aside className="lesson-access-note"><div className="inline-state"><Spinner /><span>Проверяем доступ к прогрессу…</span></div></aside>}
    {user && enrollments.isError && <aside className="lesson-access-note"><p>{getApiErrorMessage(enrollments.error)}</p><button className="button button-secondary" onClick={() => enrollments.refetch()}>Повторить</button></aside>}
    {user?.role === 'student' && enrollments.isSuccess && !enrollment && <aside className="lesson-access-note"><p>Запишитесь на курс, чтобы сохранять прогресс.</p><Link to={`/courses/${courseId}`}>Перейти к записи</Link></aside>}

    <nav className="lesson-navigation" aria-label="Переход между уроками">
      {previous ? <Link to={`/courses/${courseId}/lessons/${previous.id}`}><small>Предыдущий урок</small><span>← {previous.title}</span></Link> : <span />}
      {next ? <Link className="lesson-navigation__next" to={`/courses/${courseId}/lessons/${next.id}`}><small>Следующий урок</small><span>{next.title} →</span></Link> : <Link className="lesson-navigation__next" to={`/courses/${courseId}`}><small>Конец программы</small><span>Вернуться к курсу →</span></Link>}
    </nav>
  </main>
}
