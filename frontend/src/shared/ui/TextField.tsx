import { forwardRef, type InputHTMLAttributes } from 'react'

interface TextFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  ({ label, error, id, ...props }, ref) => {
    const inputId = id ?? props.name
    return (
      <label className="field" htmlFor={inputId}>
        <span>{label}</span>
        <input ref={ref} id={inputId} aria-invalid={Boolean(error)} {...props} />
        {error && <small role="alert">{error}</small>}
      </label>
    )
  },
)

TextField.displayName = 'TextField'
