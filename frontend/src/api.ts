export async function getJSON<T>(path: string): Promise<T> {
  const response = await fetch(path, { headers: { Accept: 'application/json' } })
  const body: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    throw new Error(errorMessage(body) ?? `Request failed with status ${response.status}`)
  }
  return body as T
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
