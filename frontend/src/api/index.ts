import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
})

// ── Types ──────────────────────────────────────────────────────────────────

export interface News {
  ID: number
  Title: string
  Slug: string
  Excerpt: string
  Thumbnail: string
  Status: string
  PublishedAt: string | null
  CreatedAt: string
  Author: { Name: string }
  Category: { Name: string } | null
  Content?: string
  MetaTitle?: string
  MetaDescription?: string
  ViewCount?: number
}

export interface Volunteer {
  ID: number
  Title: string
  Slug: string
  Division: string
  Location: string
  Type: string
  Description: string
  Requirements: string
  Benefits: string
  Status: string
  CreatedAt: string
}

export interface Banner {
  ID: number
  Title: string
  ImagePath: string
  LinkURL: string
  OrderNum: number
}

export interface Page {
  ID: number
  Slug: string
  Title: string
  Content: string
  MetaTitle: string
  MetaDescription: string
}

export interface GalleryCategory {
  ID: number
  Name: string
  Slug: string
  OrderNum: number
}

export interface GalleryItem {
  ID: number
  Title: string
  Type: 'image' | 'video'
  URL: string
  ThumbnailURL: string
  Caption: string
  Category: GalleryCategory | null
  CreatedAt: string
}

export interface Pagination {
  total: number
  page: number
  limit: number
}

// Site settings: key/value mentah dari backend. Field bertipe "json" berisi
// string JSON yang belum di-parse — gunakan helper getJSON() dari lib/settings.
export type Settings = Record<string, string>

// ── API calls ──────────────────────────────────────────────────────────────

export const portalApi = {
  // ── Home ──────────────────────────────────────────────────────────────────
  home: async () => {
    const { data } = await api.get('/home')
    return data.data as {
      banners: Banner[]
      latest_news: News[]
      latest_volunteers: Volunteer[]
    }
  },

  // ── Konten dinamis (hero, statistik, footer, kontak, tim, dst) ─────────────
  settings: async () => {
    const { data } = await api.get('/settings')
    return data.data as Settings
  },

  // ── Berita ────────────────────────────────────────────────────────────────
  newsList: async (page = 1, q = '') => {
    const { data } = await api.get('/news', { params: { page, limit: 9, q } })
    return data as { data: News[]; meta: Pagination }
  },
  newsDetail: async (slug: string) => {
    const { data } = await api.get(`/news/${slug}`)
    return data.data as News
  },

  // ── Relawan ───────────────────────────────────────────────────────────────
  volunteerList: async (page = 1) => {
    const { data } = await api.get('/volunteers', { params: { page, limit: 10 } })
    return data as { data: Volunteer[]; meta: Pagination }
  },
  volunteerDetail: async (slug: string) => {
    const { data } = await api.get(`/volunteers/${slug}`)
    return data.data as Volunteer
  },
  apply: async (slug: string, formData: FormData) => {
    const { data } = await api.post(`/volunteers/${slug}/apply`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return data as { message: string }
  },

  // ── Galeri ────────────────────────────────────────────────────────────────
  galleryCategories: async () => {
    const { data } = await api.get('/gallery/categories')
    return data.data as GalleryCategory[]
  },
  galleryList: async (page = 1, category = '') => {
    const { data } = await api.get('/gallery', {
      params: { page, limit: 12, ...(category ? { category } : {}) },
    })
    return data as { data: GalleryItem[]; meta: Pagination }
  },

  // ── Halaman statis ────────────────────────────────────────────────────────
  page: async (slug: string) => {
    const { data } = await api.get(`/pages/${slug}`)
    return data.data as Page
  },
}
