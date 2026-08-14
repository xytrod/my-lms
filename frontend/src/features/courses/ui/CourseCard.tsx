import { Link } from 'react-router-dom'
import type { CourseDto } from '../api/types'
import { courseLevelLabels, getCourseStateLabel } from '../model/courseLabels'

export function CourseCard({ course }: { course: CourseDto }) {
  return (
    <article className="course-card">
      <div className="course-card__topline"><span className="course-level">{courseLevelLabels[course.level]}</span><span className="course-state">{getCourseStateLabel(course.state)}</span></div>
      <h2><Link to={`/courses/${course.id}`}>{course.title}</Link></h2>
      <p>{course.description}</p>
      <footer>
        <span>Опубликован {new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(new Date(course.created_at))}</span>
        <Link className="course-card__link" to={`/courses/${course.id}`} aria-label={`Подробнее о курсе «${course.title}»`}>Подробнее <span aria-hidden="true">→</span></Link>
      </footer>
    </article>
  )
}
