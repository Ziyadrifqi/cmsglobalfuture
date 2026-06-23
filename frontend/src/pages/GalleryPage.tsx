import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { portalApi, GalleryItem, GalleryCategory } from '../api'

export function GalleryPage() {
  const [activeCategory, setActiveCategory] = useState<GalleryCategory | null>(null)
  const [page, setPage]                     = useState(1)
  const [lightbox, setLightbox]             = useState<GalleryItem | null>(null)

  // Ambil daftar kategori dari API
  const { data: categories = [], isLoading: loadCat } = useQuery({
    queryKey: ['gallery-categories'],
    queryFn:  portalApi.galleryCategories,
    staleTime: 5 * 60 * 1000,
  })

  const { data, isLoading } = useQuery({
    queryKey: ['gallery', page, activeCategory?.Slug ?? ''],
    queryFn:  () => portalApi.galleryList(page, activeCategory?.Slug ?? ''),
  })

  const handleCategoryChange = (cat: GalleryCategory | null) => {
    setActiveCategory(cat)
    setPage(1)
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">

      {/* Header */}
      <div className="mb-10">
        <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-2">Dokumentasi Kegiatan</p>
        <h1 className="text-4xl font-black text-teal-900 tracking-tight mb-2">Galeri</h1>
        <p className="text-slate-500 text-[15px]">Momen-momen berharga dari aksi nyata relawan Green Future Indonesia di seluruh penjuru negeri.</p>
      </div>

      {/* Category filter — dinamis dari API */}
      <div className="flex gap-2 flex-wrap mb-8">
        <CategoryButton
          label="Semua"
          active={activeCategory === null}
          onClick={() => handleCategoryChange(null)}
        />
        {loadCat
          ? Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-9 w-28 rounded-xl bg-teal-50 animate-pulse" />
            ))
          : categories.map(cat => (
              <CategoryButton
                key={cat.ID}
                label={cat.Name}
                active={activeCategory?.ID === cat.ID}
                onClick={() => handleCategoryChange(cat)}
              />
            ))}
      </div>

      {/* Grid */}
      {isLoading ? (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 auto-rows-[180px]">
          {Array.from({ length: 12 }).map((_, i) => (
            <div
              key={i}
              className={`rounded-2xl bg-gradient-to-br from-teal-50 via-white to-emerald-50 animate-pulse ${
                i === 0 || i === 7 ? 'col-span-2 row-span-2' : ''
              }`}
            />
          ))}
        </div>
      ) : !data?.data.length ? (
        <div className="text-center py-20">
          <div className="text-5xl mb-4">📷</div>
          <p className="text-slate-400 text-lg">Belum ada foto/video untuk kategori ini</p>
        </div>
      ) : (
        <GalleryGrid items={data.data} onOpen={setLightbox} />
      )}

      {/* Pagination */}
      {data && data.meta.total > 12 && (
        <div className="flex justify-center items-center gap-3 mt-12">
          <button
            disabled={page <= 1}
            onClick={() => setPage(p => p - 1)}
            className="px-4 py-2 rounded-xl border-2 border-teal-100 text-sm font-medium text-teal-700 disabled:opacity-40 hover:border-teal-400 transition-colors"
          >
            ← Sebelumnya
          </button>
          <span className="px-4 py-2 text-sm text-slate-500 font-medium">Halaman {page}</span>
          <button
            disabled={page * 12 >= data.meta.total}
            onClick={() => setPage(p => p + 1)}
            className="px-4 py-2 rounded-xl border-2 border-teal-100 text-sm font-medium text-teal-700 disabled:opacity-40 hover:border-teal-400 transition-colors"
          >
            Berikutnya →
          </button>
        </div>
      )}

      {/* Lightbox */}
      {lightbox && <Lightbox item={lightbox} onClose={() => setLightbox(null)} />}
    </div>
  )
}

// ── Category Button ───────────────────────────────────────────────────────────
function CategoryButton({
  label, active, onClick,
}: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`px-4 py-2 rounded-xl text-sm font-semibold transition-all ${
        active
          ? 'bg-teal-600 text-white shadow-sm'
          : 'bg-white border-2 border-teal-100 text-slate-600 hover:border-teal-400 hover:text-teal-700'
      }`}
    >
      {label}
    </button>
  )
}

// ── Gallery Grid ──────────────────────────────────────────────────────────────
function GalleryGrid({ items, onOpen }: { items: GalleryItem[]; onOpen: (item: GalleryItem) => void }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 auto-rows-[180px]">
      {items.map((item, i) => {
        const isLarge = i === 0 || i === 7
        const thumb   = item.ThumbnailURL || item.URL
        return (
          <div
            key={item.ID}
            onClick={() => onOpen(item)}
            className={`relative rounded-2xl overflow-hidden cursor-pointer group bg-gradient-to-br from-teal-100 to-emerald-200 ${
              isLarge ? 'col-span-2 row-span-2' : ''
            }`}
          >
            {thumb ? (
              <img
                src={thumb}
                alt={item.Title}
                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                loading="lazy"
              />
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-center gap-2">
                <span className="text-4xl">{item.Type === 'video' ? '🎥' : '🌿'}</span>
                <span className="text-xs text-teal-600 font-medium px-2 text-center">{item.Title}</span>
              </div>
            )}

            {/* Overlay */}
            <div className="absolute inset-0 bg-gradient-to-t from-teal-900/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <div className="absolute bottom-0 left-0 right-0 p-3">
                <p className="text-white text-xs font-semibold line-clamp-2">{item.Title}</p>
                {item.Caption && (
                  <p className="text-teal-200 text-[11px] mt-0.5 line-clamp-1">{item.Caption}</p>
                )}
              </div>
            </div>

            {/* Badges */}
            <div className="absolute top-2 left-2 flex gap-1 flex-wrap">
              {item.Type === 'video' && (
                <span className="bg-amber-400 text-amber-900 text-[10px] font-bold px-2 py-0.5 rounded-full">
                  ▶ VIDEO
                </span>
              )}
            </div>
            {item.Category && (
              <div className="absolute top-2 right-2">
                <span className="bg-teal-600/80 backdrop-blur-sm text-white text-[10px] font-semibold px-2 py-0.5 rounded-full">
                  {item.Category.Name}
                </span>
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}

// ── Lightbox ──────────────────────────────────────────────────────────────────
function Lightbox({ item, onClose }: { item: GalleryItem; onClose: () => void }) {
  return (
    <div
      className="fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div
        className="relative max-w-4xl w-full bg-teal-950 rounded-2xl overflow-hidden shadow-2xl"
        onClick={e => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          className="absolute top-4 right-4 z-10 w-9 h-9 rounded-full bg-white/10 hover:bg-white/20 text-white flex items-center justify-center transition-colors text-lg"
        >
          ✕
        </button>

        {/* Media */}
        <div className="aspect-video bg-teal-900 flex items-center justify-center">
          {item.Type === 'video' ? (
            item.URL ? (
              <video src={item.URL} controls autoPlay className="w-full h-full object-contain" />
            ) : (
              <div className="text-center text-teal-400">
                <div className="text-6xl mb-3">🎥</div>
                <p className="text-sm">Video tidak tersedia</p>
              </div>
            )
          ) : item.URL ? (
            <img src={item.URL} alt={item.Title} className="w-full h-full object-contain" />
          ) : item.ThumbnailURL ? (
            <img src={item.ThumbnailURL} alt={item.Title} className="w-full h-full object-contain" />
          ) : (
            <div className="text-center text-teal-400">
              <div className="text-6xl mb-3">🌿</div>
              <p className="text-sm">Gambar tidak tersedia</p>
            </div>
          )}
        </div>

        {/* Info */}
        <div className="p-5">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h3 className="font-bold text-white text-lg leading-snug">{item.Title}</h3>
              {item.Caption && (
                <p className="text-teal-300 text-sm mt-1 leading-relaxed">{item.Caption}</p>
              )}
            </div>
            <div className="flex flex-col items-end gap-1.5 shrink-0">
              {item.Category && (
                <span className="text-xs bg-teal-700 text-teal-200 px-2.5 py-0.5 rounded-full font-medium">
                  {item.Category.Name}
                </span>
              )}
              {item.Type === 'video' && (
                <span className="text-xs bg-amber-500/20 text-amber-300 px-2.5 py-0.5 rounded-full font-semibold">
                  ▶ Video
                </span>
              )}
            </div>
          </div>
          <p className="text-teal-500 text-xs mt-3">
            {new Date(item.CreatedAt).toLocaleDateString('id-ID', {
              day: 'numeric', month: 'long', year: 'numeric',
            })}
          </p>
        </div>
      </div>
    </div>
  )
}