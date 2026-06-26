import type { AnyFieldApi } from '@tanstack/react-form'
import { loginPageStyles } from './form-styles'

export function FieldError({ field }: { field: AnyFieldApi }) {
  if (!field.state.meta.isTouched || field.state.meta.isValid) return null
  return (
    <p role="alert" style={loginPageStyles.fieldError}>
      {field.state.meta.errors.join(', ')}
    </p>
  )
}
