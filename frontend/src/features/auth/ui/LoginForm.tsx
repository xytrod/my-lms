import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { getApiErrorMessage } from '../../../shared/api/error'
import { Button } from '../../../shared/ui/Button'
import { Spinner } from '../../../shared/ui/Spinner'
import { TextField } from '../../../shared/ui/TextField'
import { loginSchema, type LoginFormValues } from '../model/schemas'
import { useLogin } from '../model/useAuth'

export function LoginForm() {
  const navigate = useNavigate()
  const mutation = useLogin()
  const { register, handleSubmit, formState: { errors } } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  })

  const onSubmit = handleSubmit(async (values) => {
    await mutation.mutateAsync(values)
    navigate('/dashboard', { replace: true })
  })

  return (
    <form onSubmit={onSubmit} noValidate>
      <TextField label="Email" type="email" autoComplete="email" placeholder="you@example.com" error={errors.email?.message} {...register('email')} />
      <TextField label="Пароль" type="password" autoComplete="current-password" placeholder="Не менее 8 символов" error={errors.password?.message} {...register('password')} />
      {mutation.isError && <div className="form-error" role="alert">{getApiErrorMessage(mutation.error)}</div>}
      <Button type="submit" disabled={mutation.isPending}>{mutation.isPending ? <><Spinner /> Входим…</> : 'Войти'}</Button>
    </form>
  )
}
