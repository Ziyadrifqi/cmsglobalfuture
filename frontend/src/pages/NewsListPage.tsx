import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link, useSearchParams } from 'react-router-dom'
import { portalApi, News } from '../api'

export function NewsListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [q, setQ] = useState(searchParams.get('q') ?? '')
  const page = Number(searchParams.get('page') ?? 1)

  const { data, isLoading } = useQuery({
    queryKey: ['news-list', page, searchParams.get('q') ?? ''],
    queryFn: () => portalApi.newsList(page, searchParams.get('q') ?? ''),
  })

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setSearchParams(q ? { q } : {})
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">

      {/* Header */}
      <div className="mb-10">
        <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-2">Informasi Terkini</p>
        <h1 className="text-4xl font-black text-teal-900 tracking-tight mb-2">Berita & Artikel</h1>
        <p className="text-slate-500 text-[15px]">Cerita, insight, dan kabar terbaru dari kegiatan yayasan kami.</p>
      </div>

      {/* Search */}
      <form onSubmit={handleSearch} className="mb-10 flex gap-2 max-w-lg">
        <input
          type="text" value={q} onChange={e => setQ(e.target.value)}
          placeholder="Cari berita atau artikel..."
          className="flex-1 px-4 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors bg-white text-teal-900 placeholder:text-slate-400"
        />
        <button
          type="submit"
          className="px-5 py-2.5 bg-teal-600 text-white rounded-xl text-sm font-semibold hover:bg-teal-700 transition-colors"
        >
          Cari
        </button>
        {searchParams.get('q') && (
          <button
            type="button" onClick={() => { setQ(''); setSearchParams({}) }}
            className="px-4 py-2.5 rounded-xl border-2 border-teal-100 text-sm text-slate-500 hover:bg-slate-50 transition-colors"
          >
            Reset
          </button>
        )}
      </form>

      {/* Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map(i => (
            <div key={i} className="rounded-2xl border border-teal-100 overflow-hidden">
              <div className="h-44 bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse" />
              <div className="p-4 space-y-2">
                <div className="h-3 rounded-full bg-teal-100 animate-pulse w-1/3" />
                <div className="h-4 rounded-full bg-slate-100 animate-pulse" />
                <div className="h-4 rounded-full bg-slate-100 animate-pulse w-5/6" />
                <div className="h-3 rounded-full bg-slate-50 animate-pulse w-1/2 mt-2" />
              </div>
            </div>
          ))}
        </div>
      ) : data?.data.length === 0 ? (
        <div className="text-center py-20">
          <div className="text-5xl mb-4">🔍</div>
          <p className="text-slate-400 text-lg">Tidak ada berita ditemukan</p>
          <p className="text-sm text-slate-400 mt-1">Coba kata kunci yang berbeda</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {data?.data.map(news => <NewsCard key={news.ID} news={news} />)}
        </div>
      )}

      {/* Pagination */}
      {data && data.meta.total > 9 && (
        <div className="flex justify-center items-center gap-3 mt-12">
          <button
            disabled={page <= 1}
            onClick={() => setSearchParams({ page: String(page - 1) })}
            className="px-4 py-2 rounded-xl border-2 border-teal-100 text-sm font-medium text-teal-700 disabled:opacity-40 hover:border-teal-400 transition-colors"
          >
            ← Sebelumnya
          </button>
          <span className="px-4 py-2 text-sm text-slate-500 font-medium">Halaman {page}</span>
          <button
            disabled={page * 9 >= data.meta.total}
            onClick={() => setSearchParams({ page: String(page + 1) })}
            className="px-4 py-2 rounded-xl border-2 border-teal-100 text-sm font-medium text-teal-700 disabled:opacity-40 hover:border-teal-400 transition-colors"
          >
            Berikutnya →
          </button>
        </div>
      )}
    </div>
  )
}

function NewsCard({ news }: { news: News }) {
  const catColor: Record<string, string> = {
    Pendidikan: 'bg-teal-100 text-teal-700',
    Program:    'bg-emerald-100 text-emerald-700',
    Kemitraan:  'bg-amber-100 text-amber-800',
    Ekonomi:    'bg-green-100 text-green-700',
    Laporan:    'bg-slate-100 text-slate-600',
    default:    'bg-slate-100 text-slate-600',
  }
  const cc = (news.Category?.Name && catColor[news.Category.Name]) || catColor.default

  return (
    <Link
      to={`/berita/${news.Slug}`}
      className="group bg-white rounded-2xl border border-teal-100 overflow-hidden hover:shadow-[0_12px_40px_rgba(15,118,110,0.14)] hover:-translate-y-1 transition-all duration-200"
    >
      {news.Thumbnail ? (
        <img src={news.Thumbnail} alt={news.Title} className="w-full h-44 object-cover" />
      ) : (
        <div className="w-full h-44 bg-gradient-to-br from-teal-100 to-teal-200 flex items-center justify-center text-4xl relative">
          📰
          {news.Category && (
            <span className={`absolute top-3 right-3 text-xs font-semibold px-2.5 py-0.5 rounded-full ${cc}`}>
              {news.Category.Name}
            </span>
          )}
        </div>
      )}
      <div className="p-5">
        {news.Thumbnail && news.Category && (
          <span className={`text-xs font-semibold px-2.5 py-0.5 rounded-full ${cc}`}>{news.Category.Name}</span>
        )}
        <h3 className="font-bold text-teal-900 mt-2 line-clamp-2 leading-snug group-hover:text-teal-700 transition-colors">
          {news.Title}
        </h3>
        {news.Excerpt && (
          <p className="text-sm text-slate-500 mt-1.5 line-clamp-2 leading-relaxed">{news.Excerpt}</p>
        )}
        <div className="flex items-center justify-between mt-3 text-xs text-slate-400">
          <span>
            {news.PublishedAt
              ? new Date(news.PublishedAt).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
              : ''}
          </span>
          {news.ViewCount !== undefined && <span>👁 {news.ViewCount.toLocaleString()}</span>}
        </div>
      </div>
    </Link>
  )
}