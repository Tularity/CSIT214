async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: { Accept: 'application/json', ...init.headers },
  })
  const body: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(errorMessage(body) ?? `Request failed with status ${response.status}`)
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
