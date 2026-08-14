import { Link } from 'react-router-dom'
import { useAuthStore } from '../../auth/model/authStore'
import { useEnrollInCourse, useMyEnrollments } from '../model/enrollmentQueries'
import { useCourseProgress } from '../../progress/model/progressQueries'
import { ProgressDisplay } from '../../progress/ui/ProgressDisplay'
import { getApiErrorMessage } from '../../../shared/api/error'
import { Spinner } from '../../../shared/ui/Spinner'

export function EnrollmentAction({ courseId }: { courseId: string }) {
  const user = useAuthStore((state) => state.user)
  const enrollments = useMyEnrollments(user?.role === 'student')
  const enroll = useEnrollInCourse()
  const enrollment = enrollments.data?.find((item) => item.CourseID === courseId)
  const progress = useCourseProgress(courseId, Boolean(enrollment))

  if (!user) return <aside className="enrollment-panel"><h2>Хотите начать обучение?</h2><p>Войдите, чтобы записаться на курс и видеть свой прогресс.</p><Link className="button" to="/login" state={{ from: { pathname: `/courses/${courseId}` } }}>Войти и записаться</Link></aside>
  if (user.role !== 'student') return null
  if (enrollments.isPending) return <aside className="enrollment-panel enrollment-panel--inline"><Spinner /><span>Проверяем запись на курс…</span></aside>
  if (enrollments.isError) return <aside className="enrollment-panel"><p className="form-error">{getApiErrorMessage(enrollments.error)}</p><button className="button button-secondary" onClick={() => enrollments.refetch()}>Повторить</button></aside>
  if (enrollment) return <aside className="enrollment-panel"><div className="enrolled-heading"><span aria-hidden="true">✓</span><div><h2>Вы записаны на курс</h2><p>Статус: {enrollment.Status === 'completed' ? 'завершён' : 'активен'}</p></div></div>{progress.isPending && <div className="inline-state"><Spinner /><span>Загружаем прогресс…</span></div>}{progress.isSuccess && <ProgressDisplay progress={progress.data} />}{progress.isError && <p className="muted-text">Не удалось загрузить прогресс.</p>}<Link to="/my-courses">Перейти к моим курсам →</Link></aside>

  return <aside className="enrollment-panel"><h2>Готовы начать?</h2><p>Запишитесь на курс, чтобы он появился в разделе «Мои курсы».</p>{enroll.isError && <p className="form-error">{getApiErrorMessage(enroll.error)}</p>}<button className="button" disabled={enroll.isPending} onClick={() => enroll.mutate(courseId)}>{enroll.isPending ? <><Spinner /> Записываем…</> : 'Записаться на курс'}</button></aside>
}
