import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import type { CourseDto } from '../../courses/api/types'
import { getApiErrorMessage } from '../../../shared/api/error'
import { Spinner } from '../../../shared/ui/Spinner'
import { courseFormSchema, type CourseFormValues } from '../model/schemas'

interface Props { course?: CourseDto; pending: boolean; error?: unknown; submitLabel: string; onSubmit: (values: CourseFormValues) => Promise<void> }
export function CourseForm({ course, pending, error, submitLabel, onSubmit }: Props) {
  const form = useForm<CourseFormValues>({ resolver: zodResolver(courseFormSchema), defaultValues: { title: course?.title ?? '', description: course?.description ?? '', level: course?.level ?? 'beginner' } })
  const submit = form.handleSubmit(async (values) => {
    await onSubmit(values)
    form.reset(values)
  })
  const unchanged = Boolean(course) && !form.formState.isDirty

  return <form className="creator-form" onSubmit={submit} noValidate>
    <label className="field"><span>Название</span><input {...form.register('title')} aria-invalid={Boolean(form.formState.errors.title)} />{form.formState.errors.title && <small>{form.formState.errors.title.message}</small>}</label>
    <label className="field"><span>Описание</span><textarea rows={5} {...form.register('description')} aria-invalid={Boolean(form.formState.errors.description)} />{form.formState.errors.description && <small>{form.formState.errors.description.message}</small>}</label>
    <label className="field"><span>Уровень</span><select {...form.register('level')} aria-invalid={Boolean(form.formState.errors.level)}><option value="beginner">Начальный</option><option value="intermediate">Средний</option><option value="advanced">Продвинутый</option></select>{form.formState.errors.level && <small>{form.formState.errors.level.message}</small>}</label>
    {error != null && <div className="form-error" role="alert">{getApiErrorMessage(error)}</div>}
    <button type="submit" className="button" disabled={pending || unchanged} title={unchanged ? 'Измените хотя бы одно поле' : undefined}>{pending ? <><Spinner /> Сохраняем…</> : unchanged ? 'Нет изменений' : submitLabel}</button>
  </form>
}
