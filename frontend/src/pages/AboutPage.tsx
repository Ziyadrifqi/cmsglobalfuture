import { useQuery } from '@tanstack/react-query'
import { portalApi } from '../api'
import { useSettings } from '../hooks/useSettings'
import { getJSON } from '../lib/settings'

interface TeamMember  { name: string; role: string; icon: string }
interface ValueItem   { icon: string; title: string; desc: string }
interface Milestone   { year: string; label: string }

const FALLBACK_TEAM: TeamMember[] = [
  { name: 'Rizal Anwar, M.Env', role: 'Ketua Umum',               icon: '👨‍💼' },
  { name: 'Nadya Putri, S.Hut', role: 'Kepala Program Lapangan',   icon: '👩‍🌾' },
  { name: 'Bayu Santoso',       role: 'Kepala Komunikasi & Media', icon: '👨‍💻' },
  { name: 'Sari Dewi, M.Si',   role: 'Manajer Relawan',           icon: '👩‍💼' },
]

const FALLBACK_VALUES: ValueItem[] = [
  { icon: '🌿', title: 'Kepedulian',  desc: 'Kami percaya setiap tindakan kecil untuk lingkungan memberi dampak besar bagi bumi.' },
  { icon: '🤝', title: 'Kolaborasi',  desc: 'Perubahan nyata hanya bisa dicapai bersama — lintas komunitas, daerah, dan latar belakang.' },
  { icon: '🔬', title: 'Berbasis Data', desc: 'Setiap program dirancang dan dievaluasi berdasarkan data lingkungan yang valid dan terukur.' },
  { icon: '🌍', title: 'Inklusif',    desc: 'Siapapun bisa berkontribusi — tanpa batasan usia, profesi, atau lokasi.' },
]

const FALLBACK_MILESTONES: Milestone[] = [
  { year: '2015', label: 'Didirikan di Jakarta oleh 12 aktivis lingkungan muda' },
  { year: '2017', label: 'Program Penanaman Mangrove perdana di Kepulauan Seribu' },
  { year: '2019', label: 'Ekspansi ke 10 kota, relawan tembus 3.000 orang' },
  { year: '2021', label: 'Kemitraan resmi dengan KLHK dan 8 Pemda' },
  { year: '2023', label: 'Raih penghargaan Lingkungan Nasional dari UNEP Indonesia' },
  { year: '2025', label: '12.000+ relawan aktif, 28 kota, 850.000 pohon ditanam' },
]

const FALLBACK_ABOUT = `Green Future Indonesia adalah organisasi lingkungan nirlaba yang didirikan pada tahun 2015 oleh sekumpulan aktivis muda yang prihatin terhadap kondisi lingkungan Indonesia.

Kami percaya bahwa perubahan lingkungan yang nyata hanya bisa terwujud melalui kolaborasi antara masyarakat, pemerintah, dan sektor swasta. Karena itu, kami membangun ekosistem relawan yang kuat — dari pelajar, mahasiswa, profesional, hingga komunitas lokal.

Hingga 2025, Green Future Indonesia telah menjangkau lebih dari 28 kota dengan 12.000+ relawan aktif, dan telah berhasil menanam 850.000 pohon, memungut 120 ton sampah, serta memulihkan 45 hektar ekosistem mangrove.`

export function AboutPage() {
  const { data: about,  isLoading: loadAbout }  = useQuery({ queryKey: ['page','about'],           queryFn: () => portalApi.page('about') })
  const { data: vismis, isLoading: loadVismis } = useQuery({ queryKey: ['page','vision-mission'],  queryFn: () => portalApi.page('vision-mission') })
  const { data: settings } = useSettings()

  const team       = getJSON<TeamMember[]>(settings, 'about_team_json', FALLBACK_TEAM)
  const values     = getJSON<ValueItem[]>(settings, 'about_values_json', FALLBACK_VALUES)
  const milestones = getJSON<Milestone[]>(settings, 'about_milestones_json', FALLBACK_MILESTONES)

  const SkeletonBlock = () => (
    <div className="space-y-3">
      {[1,2,3,4].map(i => <div key={i} className={`h-4 rounded-full bg-teal-50 animate-pulse ${i===4 ? 'w-3/4' : ''}`} />)}
    </div>
  )

  return (
    <div className="max-w-4xl mx-auto px-4 py-12">

      <div className="mb-12">
        <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-2">Siapa Kami</p>
        <h1 className="text-4xl font-black text-teal-900 tracking-tight mb-2">Tentang Kami</h1>
        <p className="text-slate-500 text-[15px]">Mengenal lebih jauh Green Future Indonesia dan perjalanan kami menjaga bumi.</p>
      </div>

      {/* About */}
      <div className="bg-white rounded-2xl border border-teal-100 shadow-sm p-8 mb-6">
        <h2 className="font-black text-xl text-teal-900 mb-5">{about?.Title ?? 'Tentang Green Future Indonesia'}</h2>
        {loadAbout ? <SkeletonBlock /> : (
          <div className="text-slate-600 leading-[1.9] text-[15px] whitespace-pre-wrap">{about?.Content ?? FALLBACK_ABOUT}</div>
        )}
      </div>

      {/* Visi & Misi */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-10">
        <div className="bg-gradient-to-br from-[#134E4A] to-[#166534] rounded-2xl p-8 text-white">
          <div className="text-4xl mb-4">🔭</div>
          <h3 className="font-black text-xl text-emerald-300 mb-3">Visi</h3>
          {loadVismis ? (
            <div className="space-y-2">{[1,2,3].map(i => <div key={i} className="h-3 rounded-full bg-teal-700/50 animate-pulse" />)}</div>
          ) : (
            <p className="text-teal-100 leading-relaxed text-sm">
              Mewujudkan Indonesia yang hijau, bersih, dan lestari — di mana manusia dan alam hidup berdampingan secara harmonis untuk generasi yang akan datang.
            </p>
          )}
        </div>
        <div className="bg-teal-50 rounded-2xl p-8 border-2 border-teal-200">
          <div className="text-4xl mb-4">🎯</div>
          <h3 className="font-black text-xl text-teal-900 mb-3">Misi</h3>
          {loadVismis ? (
            <div className="space-y-2">{[1,2,3,4,5].map(i => <div key={i} className="h-3 rounded-full bg-teal-100 animate-pulse" />)}</div>
          ) : (
            <ul className="space-y-2 text-sm text-slate-600 leading-relaxed">
              {['Menggerakkan aksi restorasi lingkungan berbasis komunitas','Membangun jaringan relawan lingkungan di seluruh Indonesia','Mendorong kebijakan lingkungan yang berpihak pada alam','Mengedukasi generasi muda tentang krisis iklim','Berkolaborasi dengan pemerintah dan sektor swasta'].map(m => (
                <li key={m} className="flex items-start gap-2">
                  <span className="text-teal-500 font-bold mt-0.5">✓</span>{m}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      {/* Milestones */}
      <h2 className="font-black text-xl text-teal-900 mb-6">Perjalanan Kami</h2>
      <div className="bg-white rounded-2xl border border-teal-100 p-8 mb-10">
        <div className="space-y-0">
          {milestones.map((m, i) => (
            <div key={m.year + i} className="flex gap-5 pb-6 last:pb-0">
              <div className="flex flex-col items-center">
                <div className="w-10 h-10 rounded-full bg-teal-600 text-white flex items-center justify-center text-xs font-black shrink-0">{m.year.slice(2)}</div>
                {i < milestones.length - 1 && <div className="w-0.5 flex-1 bg-teal-100 mt-2" />}
              </div>
              <div className="pt-2 pb-2">
                <span className="text-xs font-bold text-teal-600 block mb-0.5">{m.year}</span>
                <p className="text-sm text-slate-600 leading-relaxed">{m.label}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Values */}
      <h2 className="font-black text-xl text-teal-900 mb-6">Nilai-Nilai Kami</h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
        {values.map(v => (
          <div key={v.title} className="bg-white rounded-2xl p-5 border border-teal-100 text-center hover:shadow-md hover:-translate-y-0.5 transition-all">
            <div className="text-3xl mb-3">{v.icon}</div>
            <div className="font-bold text-sm text-teal-900 mb-1.5">{v.title}</div>
            <div className="text-xs text-slate-500 leading-relaxed">{v.desc}</div>
          </div>
        ))}
      </div>

      {/* Team */}
      <h2 className="font-black text-xl text-teal-900 mb-6">Tim Pengurus</h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {team.map(m => (
          <div key={m.name} className="bg-white rounded-2xl p-6 border border-teal-100 text-center hover:shadow-md hover:-translate-y-0.5 transition-all">
            <div className="w-16 h-16 rounded-full bg-gradient-to-br from-teal-100 to-emerald-200 flex items-center justify-center text-3xl mx-auto mb-4">{m.icon}</div>
            <div className="font-bold text-sm text-teal-900 mb-1">{m.name}</div>
            <div className="text-xs text-teal-600 font-medium">{m.role}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
