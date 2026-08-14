import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import type { LessonDto } from '../../features/courses/api/types'
import { useArchiveCourse, useCreateLesson, useCreatedCourse, useCreatedCourseLessons, useDeleteLesson, usePublishCourse, useUpdateCourse, useUpdateLesson } from '../../features/creator/model/creatorQueries'
import type { CourseFormValues, LessonFormValues } from '../../features/creator/model/schemas'
import { CourseForm } from '../../features/creator/ui/CourseForm'
import { LessonForm } from '../../features/creator/ui/LessonForm'
import { getApiErrorMessage } from '../../shared/api/error'
import { Spinner } from '../../shared/ui/Spinner'
import { courseLevelLabels, getCourseStateLabel } from '../../features/courses/model/courseLabels'

export function ManageCoursePage() {
  const { courseId = '' } = useParams()
  const courseQuery = useCreatedCourse(courseId)
  const lessonsQuery = useCreatedCourseLessons(courseId)
  const updateCourse = useUpdateCourse(); const publish = usePublishCourse(); const archive = useArchiveCourse()
  const createLesson = useCreateLesson(courseId); const updateLesson = useUpdateLesson(courseId); const removeLesson = useDeleteLesson(courseId)
  const [editingLesson, setEditingLesson] = useState<LessonDto>()

  if (courseQuery.isPending) return <main className="page-container"><section className="panel-state"><Spinner /><p>Загружаем курс…</p></section></main>
  if (courseQuery.isError) return <main className="page-container"><section className="panel-state panel-state--error"><h1>Не удалось загрузить курс</h1><p>{getApiErrorMessage(courseQuery.error)}</p><Link to="/creator">Вернуться к созданным курсам</Link></section></main>
  const course = courseQuery.data
  const isDraft = course.state === 'Draft'; const isPublished = course.state === 'published'; const isClosed = course.state === 'closed'
  const orderedLessons = [...(lessonsQuery.data ?? [])].sort((a, b) => a.position - b.position || a.id.localeCompare(b.id))

  const saveCourse = async (values: CourseFormValues) => {
    const body: Partial<CourseFormValues> = {}
    if (values.title !== course.title) body.title = values.title
    if (values.description !== course.description) body.description = values.description
    if (values.level !== course.level) body.level = values.level
    if (Object.keys(body).length) await updateCourse.mutateAsync({ id: course.id, body })
  }
  const saveLesson = async (values: LessonFormValues) => {
    if (!editingLesson) { await createLesson.mutateAsync(values); return }
    const body: Partial<LessonFormValues> = {}
    if (values.title !== editingLesson.title) body.title = values.title
    if (values.content !== editingLesson.content) body.content = values.content
    if (values.position !== editingLesson.position) body.position = values.position
    if (values.is_preview !== editingLesson.is_preview) body.is_preview = values.is_preview
    if (Object.keys(body).length) await updateLesson.mutateAsync({ id: editingLesson.id, body })
    setEditingLesson(undefined)
  }
  const deleteItem = async (lesson: LessonDto) => {
    if (window.confirm(`Удалить урок «${lesson.title}»?`)) await removeLesson.mutateAsync(lesson.id)
  }
  const publishCourse = () => {
    if (orderedLessons.length === 0) return
    publish.mutate(course.id)
  }

  return <main className="page-container creator-manage"><Link className="back-link" to="/creator">← Созданные курсы</Link>
    <header className="manage-heading"><div><span className={`creator-state creator-state--${course.state}`}>{getCourseStateLabel(course.state)}</span><h1>{course.title}</h1></div><div className="state-actions">{isDraft && <button className="button" disabled={publish.isPending || lessonsQuery.isPending || orderedLessons.length === 0} title={orderedLessons.length === 0 ? 'Добавьте хотя бы один урок' : undefined} onClick={publishCourse}>{publish.isPending ? <Spinner /> : 'Опубликовать'}</button>}{!isClosed && <button className="button button-danger" disabled={archive.isPending} onClick={() => window.confirm('Архивировать курс? Вернуть его в публикацию через текущий API нельзя.') && archive.mutate(course.id)}>{archive.isPending ? <Spinner /> : 'Архивировать'}</button>}</div></header>
    {isDraft && lessonsQuery.isSuccess && orderedLessons.length === 0 && <div className="creator-limitation">Для публикации добавьте хотя бы один урок.</div>}
    {(publish.isError || archive.isError) && <div className="form-error" role="alert">{getApiErrorMessage(publish.error ?? archive.error)}</div>}
    {isPublished && <div className="creator-limitation" role="note">Структура опубликованного курса заблокирована. Уроки доступны только для просмотра.</div>}
    {isClosed && <div className="creator-warning">Архивированный курс доступен только для просмотра. Метаданные и структура уроков заблокированы.</div>}

    <section className="creator-section"><h2>Информация о курсе</h2>{isClosed ? <dl className="readonly-course"><div><dt>Описание</dt><dd>{course.description}</dd></div><div><dt>Уровень</dt><dd>{courseLevelLabels[course.level]}</dd></div></dl> : <CourseForm key={course.updated_at} course={course} pending={updateCourse.isPending} error={updateCourse.error} submitLabel="Сохранить изменения" onSubmit={saveCourse} />}</section>

    <section className="creator-section"><h2>Уроки</h2>
      {lessonsQuery.isPending && <div className="inline-state"><Spinner /><span>Загружаем уроки…</span></div>}
      {lessonsQuery.isError && <div className="inline-error"><p>{getApiErrorMessage(lessonsQuery.error)}</p><button className="button button-secondary" onClick={() => lessonsQuery.refetch()}>Повторить</button></div>}
      {isDraft && lessonsQuery.isSuccess && <div className="lesson-editor"><h3>{editingLesson ? 'Редактировать урок' : 'Добавить урок'}</h3><LessonForm key={editingLesson?.id ?? `new-${orderedLessons.length}`} lesson={editingLesson} nextPosition={orderedLessons.length + 1} pending={createLesson.isPending || updateLesson.isPending} error={createLesson.error ?? updateLesson.error} onCancel={editingLesson ? () => setEditingLesson(undefined) : undefined} onSubmit={saveLesson} /></div>}
      {lessonsQuery.isSuccess && orderedLessons.length === 0 && <div className="empty-lessons"><h3>Уроков пока нет</h3><p>{isDraft ? 'Добавьте первый урок, чтобы подготовить курс к публикации.' : 'В этом курсе нет уроков.'}</p></div>}
      {lessonsQuery.isSuccess && orderedLessons.length > 0 && <ol className="creator-lessons">{orderedLessons.map((lesson) => <li key={lesson.id}><span className="lesson-position">{lesson.position}</span><div><h3>{lesson.title}</h3><p>{lesson.content}</p>{lesson.is_preview && <span className="preview-badge">Открытый урок</span>}</div>{isDraft && <div className="lesson-actions"><button className="button button-secondary" onClick={() => setEditingLesson(lesson)}>Изменить</button><button className="button button-danger" disabled={removeLesson.isPending} onClick={() => void deleteItem(lesson)}>Удалить</button></div>}</li>)}</ol>}
      {(createLesson.isError || updateLesson.isError || removeLesson.isError) && <div className="form-error" role="alert">{getApiErrorMessage(createLesson.error ?? updateLesson.error ?? removeLesson.error)}</div>}
    </section>
  </main>
}
