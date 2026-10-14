export class ApiError extends Error {
  readonly status: number
  readonly body: unknown

  constructor(message: string, status: number, body: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.body = body
  }
}

const unreachable = 'The server cannot be reached. Check that it is running, then try again.'

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  let response: Response
  try {
    response = await fetch(path, {
      ...init,
      headers: { Accept: 'application/json', ...init.headers },
    })
  } catch {
    throw new Error(unreachable)
  }
  const body: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    const message = errorMessage(body)
    if (message) throw new ApiError(message, response.status, body)
    // Without an {"error": ...} body a gateway error means the development
    // proxy could not reach the backend at all.
    if (response.status >= 502 && response.status <= 504) throw new Error(unreachable)
    throw new Error(`Request failed with status ${response.status}`)
  }
  return body as T
}

export function getJSON<T>(path: string): Promise<T> {
  return request<T>(path)
}

export function postJSON<T>(path: string, payload?: unknown): Promise<T> {
  return request<T>(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: payload === undefined ? undefined : JSON.stringify(payload),
  })
}

// The API reports failures as {"error": "..."}; that wording is shown to the
// operator unchanged so the interface and the backend never disagree.
function errorMessage(body: unknown): string | null {
  if (body && typeof body === 'object' && 'error' in body) {
    const message = (body as { error: unknown }).error
    if (typeof message === 'string' && message !== '') {
      return message
    }
  }
  return null
}
