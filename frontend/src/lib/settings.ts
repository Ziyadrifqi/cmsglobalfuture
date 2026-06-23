import { Settings } from '../api'

/**
 * Ambil nilai teks dari settings. Kalau key belum ada / kosong (mis. CMS
 * belum sempat di-restart, atau admin mengosongkan field), pakai fallback
 * supaya tampilan tidak pernah blank.
 */
export function getText(settings: Settings | undefined, key: string, fallback: string): string {
  const v = settings?.[key]
  return v && v.trim() !== '' ? v : fallback
}

/**
 * Ambil & parse nilai JSON (list) dari settings. Kalau key kosong atau
 * format JSON-nya tidak valid (admin salah edit di CMS), otomatis jatuh
 * ke nilai fallback supaya frontend tidak pernah crash gara-gara konten.
 */
export function getJSON<T>(settings: Settings | undefined, key: string, fallback: T): T {
  const raw = settings?.[key]
  if (!raw) return fallback
  try {
    const parsed = JSON.parse(raw)
    return parsed as T
  } catch {
    return fallback
  }
}
