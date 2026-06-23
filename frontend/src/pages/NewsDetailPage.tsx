import { useQuery } from '@tanstack/react-query'
import { useParams, Link } from 'react-router-dom'
import { portalApi } from '../api'

export function NewsDetailPage() {
  const { slug } = useParams<{ slug: string }>()

  const { data: news, isLoading, isError } = useQuery({
    queryKey: ['news-detail', slug],
    queryFn: () => portalApi.newsDetail(slug!),
    enabled: !!slug,
  })

  if (isLoading) return (
    <div className="max-w-3xl mx-auto px-4 py-12 space-y-4">
      <div className="flex gap-2 items-center">
        <div className="h-3 w-16 rounded-full bg-teal-100 animate-pulse" />
        <div className="h-3 w-3 rounded-full bg-teal-100 animate-pulse" />
        <div className="h-3 w-24 rounded-full bg-teal-100 animate-pulse" />
      </div>
      <div className="h-8 bg-slate-100 rounded-xl animate-pulse w-3/4" />
      <div className="h-6 bg-slate-100 rounded-xl animate-pulse w-1/2" />
      <div className="h-4 bg-teal-50 rounded-full animate-pulse w-40" />
      <div className="h-72 bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse rounded-2xl" />
      <div className="space-y-3 pt-2">
        {[1, 2, 3, 4, 5, 6, 7].map(i => (
          <div key={i} className={`h-4 rounded-full animate-pulse ${i % 3 === 0 ? 'bg-slate-100 w-4/5' : 'bg-slate-100'}`} />
        ))}
      </div>
    </div>
  )

  if (isError || !news) return (
    <div className="max-w-3xl mx-auto px-4 py-24 text-center">
      <div className="text-6xl mb-5">📭</div>
      <h1 className="text-xl font-bold text-teal-900 mb-3">Berita tidak ditemukan</h1>
      <Link to="/berita" className="text-teal-600 hover:underline text-sm font-medium">← Kembali ke daftar berita</Link>
    </div>
  )

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">

      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-sm text-slate-400 mb-8">
        <Link to="/" className="hover:text-teal-600 transition-colors">Beranda</Link>
        <span>/</span>
        <Link to="/berita" className="hover:text-teal-600 transition-colors">Berita</Link>
        <span>/</span>
        <span className="text-slate-600 line-clamp-1 max-w-[200px]">{news.Title}</span>
      </nav>

      {/* Category */}
      {news.Category && (
        <span className="inline-block px-3 py-1 bg-teal-100 text-teal-700 rounded-full text-xs font-bold">
          {news.Category.Name}
        </span>
      )}

      {/* Title */}
      <h1 className="text-3xl md:text-4xl font-black text-teal-900 mt-4 mb-5 leading-tight tracking-tight">
        {news.Title}
      </h1>

      {/* Meta */}
      <div className="flex flex-wrap items-center gap-4 text-sm text-slate-400 pb-6 mb-8 border-b border-teal-100">
        <span>✍️ {news.Author?.Name}</span>
        {news.PublishedAt && (
          <span>📅 {new Date(news.PublishedAt).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}</span>
        )}
        {news.ViewCount !== undefined && (
          <span>👁 {news.ViewCount.toLocaleString()} kali dibaca</span>
        )}
      </div>

      {/* Thumbnail */}
      {news.Thumbnail ? (
        <img
          src={news.Thumbnail} alt={news.Title}
          className="w-full rounded-2xl mb-10 object-cover max-h-96"
        />
      ) : (
        <div className="w-full h-64 bg-gradient-to-br from-teal-100 to-teal-200 rounded-2xl flex items-center justify-center text-6xl mb-10">
          📰
        </div>
      )}

      {/* Content */}
      <div className="prose prose-slate max-w-none text-slate-700 leading-[1.9] whitespace-pre-wrap text-[15px]">
        {news.Content}
      </div>

      {/* Back */}
      <div className="mt-12 pt-6 border-t border-teal-100">
        <Link to="/berita" className="text-teal-600 hover:text-teal-800 font-semibold text-sm transition-colors">
          ← Kembali ke daftar berita
        </Link>
      </div>
    </div>
  )
}