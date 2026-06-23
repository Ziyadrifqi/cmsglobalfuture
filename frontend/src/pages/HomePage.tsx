import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { portalApi, News, Volunteer, Banner } from '../api'
import { useSettings } from '../hooks/useSettings'
import { getText, getJSON } from '../lib/settings'

const TYPE_LABEL: Record<string, string> = {
  regular:  'Reguler',
  event:    'Event',
  remote:   'Remote',
  training: 'Pelatihan',
}

interface StatItem    { value: string; label: string; icon: string }
interface ProgramItem { icon: string; title: string; desc: string }
interface ImpactItem  { icon: string; num: string; unit: string; desc: string }

const FALLBACK_STATS: StatItem[] = [
  { value: '12K+', label: 'Relawan Aktif',      icon: '🌿' },
  { value: '340+', label: 'Kegiatan Terlaksana', icon: '🏕️' },
  { value: '28',   label: 'Kota Terjangkau',     icon: '🗺️' },
  { value: '850T', label: 'Pohon Ditanam',        icon: '🌳' },
]

const FALLBACK_PROGRAMS: ProgramItem[] = [
  { icon: '🌊', title: 'Bersih Pantai & Sungai',  desc: 'Aksi bersih rutin di pesisir dan daerah aliran sungai untuk mengurangi sampah plastik di perairan Indonesia.' },
  { icon: '🌱', title: 'Penanaman Mangrove',       desc: 'Restorasi ekosistem mangrove di wilayah pesisir yang terdegradasi untuk melindungi garis pantai dan habitat satwa.' },
  { icon: '📚', title: 'Edukasi Lingkungan',       desc: 'Program edukasi lingkungan ke sekolah-sekolah dan komunitas untuk membangun kesadaran sejak dini.' },
]

const FALLBACK_IMPACTS: ImpactItem[] = [
  { num: '850.000', unit: 'Pohon',    desc: 'ditanam sejak 2015', icon: '🌳' },
  { num: '120 ton', unit: 'Sampah',   desc: 'berhasil dipungut',  icon: '🗑️' },
  { num: '45 ha',   unit: 'Mangrove', desc: 'berhasil dipulihkan', icon: '🌿' },
]

export function HomePage() {
  const { data, isLoading } = useQuery({
    queryKey: ['home'],
    queryFn:  portalApi.home,
  })
  const { data: settings } = useSettings()

  const heroBadge      = getText(settings, 'hero_badge_text', '12.000+ Relawan Aktif di 28 Kota')
  const heroTitleMain   = getText(settings, 'hero_title_main', 'Bersama Jaga Bumi')
  const heroTitleHi     = getText(settings, 'hero_title_highlight', 'untuk Generasi Depan')
  const heroSubtitle    = getText(settings, 'hero_subtitle',
    'Green Future Indonesia adalah gerakan lingkungan yang mengajak semua orang untuk beraksi nyata — menanam, membersihkan, dan mendidik demi Indonesia yang lebih hijau.')
  const heroCtaPrimary   = getText(settings, 'hero_cta_primary_text', 'Jadi Relawan')
  const heroCtaSecondary = getText(settings, 'hero_cta_secondary_text', 'Tentang Kami')

  const stats    = getJSON<StatItem[]>(settings, 'home_stats_json', FALLBACK_STATS)
  const programs = getJSON<ProgramItem[]>(settings, 'home_programs_json', FALLBACK_PROGRAMS)
  const impacts  = getJSON<ImpactItem[]>(settings, 'home_impacts_json', FALLBACK_IMPACTS)

  const ctaTitle    = getText(settings, 'home_cta_title', 'Siap Beraksi untuk Bumi?')
  const ctaSubtitle = getText(settings, 'home_cta_subtitle',
    'Bergabunglah bersama 12.000+ relawan Green Future Indonesia dan jadilah bagian dari perubahan nyata.')

  return (
    <div>
      {/* ── HERO ─────────────────────────────────────────────────────── */}
      <section className="relative overflow-hidden bg-gradient-to-br from-[#134E4A] via-[#0F766E] to-[#166534] py-24 px-4">
        <div className="absolute -top-24 -right-24 w-[420px] h-[420px] rounded-full bg-white/[0.04] pointer-events-none" />
        <div className="absolute -bottom-16 -left-16 w-[300px] h-[300px] rounded-full bg-amber-400/[0.07] pointer-events-none" />
        <span className="absolute top-1/3 left-[8%] w-2 h-2 rounded-full bg-amber-400/60" />
        <span className="absolute top-2/3 right-[12%] w-1.5 h-1.5 rounded-full bg-emerald-300/50" />

        <div className="max-w-3xl mx-auto text-center relative">
          <div className="inline-flex items-center gap-2 bg-white/10 border border-white/20 rounded-full px-4 py-1.5 mb-8">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse inline-block" />
            <span className="text-xs text-teal-200 font-medium tracking-widest uppercase">
              {heroBadge}
            </span>
          </div>
          <h1 className="text-5xl md:text-6xl font-black text-white leading-[1.1] tracking-tight mb-6">
            {heroTitleMain}<br />
            <span className="text-amber-400">{heroTitleHi}</span>
          </h1>
          <p className="text-teal-200 text-lg leading-relaxed mb-10 max-w-xl mx-auto">
            {heroSubtitle}
          </p>
          <div className="flex flex-wrap gap-3 justify-center">
            <Link
              to="/relawan"
              className="px-7 py-3.5 rounded-xl bg-amber-400 text-teal-900 font-bold text-[15px] hover:bg-amber-300 transition-all shadow-lg shadow-amber-400/30 hover:-translate-y-0.5"
            >
              {heroCtaPrimary} →
            </Link>
            <Link
              to="/tentang"
              className="px-7 py-3.5 rounded-xl border-2 border-white/40 text-white font-bold text-[15px] hover:bg-white/10 hover:border-white/70 transition-all"
            >
              {heroCtaSecondary}
            </Link>
          </div>
        </div>
      </section>

      {/* ── STATS ────────────────────────────────────────────────────── */}
      <div className="max-w-4xl mx-auto px-4 -mt-10 relative z-10">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {stats.map(s => (
            <div key={s.label} className="bg-white rounded-2xl p-5 text-center shadow-[0_4px_24px_rgba(15,118,110,0.12)] border border-teal-50">
              <div className="text-3xl mb-2">{s.icon}</div>
              <div className="text-3xl font-black text-teal-600 leading-none">{s.value}</div>
              <div className="text-xs text-slate-500 mt-1.5 font-medium">{s.label}</div>
            </div>
          ))}
        </div>
      </div>

      {/* ── BANNER ───────────────────────────────────────────────────── */}
      {!isLoading && data?.banners && data.banners.length > 0 && (
        <BannerCarousel banners={data.banners} />
      )}

      {/* ── BERITA TERBARU ───────────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 pt-20 pb-4">
        <div className="flex items-end justify-between mb-8">
          <div>
            <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-1.5">Kabar Terbaru</p>
            <h2 className="text-3xl font-black text-teal-900 tracking-tight">Berita &amp; Artikel</h2>
          </div>
          <Link to="/berita" className="text-sm text-teal-600 font-semibold hover:text-teal-800 transition-colors">
            Lihat semua →
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[1, 2, 3].map(i => (
              <div key={i} className="rounded-2xl overflow-hidden border border-teal-100">
                <div className="h-44 bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse" />
                <div className="p-4 space-y-2">
                  <div className="h-3 rounded-full bg-teal-100 animate-pulse w-1/3" />
                  <div className="h-4 rounded-full bg-slate-100 animate-pulse" />
                  <div className="h-4 rounded-full bg-slate-100 animate-pulse w-5/6" />
                </div>
              </div>
            ))}
          </div>
        ) : !data?.latest_news?.length ? (
          <p className="text-slate-400 text-sm text-center py-10">Belum ada berita</p>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {data.latest_news.slice(0, 3).map(news => <NewsCard key={news.ID} news={news} />)}
          </div>
        )}
      </section>

      {/* ── PROGRAM UNGGULAN ─────────────────────────────────────────── */}
      <section className="bg-teal-50/60 py-20 px-4 mt-16">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-12">
            <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-1.5">Apa yang Kami Lakukan</p>
            <h2 className="text-3xl font-black text-teal-900 tracking-tight mb-3">Program Unggulan</h2>
            <p className="text-slate-500 max-w-md mx-auto text-[15px]">
              Tiga program inti kami yang berjalan aktif di seluruh Indonesia bersama ribuan relawan.
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {programs.map(p => (
              <div
                key={p.title}
                className="bg-white rounded-2xl p-8 border border-teal-100 shadow-sm hover:shadow-[0_8px_32px_rgba(15,118,110,0.12)] hover:-translate-y-1 transition-all duration-200"
              >
                <div className="w-14 h-14 rounded-2xl bg-teal-100 flex items-center justify-center text-3xl mb-5">
                  {p.icon}
                </div>
                <h3 className="font-bold text-lg text-teal-900 mb-2">{p.title}</h3>
                <p className="text-slate-500 text-sm leading-relaxed">{p.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── DAMPAK NYATA ─────────────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="text-center mb-12">
          <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-1.5">Jejak Kami</p>
          <h2 className="text-3xl font-black text-teal-900 tracking-tight">Dampak Nyata</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {impacts.map(im => (
            <div key={im.unit} className="bg-gradient-to-br from-[#134E4A] to-[#0F766E] rounded-2xl p-8 text-center">
              <div className="text-5xl mb-4">{im.icon}</div>
              <div className="text-4xl font-black text-amber-400 mb-1">{im.num}</div>
              <div className="text-lg font-bold text-white mb-1">{im.unit}</div>
              <div className="text-sm text-teal-200">{im.desc}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── REKRUTMEN RELAWAN ────────────────────────────────────────── */}
      <section className="bg-teal-50/60 py-16 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-end justify-between mb-8">
            <div>
              <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-1.5">Bergabung Sekarang</p>
              <h2 className="text-3xl font-black text-teal-900 tracking-tight">Rekrutmen Relawan</h2>
            </div>
            <Link to="/relawan" className="text-sm text-teal-600 font-semibold hover:text-teal-800 transition-colors">
              Lihat semua →
            </Link>
          </div>

          {isLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map(i => (
                <div key={i} className="h-20 rounded-2xl bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse" />
              ))}
            </div>
          ) : !data?.latest_volunteers?.length ? (
            <p className="text-slate-400 text-sm text-center py-10">Tidak ada rekrutmen aktif saat ini</p>
          ) : (
            <div className="space-y-3">
              {data.latest_volunteers.map(v => <VolunteerCard key={v.ID} volunteer={v} />)}
            </div>
          )}
        </div>
      </section>

      {/* ── CTA ──────────────────────────────────────────────────────── */}
      <div className="max-w-6xl mx-auto px-4 py-20">
        <div className="relative overflow-hidden bg-gradient-to-br from-[#134E4A] to-[#166534] rounded-3xl px-8 py-14 text-center">
          <div className="absolute -top-10 -right-10 w-48 h-48 rounded-full bg-amber-400/10 pointer-events-none" />
          <div className="relative">
            <div className="text-5xl mb-4">🌍</div>
            <h2 className="text-3xl font-black text-white mb-4 tracking-tight">{ctaTitle}</h2>
            <p className="text-teal-200 mb-8 max-w-md mx-auto text-[15px] leading-relaxed">
              {ctaSubtitle}
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
              <Link
                to="/relawan"
                className="px-6 py-3 rounded-xl bg-amber-400 text-teal-900 font-bold text-sm hover:bg-amber-300 transition-colors shadow-lg shadow-amber-400/20"
              >
                Daftar Jadi Relawan
              </Link>
              <Link
                to="/galeri"
                className="px-6 py-3 rounded-xl border-2 border-white/30 text-white font-bold text-sm hover:bg-white/10 transition-colors"
              >
                Lihat Galeri Kegiatan
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

// ── Sub-components ────────────────────────────────────────────────────────────

function BannerCarousel({ banners }: { banners: Banner[] }) {
  const [index, setIndex] = useState(0)

  useEffect(() => {
    if (banners.length <= 1) return
    const timer = setInterval(() => setIndex(i => (i + 1) % banners.length), 5000)
    return () => clearInterval(timer)
  }, [banners.length])

  const active = banners[index]
  if (!active) return null

  const slide = (
    <div className="relative rounded-2xl overflow-hidden h-48 md:h-64 bg-gradient-to-br from-teal-100 to-emerald-200 group">
      <img src={active.ImagePath} alt={active.Title} className="w-full h-full object-cover" />
      <div className="absolute inset-0 bg-gradient-to-t from-black/55 via-black/0 to-transparent flex items-end p-5 md:p-6">
        <p className="text-white font-bold text-lg md:text-xl drop-shadow">{active.Title}</p>
      </div>
    </div>
  )

  return (
    <section className="max-w-6xl mx-auto px-4 pt-12">
      {active.LinkURL ? (
        <a href={active.LinkURL} target="_blank" rel="noreferrer noopener">{slide}</a>
      ) : slide}

      {banners.length > 1 && (
        <div className="flex justify-center gap-1.5 mt-3">
          {banners.map((b, i) => (
            <button
              key={b.ID}
              onClick={() => setIndex(i)}
              aria-label={`Tampilkan banner ${i + 1}`}
              className={`h-2 rounded-full transition-all ${i === index ? 'w-6 bg-teal-600' : 'w-2 bg-teal-200 hover:bg-teal-300'}`}
            />
          ))}
        </div>
      )}
    </section>
  )
}

function NewsCard({ news }: { news: News }) {
  const catColor: Record<string, string> = {
    Lingkungan: 'bg-teal-100 text-teal-700',
    Program:    'bg-emerald-100 text-emerald-700',
    Kampanye:   'bg-amber-100 text-amber-800',
    Edukasi:    'bg-green-100 text-green-700',
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
        <div className="w-full h-44 bg-gradient-to-br from-teal-100 to-emerald-200 flex items-center justify-center text-4xl relative">
          🌿
          {news.Category && (
            <span className={`absolute top-3 right-3 text-xs font-semibold px-2.5 py-0.5 rounded-full ${cc}`}>
              {news.Category.Name}
            </span>
          )}
        </div>
      )}
      <div className="p-5">
        {news.Thumbnail && news.Category && (
          <span className={`text-xs font-semibold px-2.5 py-0.5 rounded-full ${cc}`}>
            {news.Category.Name}
          </span>
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

function VolunteerCard({ volunteer }: { volunteer: Volunteer }) {
  return (
    <Link
      to={`/relawan/${volunteer.Slug}`}
      className="flex items-center justify-between gap-4 bg-white rounded-2xl border-2 border-teal-100 px-5 py-4 hover:border-teal-400 hover:shadow-[0_4px_20px_rgba(15,118,110,0.10)] transition-all duration-200"
    >
      <div className="min-w-0">
        <div className="font-bold text-teal-900">{volunteer.Title}</div>
        <div className="flex flex-wrap gap-1.5 mt-2">
          {volunteer.Division && (
            <span className="text-xs bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full">
              {volunteer.Division}
            </span>
          )}
          {volunteer.Location && (
            <span className="text-xs bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full">
              📍 {volunteer.Location}
            </span>
          )}
          <span className="text-xs bg-teal-100 text-teal-700 font-semibold px-2 py-0.5 rounded-full">
            {TYPE_LABEL[volunteer.Type] ?? volunteer.Type}
          </span>
        </div>
      </div>
      <span className="shrink-0 px-4 py-2 bg-teal-600 text-white text-sm font-semibold rounded-xl">Daftar</span>
    </Link>
  )
}
