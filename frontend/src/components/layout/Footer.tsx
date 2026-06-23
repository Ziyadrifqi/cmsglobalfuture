import { Link } from 'react-router-dom'
import { useSettings } from '../../hooks/useSettings'
import { getText, getJSON } from '../../lib/settings'

const navLinks: [string, string][] = [
  ['/', 'Beranda'],
  ['/tentang', 'Tentang Kami'],
  ['/berita', 'Berita & Artikel'],
  ['/galeri', 'Galeri'],
  ['/relawan', 'Relawan'],
  ['/kontak', 'Kontak'],
]

interface SocialItem { icon: string; label: string; link: string }

const FALLBACK_PROGRAMS = [
  'Penanaman Mangrove',
  'Bersih Sungai & Pantai',
  'Edukasi Lingkungan Sekolah',
  'Daur Ulang Sampah Plastik',
  'Pemantauan Kualitas Udara',
]

const FALLBACK_SOCIALS: SocialItem[] = [
  { icon: '𝕏', label: 'Twitter/X', link: '#' },
  { icon: 'in', label: 'LinkedIn', link: '#' },
  { icon: 'f', label: 'Facebook', link: '#' },
  { icon: '▶', label: 'YouTube', link: '#' },
]

export function Footer() {
  const { data: settings } = useSettings()

  const description = getText(settings, 'footer_description',
    'Bersama menjaga bumi untuk generasi mendatang melalui aksi nyata, edukasi, dan kolaborasi komunitas.')
  const programs = getJSON<string[]>(settings, 'footer_programs_json', FALLBACK_PROGRAMS)
  const socials  = getJSON<SocialItem[]>(settings, 'footer_social_json', FALLBACK_SOCIALS)

  const email1 = getText(settings, 'contact_email_1', 'info@greenfuture.id')
  const email2 = getText(settings, 'contact_email_2', 'media@greenfuture.id')
  const phone  = getText(settings, 'contact_phone', '(021) 456-7890')
  const phoneNote = getText(settings, 'contact_phone_note', 'Senin–Jumat 08.00–17.00')
  const address1 = getText(settings, 'contact_address_1', 'Jl. Kemang Raya No. 45')
  const address2 = getText(settings, 'contact_address_2', 'Jakarta Selatan 12730')
  const hours1 = getText(settings, 'contact_hours_1', 'Senin – Jumat')
  const hours2 = getText(settings, 'contact_hours_2', '08.00 – 17.00 WIB')

  const siteName    = getText(settings, 'site_name', 'Green Future')
  const siteNameSub = getText(settings, 'site_name_sub', 'Indonesia')

  const contactRows: [string, string, string][] = [
    ['📧', email1, email2],
    ['📞', phone, phoneNote],
    ['📍', address1, address2],
    ['🕐', hours1, hours2],
  ]

  return (
    <footer className="bg-[#0A3330] text-teal-300 mt-20">
      <div className="max-w-6xl mx-auto px-4 py-14 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-10">

        {/* Brand */}
        <div>
          <div className="flex items-center gap-2.5 mb-4">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-teal-600 to-teal-800 flex items-center justify-center text-lg">🌿</div>
            <div className="leading-tight">
              <span className="font-extrabold text-base text-white block leading-none">{siteName}</span>
              <span className="text-[10px] text-emerald-400 font-semibold tracking-widest uppercase">{siteNameSub}</span>
            </div>
          </div>
          <p className="text-sm text-teal-400 leading-relaxed mb-5 max-w-[220px]">
            {description}
          </p>
          <div className="flex gap-2">
            {socials.map((s, i) => (
              <a
                key={i}
                href={s.link}
                title={s.label}
                className="w-8 h-8 rounded-lg bg-teal-900 hover:bg-teal-700 flex items-center justify-center cursor-pointer text-xs font-bold text-teal-300 transition-colors"
              >
                {s.icon}
              </a>
            ))}
          </div>
        </div>

        {/* Navigasi */}
        <div>
          <h4 className="text-[11px] font-bold tracking-widest text-emerald-400 uppercase mb-4">Navigasi</h4>
          <div className="space-y-1.5">
            {navLinks.map(([to, label]) => (
              <Link key={to} to={to} className="block text-sm text-teal-400 hover:text-white transition-colors">
                {label}
              </Link>
            ))}
          </div>
        </div>

        {/* Program */}
        <div>
          <h4 className="text-[11px] font-bold tracking-widest text-emerald-400 uppercase mb-4">Program Kami</h4>
          <div className="space-y-1.5">
            {programs.map(p => (
              <div key={p} className="text-sm text-teal-400">{p}</div>
            ))}
          </div>
        </div>

        {/* Kontak */}
        <div>
          <h4 className="text-[11px] font-bold tracking-widest text-emerald-400 uppercase mb-4">Kontak</h4>
          <div className="space-y-3">
            {contactRows.map(([icon, v1, v2]) => (
              <div key={v1} className="flex items-start gap-3">
                <span className="text-base mt-0.5">{icon}</span>
                <div className="text-sm text-teal-400 leading-relaxed">{v1}<br/>{v2}</div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="border-t border-teal-900 px-4 py-5 text-center text-xs text-teal-600">
        © {new Date().getFullYear()} {siteName} {siteNameSub}. Hak cipta dilindungi. 🌱 Untuk Bumi yang Lebih Hijau.
      </div>
    </footer>
  )
}
