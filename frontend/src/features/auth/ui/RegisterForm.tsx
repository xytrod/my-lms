import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { getApiErrorMessage } from '../../../shared/api/error'
import { Button } from '../../../shared/ui/Button'
import { Spinner } from '../../../shared/ui/Spinner'
import { TextField } from '../../../shared/ui/TextField'
import { registerSchema, type RegisterFormValues } from '../model/schemas'
import { useRegister } from '../model/useAuth'

export function RegisterForm() {
  const navigate = useNavigate()
  const mutation = useRegister()
  const { register, handleSubmit, formState: { errors } } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: '', password: '', username: '', first_name: '', last_name: '' },
  })

  const onSubmit = handleSubmit(async (values) => {
    await mutation.mutateAsync(values)
    navigate('/dashboard', { replace: true })
  })

  return (
    <form onSubmit={onSubmit} noValidate>
      <div className="field-row">
        <TextField label="Имя" autoComplete="given-name" error={errors.first_name?.message} {...register('first_name')} />
        <TextField label="Фамилия" autoComplete="family-name" error={errors.last_name?.message} {...register('last_name')} />
      </div>
      <TextField label="Имя пользователя" autoComplete="username" placeholder="Минимум 7 символов" error={errors.username?.message} {...register('username')} />
      <TextField label="Email" type="email" autoComplete="email" placeholder="you@example.com" error={errors.email?.message} {...register('email')} />
      <TextField label="Пароль" type="password" autoComplete="new-password" placeholder="От 8 до 20 символов" error={errors.password?.message} {...register('password')} />
      {mutation.isError && <div className="form-error" role="alert">{getApiErrorMessage(mutation.error)}</div>}
      <Button type="submit" disabled={mutation.isPending}>{mutation.isPending ? <><Spinner /> Создаём…</> : 'Создать аккаунт'}</Button>
    </form>
  )
}
