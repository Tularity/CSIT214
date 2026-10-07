// Times arrive as UTC ISO 8601 and are only converted for display, so sorting
// and comparison always happen on the stored value.
export function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}
