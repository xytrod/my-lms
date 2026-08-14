import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import type { LessonDto } from '../../courses/api/types'
import { getApiErrorMessage } from '../../../shared/api/error'
import { Spinner } from '../../../shared/ui/Spinner'
import { lessonFormSchema, type LessonFormValues } from '../model/schemas'

interface Props { lesson?: LessonDto; nextPosition?: number; pending: boolean; error?: unknown; onCancel?: () => void; onSubmit: (values: LessonFormValues) => Promise<void> }
export function LessonForm({ lesson, nextPosition = 1, pending, error, onCancel, onSubmit }: Props) {
  const form = useForm<LessonFormValues>({ resolver: zodResolver(lessonFormSchema), defaultValues: { title: lesson?.title ?? '', content: lesson?.content ?? '', position: lesson?.position ?? nextPosition, is_preview: lesson?.is_preview ?? false } })
  return <form className="creator-form lesson-form" onSubmit={form.handleSubmit(onSubmit)} noValidate>
    <div className="field-row"><label className="field"><span>Название урока</span><input {...form.register('title')} />{form.formState.errors.title && <small>{form.formState.errors.title.message}</small>}</label><label className="field"><span>Позиция</span><input type="number" min={1} {...form.register('position', { valueAsNumber: true })} />{form.formState.errors.position && <small>{form.formState.errors.position.message}</small>}</label></div>
    <label className="field"><span>Содержание</span><textarea rows={8} {...form.register('content')} />{form.formState.errors.content && <small>{form.formState.errors.content.message}</small>}</label>
    <label className="check-field"><input type="checkbox" {...form.register('is_preview')} /><span>Отметить как открытый урок</span></label>
    {error != null && <div className="form-error" role="alert">{getApiErrorMessage(error)}</div>}
    <div className="form-actions"><button className="button" disabled={pending}>{pending ? <><Spinner /> Сохраняем…</> : lesson ? 'Сохранить урок' : 'Добавить урок'}</button>{onCancel && <button type="button" className="button button-secondary" onClick={onCancel}>Отмена</button>}</div>
  </form>
}
