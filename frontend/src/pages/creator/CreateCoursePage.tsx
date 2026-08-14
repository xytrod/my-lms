import { Link, useNavigate } from 'react-router-dom'
import { CourseForm } from '../../features/creator/ui/CourseForm'
import { useCreateCourse } from '../../features/creator/model/creatorQueries'
import type { CourseFormValues } from '../../features/creator/model/schemas'

export function CreateCoursePage() {
  const mutation = useCreateCourse(); const navigate = useNavigate()
  const submit = async (values: CourseFormValues) => { const course = await mutation.mutateAsync(values); navigate(`/creator/courses/${course.id}`, { replace: true }) }
  return <main className="page-container narrow-page"><Link className="back-link" to="/creator">← Созданные курсы</Link><header className="form-page-heading"><p className="eyebrow">Новая программа</p><h1>Создать курс</h1><p>Курс сохранится как черновик. Для публикации добавьте хотя бы один урок.</p></header><CourseForm pending={mutation.isPending} error={mutation.error} submitLabel="Создать курс" onSubmit={submit} /></main>
}
