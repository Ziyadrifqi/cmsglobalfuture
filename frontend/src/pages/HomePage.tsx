import { useState, useEffect, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { portalApi, News, Volunteer, Banner } from '../api'
import { useSettings } from '../hooks/useSettings'
import { getText, getJSON } from '../lib/settings'
import { imgUrl } from '../lib/imageUrl'

const TYPE_LABEL: Record<string, string> = {
  regular: 'Reguler', event: 'Event', remote: 'Remote', training: 'Pelatihan',
}

interface StatItem    { icon: string; value: string; label: string }
interface ProgramItem { icon: string; title: string; desc: string }
interface ImpactItem  { icon: string; num: string; unit: string; desc: string }

const FALLBACK_STATS: StatItem[] = [
  { icon: '🌿', value: '12K+', label: 'Relawan Aktif' },
  { icon: '🏕️', value: '340+', label: 'Kegiatan' },
  { icon: '🗺️', value: '28',   label: 'Kota' },
  { icon: '🌳', value: '850T', label: 'Pohon Ditanam' },
]
const FALLBACK_PROGRAMS: ProgramItem[] = [
  { icon: '🌊', title: 'Bersih Pantai & Sungai',  desc: 'Aksi bersih rutin di pesisir dan daerah aliran sungai.' },
  { icon: '🌱', title: 'Penanaman Mangrove',       desc: 'Restorasi ekosistem mangrove di wilayah pesisir.' },
  { icon: '📚', title: 'Edukasi Lingkungan',       desc: 'Program edukasi ke sekolah dan komunitas lokal.' },
]
const FALLBACK_IMPACTS: ImpactItem[] = [
  { icon: '🌳', num: '850.000', unit: 'Pohon',    desc: 'ditanam sejak 2015' },
  { icon: '🗑️', num: '120 ton', unit: 'Sampah',   desc: 'berhasil dipungut' },
  { icon: '🌿', num: '45 ha',   unit: 'Mangrove', desc: 'berhasil dipulihkan' },
]

export function HomePage() {
  const { data, isLoading } = useQuery({ queryKey: ['home'], queryFn: portalApi.home })
  const { data: settings }  = useSettings()

  const heroBadge       = getText(settings, 'hero_badge_text',       '12.000+ Relawan Aktif di 28 Kota')
  const heroTitleMain   = getText(settings, 'hero_title_main',        'Bersama Jaga Bumi')
  const heroTitleHi     = getText(settings, 'hero_title_highlight',   'untuk Generasi Depan')
  const heroSubtitle    = getText(settings, 'hero_subtitle',          'Gerakan lingkungan yang mengajak semua orang beraksi nyata — menanam, membersihkan, dan mendidik demi Indonesia yang lebih hijau.')
  const heroCtaPrimary  = getText(settings, 'hero_cta_primary_text',  'Jadi Relawan')
  const heroCtaSecondary= getText(settings, 'hero_cta_secondary_text','Tentang Kami')
  const ctaTitle        = getText(settings, 'home_cta_title',         'Siap Beraksi untuk Bumi?')
  const ctaSubtitle     = getText(settings, 'home_cta_subtitle',      'Bergabunglah bersama ribuan relawan Green Future Indonesia.')

  const stats    = getJSON<StatItem[]>(settings,    'home_stats_json',    FALLBACK_STATS)
  const programs = getJSON<ProgramItem[]>(settings, 'home_programs_json', FALLBACK_PROGRAMS)
  const impacts  = getJSON<ImpactItem[]>(settings,  'home_impacts_json',  FALLBACK_IMPACTS)

  const banners = data?.banners ?? []

  return (
    <div className="bg-[#F8FAFC]">

      {/* ── 1. HERO — banner jadi fullscreen background ────────────────── */}
      <HeroBanner
        banners={banners}
        isLoading={isLoading}
        badge={heroBadge}
        titleMain={heroTitleMain}
        titleHighlight={heroTitleHi}
        subtitle={heroSubtitle}
        ctaPrimary={heroCtaPrimary}
        ctaSecondary={heroCtaSecondary}
      />

      {/* ── 2. STATS — floating overlap dari hero ─────────────────────── */}
      <section className="max-w-5xl mx-auto px-4 -mt-10 relative z-10">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          {stats.map((s, i) => (
            <div key={i}
              className="bg-white rounded-2xl p-4 md:p-5 text-center shadow-[0_4px_24px_rgba(15,118,110,0.13)] border border-white/80 backdrop-blur">
              <div className="text-2xl md:text-3xl mb-1.5">{s.icon}</div>
              <div className="text-2xl md:text-3xl font-black text-teal-700 leading-none">{s.value}</div>
              <div className="text-[11px] md:text-xs text-slate-500 mt-1 font-medium">{s.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── 3. PROGRAM UNGGULAN ────────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 pt-20 pb-4">
        <SectionHeader
          label="Apa yang Kami Lakukan"
          title="Program Unggulan"
          subtitle="Tiga program inti yang berjalan aktif di seluruh Indonesia bersama ribuan relawan."
        />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mt-8">
          {programs.map((p, i) => (
            <div key={i}
              className="group bg-white rounded-2xl p-7 border border-teal-100/60 shadow-sm
                hover:shadow-[0_8px_32px_rgba(15,118,110,0.13)] hover:-translate-y-1.5
                transition-all duration-300 cursor-default">
              <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-teal-100 to-emerald-100
                flex items-center justify-center text-3xl mb-5
                group-hover:scale-110 transition-transform duration-300">
                {p.icon}
              </div>
              <h3 className="font-bold text-[17px] text-teal-900 mb-2.5">{p.title}</h3>
              <p className="text-slate-500 text-sm leading-relaxed">{p.desc}</p>
              <div className="mt-5 w-8 h-0.5 bg-teal-400 rounded-full
                group-hover:w-16 transition-all duration-300"/>
            </div>
          ))}
        </div>
      </section>

      {/* ── 4. BERITA TERBARU ──────────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 pt-16 pb-4">
        <div className="flex items-end justify-between mb-8">
          <SectionHeader label="Kabar Terbaru" title="Berita & Artikel" />
          <Link to="/berita"
            className="hidden md:flex items-center gap-1 text-sm text-teal-600 font-semibold
              hover:text-teal-800 transition-colors">
            Lihat semua <span className="text-base">→</span>
          </Link>
        </div>
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
            {[1,2,3].map(i => <NewsCardSkeleton key={i}/>)}
          </div>
        ) : !data?.latest_news?.length ? (
          <EmptyState icon="📰" text="Belum ada berita"/>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
            {data.latest_news.slice(0,3).map(n => <NewsCard key={n.ID} news={n}/>)}
          </div>
        )}
        <Link to="/berita"
          className="md:hidden flex items-center justify-center gap-1 mt-6
            text-sm text-teal-600 font-semibold">
          Lihat semua berita →
        </Link>
      </section>

      {/* ── 5. DAMPAK NYATA ────────────────────────────────────────────── */}
      <section className="mt-16 py-20 px-4 bg-gradient-to-br from-[#0D3B36] via-[#134E4A] to-[#166534]
        relative overflow-hidden">
        <div className="absolute inset-0 opacity-10"
          style={{backgroundImage:'radial-gradient(circle at 20% 50%, #34d399 0%, transparent 50%), radial-gradient(circle at 80% 20%, #6ee7b7 0%, transparent 40%)'}}/>
        <div className="max-w-6xl mx-auto relative">
          <SectionHeader
            label="Jejak Kami"
            title="Dampak Nyata"
            subtitle="Angka-angka ini bukan sekadar statistik — ini adalah hasil kerja nyata ribuan relawan."
            dark
          />
          <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mt-10">
            {impacts.map((im, i) => (
              <CounterCard key={i} item={im}/>
            ))}
          </div>
        </div>
      </section>

      {/* ── 6. REKRUTMEN RELAWAN ───────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 pt-16 pb-4">
        <div className="flex items-end justify-between mb-8">
          <SectionHeader label="Bergabung Sekarang" title="Rekrutmen Relawan"/>
          <Link to="/relawan"
            className="hidden md:flex items-center gap-1 text-sm text-teal-600 font-semibold
              hover:text-teal-800 transition-colors">
            Lihat semua <span className="text-base">→</span>
          </Link>
        </div>
        {isLoading ? (
          <div className="space-y-3">{[1,2,3].map(i=><VolSkeleton key={i}/>)}</div>
        ) : !data?.latest_volunteers?.length ? (
          <EmptyState icon="🌱" text="Tidak ada rekrutmen aktif saat ini"/>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {data.latest_volunteers.map(v=><VolunteerCard key={v.ID} volunteer={v}/>)}
          </div>
        )}
        <Link to="/relawan"
          className="md:hidden flex items-center justify-center gap-1 mt-6
            text-sm text-teal-600 font-semibold">
          Lihat semua rekrutmen →
        </Link>
      </section>

      {/* ── 7. CTA BAWAH ───────────────────────────────────────────────── */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="relative overflow-hidden rounded-3xl
          bg-gradient-to-br from-teal-600 via-teal-700 to-emerald-700
          px-8 py-16 text-center shadow-[0_20px_60px_rgba(13,62,54,0.25)]">
          {/* Decorative circles */}
          <div className="absolute -top-16 -right-16 w-56 h-56 rounded-full bg-white/5"/>
          <div className="absolute -bottom-12 -left-12 w-44 h-44 rounded-full bg-emerald-400/10"/>
          <div className="absolute top-1/2 left-1/4 w-2 h-2 rounded-full bg-amber-400/60"/>
          <div className="absolute top-1/3 right-1/4 w-1.5 h-1.5 rounded-full bg-white/40"/>

          <div className="relative">
            <div className="inline-flex items-center gap-2 bg-white/15 border border-white/20
              rounded-full px-4 py-1.5 mb-6">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-300 animate-pulse"/>
              <span className="text-xs text-white/90 font-medium tracking-widest uppercase">
                {heroBadge}
              </span>
            </div>
            <h2 className="text-3xl md:text-4xl font-black text-white tracking-tight mb-4">
              {ctaTitle}
            </h2>
            <p className="text-teal-100 mb-10 max-w-lg mx-auto text-[15px] leading-relaxed">
              {ctaSubtitle}
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
              <Link to="/relawan"
                className="px-7 py-3.5 rounded-xl bg-white text-teal-700 font-bold text-sm
                  hover:bg-teal-50 transition-all shadow-lg shadow-black/10 hover:-translate-y-0.5">
                Daftar Jadi Relawan →
              </Link>
              <Link to="/galeri"
                className="px-7 py-3.5 rounded-xl border-2 border-white/30 text-white font-bold text-sm
                  hover:bg-white/10 hover:border-white/60 transition-all">
                Lihat Galeri Kegiatan
              </Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

/* ── Hero Banner (banner jadi fullscreen hero) ─────────────────────────────── */
function HeroBanner({
  banners, isLoading, badge, titleMain, titleHighlight,
  subtitle, ctaPrimary, ctaSecondary,
}: {
  banners: Banner[]
  isLoading: boolean
  badge: string
  titleMain: string
  titleHighlight: string
  subtitle: string
  ctaPrimary: string
  ctaSecondary: string
}) {
  const [index, setIndex] = useState(0)
  const [transitioning, setTransitioning] = useState(false)
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const hasBanners = banners.length > 0

  const goTo = (i: number) => {
    if (i === index || transitioning) return
    setTransitioning(true)
    setTimeout(() => { setIndex(i); setTransitioning(false) }, 400)
  }

  useEffect(() => {
    if (!hasBanners || banners.length <= 1) return
    timerRef.current = setInterval(() => {
      setTransitioning(true)
      setTimeout(() => {
        setIndex(i => (i + 1) % banners.length)
        setTransitioning(false)
      }, 400)
    }, 5000)
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [banners.length, hasBanners])

  const active = banners[index]

  return (
    <section className="relative h-[92vh] min-h-[580px] overflow-hidden">

      {/* Background: gambar banner atau fallback gradient */}
      {hasBanners && active ? (
        <>
          <img
            key={index}
            src={imgUrl(active.ImagePath)}
            alt={active.Title}
            className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-500
              ${transitioning ? 'opacity-0' : 'opacity-100'}`}
          />
          {/* Overlay gradient supaya teks terbaca */}
          <div className="absolute inset-0 bg-gradient-to-b
            from-black/50 via-black/35 to-black/75"/>
        </>
      ) : (
        /* Fallback kalau belum upload banner */
        <div className="absolute inset-0 bg-gradient-to-br
          from-[#0D3B36] via-[#134E4A] to-[#0f5132]">
          <div className="absolute inset-0 opacity-20"
            style={{backgroundImage:'radial-gradient(circle at 70% 30%, #34d399 0%, transparent 50%)'}}/>
        </div>
      )}

      {/* Partikel dekoratif */}
      <span className="absolute top-1/3 left-[8%] w-2 h-2 rounded-full bg-amber-400/50 blur-[1px]"/>
      <span className="absolute top-2/3 right-[12%] w-1.5 h-1.5 rounded-full bg-emerald-300/60"/>
      <span className="absolute top-1/2 left-[55%] w-1 h-1 rounded-full bg-white/30"/>

      {/* Konten hero */}
      <div className="relative h-full flex flex-col items-center justify-center text-center px-4 pb-16">

        {/* Badge */}
        <div className="inline-flex items-center gap-2 bg-white/10 border border-white/25
          backdrop-blur-sm rounded-full px-4 py-1.5 mb-7
          animate-[fadeUp_0.6s_ease_both]">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"/>
          <span className="text-xs text-white/90 font-medium tracking-widest uppercase">{badge}</span>
        </div>

        {/* Judul */}
        <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-black text-white
          leading-[1.05] tracking-tight mb-5 max-w-4xl
          animate-[fadeUp_0.7s_ease_0.1s_both]">
          {titleMain}<br/>
          <span className="text-amber-400">{titleHighlight}</span>
        </h1>

        {/* Subtitle */}
        <p className="text-white/80 text-base md:text-lg leading-relaxed mb-10
          max-w-xl animate-[fadeUp_0.7s_ease_0.2s_both]">
          {subtitle}
        </p>

        {/* CTA Buttons */}
        <div className="flex flex-wrap gap-3 justify-center animate-[fadeUp_0.7s_ease_0.3s_both]">
          <Link to="/relawan"
            className="px-8 py-4 rounded-xl bg-amber-400 text-teal-900 font-bold text-[15px]
              hover:bg-amber-300 transition-all shadow-lg shadow-amber-400/30
              hover:-translate-y-0.5 active:translate-y-0">
            {ctaPrimary} →
          </Link>
          <Link to="/tentang"
            className="px-8 py-4 rounded-xl border-2 border-white/40 text-white font-bold
              text-[15px] hover:bg-white/10 hover:border-white/70 transition-all
              backdrop-blur-sm">
            {ctaSecondary}
          </Link>
        </div>

        {/* Dot navigation — hanya kalau ada >1 banner */}
        {hasBanners && banners.length > 1 && (
          <div className="absolute bottom-10 left-1/2 -translate-x-1/2 flex gap-2 items-center">
            {banners.map((_, i) => (
              <button
                key={i}
                onClick={() => goTo(i)}
                aria-label={`Banner ${i+1}`}
                className={`rounded-full transition-all duration-300 ${
                  i === index
                    ? 'w-7 h-2.5 bg-amber-400'
                    : 'w-2.5 h-2.5 bg-white/40 hover:bg-white/70'
                }`}
              />
            ))}
          </div>
        )}

        {/* Arrow kiri kanan — hanya kalau ada >1 banner */}
        {hasBanners && banners.length > 1 && (
          <>
            <button
              onClick={() => goTo((index - 1 + banners.length) % banners.length)}
              className="absolute left-4 md:left-8 top-1/2 -translate-y-1/2
                w-10 h-10 md:w-12 md:h-12 rounded-full bg-white/10 border border-white/20
                backdrop-blur-sm text-white hover:bg-white/20 transition-all
                flex items-center justify-center text-lg">
              ‹
            </button>
            <button
              onClick={() => goTo((index + 1) % banners.length)}
              className="absolute right-4 md:right-8 top-1/2 -translate-y-1/2
                w-10 h-10 md:w-12 md:h-12 rounded-full bg-white/10 border border-white/20
                backdrop-blur-sm text-white hover:bg-white/20 transition-all
                flex items-center justify-center text-lg">
              ›
            </button>
          </>
        )}
      </div>
    </section>
  )
}

/* ── Section Header ──────────────────────────────────────────────────────────── */
function SectionHeader({
  label, title, subtitle, dark = false,
}: { label: string; title: string; subtitle?: string; dark?: boolean }) {
  return (
    <div className={dark ? 'text-center mb-2' : 'mb-2'}>
      <p className={`text-[11px] font-bold tracking-widest uppercase mb-2 ${
        dark ? 'text-emerald-400' : 'text-teal-600'}`}>
        {label}
      </p>
      <h2 className={`text-3xl md:text-4xl font-black tracking-tight ${
        dark ? 'text-white' : 'text-teal-900'}`}>
        {title}
      </h2>
      {subtitle && (
        <p className={`mt-3 text-[15px] leading-relaxed max-w-lg ${
          dark ? 'text-teal-200 mx-auto' : 'text-slate-500'}`}>
          {subtitle}
        </p>
      )}
    </div>
  )
}

/* ── Counter Card (animasi angka saat scroll ke view) ────────────────────────── */
function CounterCard({ item }: { item: ImpactItem }) {
  const ref = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const el = ref.current
    if (!el) return
    const obs = new IntersectionObserver(
      ([e]) => { if (e.isIntersecting) { setVisible(true); obs.disconnect() } },
      { threshold: 0.3 }
    )
    obs.observe(el)
    return () => obs.disconnect()
  }, [])

  return (
    <div ref={ref}
      className={`bg-white/5 border border-white/10 rounded-2xl p-8 text-center
        backdrop-blur transition-all duration-700 ${
          visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
        }`}>
      <div className="text-5xl mb-4">{item.icon}</div>
      <div className="text-4xl md:text-5xl font-black text-amber-400 mb-1">{item.num}</div>
      <div className="text-lg font-bold text-white mb-1">{item.unit}</div>
      <div className="text-sm text-teal-300">{item.desc}</div>
    </div>
  )
}

/* ── News Card ────────────────────────────────────────────────────────────────── */
function NewsCard({ news }: { news: News }) {
  const catColor: Record<string, string> = {
    Lingkungan: 'bg-teal-100 text-teal-700',
    Program:    'bg-emerald-100 text-emerald-700',
    Kampanye:   'bg-amber-100 text-amber-800',
    Edukasi:    'bg-green-100 text-green-700',
    default:    'bg-slate-100 text-slate-600',
  }
  const cc = (news.Category?.Name && catColor[news.Category.Name]) || catColor.default

  return (
    <Link to={`/berita/${news.Slug}`}
      className="group bg-white rounded-2xl overflow-hidden border border-teal-50
        hover:shadow-[0_12px_40px_rgba(15,118,110,0.14)] hover:-translate-y-1.5
        transition-all duration-300">
      <div className="relative overflow-hidden h-48">
        {news.Thumbnail ? (
          <img src={imgUrl(news.Thumbnail)} alt={news.Title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"/>
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-teal-100 to-emerald-200
            flex items-center justify-center text-5xl">🌿</div>
        )}
        {news.Category && (
          <span className={`absolute top-3 left-3 text-xs font-semibold px-2.5 py-1
            rounded-full backdrop-blur-sm ${cc}`}>
            {news.Category.Name}
          </span>
        )}
      </div>
      <div className="p-5">
        <h3 className="font-bold text-teal-900 text-[15px] leading-snug line-clamp-2
          group-hover:text-teal-600 transition-colors mb-2">
          {news.Title}
        </h3>
        {news.Excerpt && (
          <p className="text-sm text-slate-500 line-clamp-2 leading-relaxed mb-4">
            {news.Excerpt}
          </p>
        )}
        <div className="flex items-center justify-between text-xs text-slate-400">
          <span className="flex items-center gap-1.5">
            <span className="text-slate-300">✍</span>
            {news.Author?.Name}
          </span>
          {news.PublishedAt && (
            <span>{new Date(news.PublishedAt).toLocaleDateString('id-ID',{day:'numeric',month:'short',year:'numeric'})}</span>
          )}
        </div>
      </div>
    </Link>
  )
}

/* ── Volunteer Card ────────────────────────────────────────────────────────────── */
function VolunteerCard({ volunteer }: { volunteer: Volunteer }) {
  return (
    <Link to={`/relawan/${volunteer.Slug}`}
      className="group flex items-center justify-between gap-4 bg-white rounded-2xl
        border border-teal-100 px-5 py-4
        hover:border-teal-300 hover:shadow-[0_4px_20px_rgba(15,118,110,0.10)]
        transition-all duration-200">
      <div className="min-w-0">
        <h3 className="font-bold text-teal-900 text-[15px] group-hover:text-teal-700 transition-colors">
          {volunteer.Title}
        </h3>
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
          <span className="text-xs bg-teal-50 text-teal-700 font-semibold px-2 py-0.5 rounded-full">
            {TYPE_LABEL[volunteer.Type] ?? volunteer.Type}
          </span>
        </div>
      </div>
      <span className="shrink-0 px-4 py-2 bg-teal-600 text-white text-sm font-semibold
        rounded-xl group-hover:bg-teal-700 transition-colors">
        Daftar
      </span>
    </Link>
  )
}

/* ── Skeleton loaders ──────────────────────────────────────────────────────────── */
function NewsCardSkeleton() {
  return (
    <div className="bg-white rounded-2xl overflow-hidden border border-teal-50">
      <div className="h-48 bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse"/>
      <div className="p-5 space-y-3">
        <div className="h-3 rounded-full bg-teal-50 animate-pulse w-1/4"/>
        <div className="h-4 rounded-full bg-slate-100 animate-pulse"/>
        <div className="h-4 rounded-full bg-slate-100 animate-pulse w-5/6"/>
        <div className="h-3 rounded-full bg-slate-50 animate-pulse w-1/2 mt-4"/>
      </div>
    </div>
  )
}
function VolSkeleton() {
  return <div className="h-20 rounded-2xl bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse"/>
}
function EmptyState({ icon, text }: { icon: string; text: string }) {
  return (
    <div className="text-center py-16">
      <div className="text-5xl mb-3 opacity-40">{icon}</div>
      <p className="text-slate-400">{text}</p>
    </div>
  )
}
