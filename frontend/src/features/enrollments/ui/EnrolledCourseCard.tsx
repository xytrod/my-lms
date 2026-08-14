import { Link } from 'react-router-dom'
import type { CourseDto } from '../../courses/api/types'
import type { ProgressDto } from '../../progress/api/types'
import { ProgressDisplay } from '../../progress/ui/ProgressDisplay'
import type { EnrollmentDto } from '../api/types'
import { courseLevelLabels } from '../../courses/model/courseLabels'

interface Props { enrollment: EnrollmentDto; course: CourseDto; progress?: ProgressDto; progressUnavailable?: boolean }

export function EnrolledCourseCard({ enrollment, course, progress, progressUnavailable }: Props) {
  return <article className="my-course-card">
    <div className="course-card__topline"><span className="course-level">{courseLevelLabels[course.level]}</span><span className="course-state">{enrollment.Status === 'completed' ? 'Завершён' : 'Активен'}</span></div>
    <h2><Link to={`/courses/${course.id}`}>{course.title}</Link></h2><p>{course.description}</p>
    <small>Запись от {new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(new Date(enrollment.StartedAt))}</small>
    {progress && <ProgressDisplay progress={progress} />}{progressUnavailable && <p className="muted-text">Прогресс временно недоступен.</p>}
    <Link className="course-card__link" to={`/courses/${course.id}`}>Открыть курс →</Link>
  </article>
}
