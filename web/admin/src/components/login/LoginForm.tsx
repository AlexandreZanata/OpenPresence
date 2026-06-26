import { useState } from 'react'
import { useForm } from '@tanstack/react-form'
import { ApiError } from '~/lib/api/client'
import type { AuthState } from '~/lib/auth/types'
import { FieldError } from './FieldError'
import { loginPageStyles as styles } from './form-styles'

type LoginFormValues = {
  registrationId: string
  password: string
}

type Props = {
  auth: AuthState
  onSuccess: () => void
}

function required(value: string): string | undefined {
  return value.trim() ? undefined : 'This field is required'
}

function minPassword(value: string): string | undefined {
  if (!value.trim()) return 'This field is required'
  if (value.length < 4) return 'Password must be at least 4 characters'
  return undefined
}

export function LoginForm({ auth, onSuccess }: Props) {
  const [apiError, setApiError] = useState<string | null>(null)

  const form = useForm({
    defaultValues: {
      registrationId: '',
      password: '',
    } satisfies LoginFormValues,
    onSubmit: async ({ value }) => {
      setApiError(null)
      try {
        await auth.login(value)
        onSuccess()
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          setApiError('Invalid credentials')
          return
        }
        if (err instanceof ApiError) {
          setApiError(err.message)
          return
        }
        setApiError('Login failed')
      }
    },
  })

  return (
    <form
      style={styles.form}
      onSubmit={(event) => {
        event.preventDefault()
        event.stopPropagation()
        void form.handleSubmit()
      }}
      noValidate
    >
      <form.Field
        name="registrationId"
        validators={{ onChange: ({ value }) => required(value) }}
      >
        {(field) => (
          <label style={styles.label} htmlFor={field.name}>
            Registration ID
            <input
              id={field.name}
              name={field.name}
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
              autoComplete="username"
              aria-invalid={field.state.meta.isTouched && !field.state.meta.isValid}
              style={{
                ...styles.input,
                ...(field.state.meta.isTouched && !field.state.meta.isValid
                  ? styles.inputInvalid
                  : {}),
              }}
            />
            <FieldError field={field} />
          </label>
        )}
      </form.Field>

      <form.Field
        name="password"
        validators={{ onChange: ({ value }) => minPassword(value) }}
      >
        {(field) => (
          <label style={styles.label} htmlFor={field.name}>
            Password
            <input
              id={field.name}
              name={field.name}
              type="password"
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
              autoComplete="current-password"
              aria-invalid={field.state.meta.isTouched && !field.state.meta.isValid}
              style={{
                ...styles.input,
                ...(field.state.meta.isTouched && !field.state.meta.isValid
                  ? styles.inputInvalid
                  : {}),
              }}
            />
            <FieldError field={field} />
          </label>
        )}
      </form.Field>

      {apiError ? (
        <p role="alert" style={styles.apiError}>
          {apiError}
        </p>
      ) : null}

      <form.Subscribe selector={(state) => [state.canSubmit, state.isSubmitting]}>
        {([canSubmit, isSubmitting]) => (
          <button
            type="submit"
            disabled={!canSubmit || isSubmitting}
            style={{
              ...styles.button,
              ...(!canSubmit || isSubmitting ? styles.buttonDisabled : {}),
            }}
          >
            {isSubmitting ? 'Signing in…' : 'Sign in'}
          </button>
        )}
      </form.Subscribe>
    </form>
  )
}
