export function imgUrl(path: string | null | undefined): string {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  const base = import.meta.env.VITE_BACKEND_URL ?? ''
  return `${base}${path}`
}
