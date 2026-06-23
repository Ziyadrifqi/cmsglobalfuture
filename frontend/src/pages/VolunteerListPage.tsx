import { useQuery } from '@tanstack/react-query'
import { Link, useSearchParams } from 'react-router-dom'
import { portalApi, Volunteer } from '../api'
import { useSettings } from '../hooks/useSettings'
import { getJSON } from '../lib/settings'

const TYPE_LABEL: Record<string, string> = {
  regular: 'Reguler',
  event:   'Event',
  remote:  'Remote',
  training: 'Pelatihan',
}

interface WhyJoinItem { icon: string; title: string; desc: string }

const FALLBACK_WHY_JOIN: WhyJoinItem[] = [
  { icon: '🌿', title: 'Dampak Nyata',    desc: 'Setiap aksi Anda langsung terasa — pohon tertanam, pantai lebih bersih, alam terjaga.' },
  { icon: '🤝', title: 'Komunitas Solid', desc: 'Bergabung dengan 12.000+ relawan yang peduli dan saling mendukung.' },
  { icon: '📜', title: 'Sertifikat Resmi', desc: 'Dapatkan sertifikat relawan yang diakui dan bisa dicantumkan di CV.' },
  { icon: '🎓', title: 'Pelatihan Gratis', desc: 'Akses workshop lingkungan, manajemen komunitas, dan keterampilan lapangan.' },
]

export function VolunteerListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const page = Number(searchParams.get('page') ?? 1)

  const { data, isLoading } = useQuery({
    queryKey: ['volunteers-list', page],
    queryFn:  () => portalApi.volunteerList(page),
  })
  const { data: settings } = useSettings()
  const whyJoin = getJSON<WhyJoinItem[]>(settings, 'volunteer_why_join_json', FALLBACK_WHY_JOIN)

  return (
    <div className="max-w-4xl mx-auto px-4 py-12">

      <div className="mb-10">
        <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-2">Bergabung Bersama Kami</p>
        <h1 className="text-4xl font-black text-teal-900 tracking-tight mb-2">Rekrutmen Relawan</h1>
        <p className="text-slate-500 text-[15px]">Jadilah bagian dari gerakan nyata menjaga lingkungan Indonesia.</p>
      </div>

      {/* Why join */}
      <div className="bg-teal-50 rounded-2xl p-6 mb-10 grid grid-cols-2 md:grid-cols-4 gap-5">
        {whyJoin.map(w => (
          <div key={w.title}>
            <div className="text-2xl mb-1.5">{w.icon}</div>
            <div className="font-bold text-sm text-teal-900 mb-1">{w.title}</div>
            <div className="text-xs text-slate-500 leading-relaxed">{w.desc}</div>
          </div>
        ))}
      </div>

      {/* List */}
      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="h-24 rounded-2xl bg-gradient-to-r from-teal-50 via-white to-teal-50 animate-pulse" />
          ))}
        </div>
      ) : !data?.data.length ? (
        <div className="text-center py-20">
          <div className="text-5xl mb-4">🌱</div>
          <p className="text-slate-400 text-lg">Belum ada rekrutmen aktif saat ini</p>
          <p className="text-sm text-slate-400 mt-1">Pantau terus halaman ini ya!</p>
        </div>
      ) : (
        <div className="space-y-3">
          {data.data.map(v => <VolunteerItem key={v.ID} volunteer={v} />)}
        </div>
      )}

      {/* Pagination */}
      {data && data.meta.total > 10 && (
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
            disabled={page * 10 >= (data?.meta.total ?? 0)}
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

function VolunteerItem({ volunteer }: { volunteer: Volunteer }) {
  return (
    <Link
      to={`/relawan/${volunteer.Slug}`}
      className="flex items-center justify-between gap-4 bg-white rounded-2xl border-2 border-teal-100 px-5 py-4 hover:border-teal-400 hover:shadow-[0_4px_20px_rgba(15,118,110,0.10)] transition-all duration-200"
    >
      <div className="flex-1 min-w-0">
        <h3 className="font-bold text-teal-900 text-[15px]">{volunteer.Title}</h3>
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
      <span className="shrink-0 px-4 py-2 bg-teal-600 text-white text-sm font-semibold rounded-xl hover:bg-teal-700 transition-colors">
        Daftar
      </span>
    </Link>
  )
}
