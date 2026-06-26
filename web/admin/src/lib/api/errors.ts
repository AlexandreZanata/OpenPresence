export type ApiErrorBody = {
  error: {
    code: string
    message: string
    correlationId?: string
  }
}

export class ApiError extends Error {
  readonly code: string
  readonly status: number
  readonly correlationId?: string

  constructor(code: string, message: string, status: number, correlationId?: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
    this.correlationId = correlationId
  }
}

export async function parseApiError(response: Response): Promise<ApiError> {
  try {
    const body = (await response.json()) as ApiErrorBody
    if (body?.error?.code) {
      return new ApiError(
        body.error.code,
        body.error.message,
        response.status,
        body.error.correlationId,
      )
    }
  } catch {
    // fall through
  }
  return new ApiError('HTTP_ERROR', response.statusText || 'Request failed', response.status)
}
