import { Link } from 'react-router-dom'
import type { LessonDto } from '../api/types'

export function LessonList({ lessons, courseId }: { lessons: LessonDto[]; courseId: string }) {
  const orderedLessons = [...lessons].sort((a, b) => a.position - b.position || a.id.localeCompare(b.id))
  return <ol className="lesson-list">{orderedLessons.map((lesson) => (
    <li key={lesson.id} className="lesson-item">
      <span className="lesson-position">{lesson.position}</span>
      <div><div className="lesson-heading"><h3><Link to={`/courses/${courseId}/lessons/${lesson.id}`}>{lesson.title}</Link></h3>{lesson.is_preview && <span className="preview-badge">Открытый урок</span>}</div><Link className="lesson-open-link" to={`/courses/${courseId}/lessons/${lesson.id}`}>Открыть урок →</Link></div>
    </li>
  ))}</ol>
}
