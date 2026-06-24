import { useState } from 'react'
import { useSettings } from '../hooks/useSettings'
import { getText, getJSON } from '../lib/settings'

interface SocialItem { icon: string; platform: string; handle: string }

const FALLBACK_SOCIALS: SocialItem[] = [
  { icon: '🌿', platform: 'Instagram', handle: '@greenfuture.id' },
  { icon: '▶',  platform: 'YouTube',   handle: 'Green Future ID' },
  { icon: '𝕏',  platform: 'Twitter/X', handle: '@GreenFutureID' },
  { icon: 'in', platform: 'LinkedIn',  handle: 'Green Future Indonesia' },
]

export function ContactPage() {
  const [sending, setSending] = useState(false)
  const [sent, setSent]       = useState(false)
  const { data: settings } = useSettings()

  // Semua info kontak dari settings — tidak bergantung pada Page "contact"
  const email1    = getText(settings, 'contact_email_1', 'info@greenfuture.id')
  const email2    = getText(settings, 'contact_email_2', 'media@greenfuture.id')
  const phone     = getText(settings, 'contact_phone', '(021) 456-7890')
  const phoneNote = getText(settings, 'contact_phone_note', 'Senin–Jumat, 08.00–17.00')
  const address1  = getText(settings, 'contact_address_1', 'Jl. Kemang Raya No. 45')
  const address2  = getText(settings, 'contact_address_2', 'Jakarta Selatan 12730')
  const hours1    = getText(settings, 'contact_hours_1', 'Senin – Jumat: 08.00 – 17.00')
  const hours2    = getText(settings, 'contact_hours_2', 'Sabtu: 09.00 – 13.00 WIB')
  const mapLabel  = getText(settings, 'contact_map_label', 'Kemang, Jakarta Selatan')
  const socials   = getJSON<SocialItem[]>(settings, 'contact_social_json', FALLBACK_SOCIALS)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setSending(true)
    setTimeout(() => { setSending(false); setSent(true) }, 1600)
  }

  const inputCls = 'w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors bg-white text-slate-700 placeholder:text-slate-400'

  const infoRows: [string, string, string, string][] = [
    ['📧', 'Email',           email1, email2],
    ['📞', 'Telepon',         phone,  phoneNote],
    ['📍', 'Alamat',          address1, address2],
    ['🕐', 'Jam Operasional', hours1, hours2],
  ]

  return (
    <div className="max-w-5xl mx-auto px-4 py-12">

      <div className="mb-12">
        <p className="text-[11px] font-bold tracking-widest text-teal-600 uppercase mb-2">Hubungi Kami</p>
        <h1 className="text-4xl font-black text-teal-900 tracking-tight mb-2">Kontak</h1>
        <p className="text-slate-500 text-[15px]">Ada pertanyaan tentang program atau ingin berkolaborasi? Kami siap mendengar!</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">

        {/* Info */}
        <div className="space-y-5">
          <div className="bg-white rounded-2xl border border-teal-100 shadow-sm p-7">
            <h2 className="font-bold text-teal-900 text-[17px] mb-6">Informasi Kontak</h2>
            <div className="space-y-5">
              {infoRows.map(([icon, label, v1, v2]) => (
                <div key={label} className="flex items-start gap-4">
                  <div className="w-10 h-10 rounded-xl bg-teal-50 border border-teal-100 flex items-center justify-center text-[18px] shrink-0">{icon}</div>
                  <div>
                    <div className="font-semibold text-sm text-teal-900 mb-0.5">{label}</div>
                    <div className="text-sm text-slate-500 leading-relaxed">{v1}<br/>{v2}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Sosial media */}
          <div className="bg-white rounded-2xl border border-teal-100 p-6">
            <h3 className="font-bold text-teal-900 text-sm mb-4">Ikuti Kami</h3>
            <div className="grid grid-cols-2 gap-3">
              {socials.map(s => (
                <div key={s.platform} className="flex items-center gap-2.5 p-2.5 rounded-xl bg-teal-50 border border-teal-100">
                  <span className="text-xl">{s.icon}</span>
                  <div>
                    <div className="text-xs font-bold text-teal-900">{s.platform}</div>
                    <div className="text-[11px] text-teal-600">{s.handle}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Peta placeholder */}
          <div className="bg-gradient-to-br from-teal-100 to-emerald-200 rounded-2xl h-40 flex flex-col items-center justify-center gap-2 text-teal-700">
            <span className="text-4xl">🗺️</span>
            <span className="text-sm font-semibold">{mapLabel}</span>
          </div>
        </div>

        {/* Form */}
        <div className="bg-white rounded-2xl border border-teal-100 shadow-sm p-7">
          <h2 className="font-bold text-teal-900 text-[17px] mb-6">Kirim Pesan</h2>

          {sent ? (
            <div className="flex flex-col items-center justify-center py-10 text-center gap-3">
              <div className="w-16 h-16 rounded-full bg-emerald-100 flex items-center justify-center text-4xl animate-bounce">🌿</div>
              <p className="font-bold text-emerald-800 text-lg">Pesan Terkirim!</p>
              <p className="text-sm text-slate-500">Tim kami akan membalas dalam 1–2 hari kerja.</p>
              <button onClick={() => setSent(false)} className="mt-3 px-5 py-2 rounded-xl border-2 border-teal-200 text-sm font-semibold text-teal-700 hover:bg-teal-50 transition-colors">
                Kirim Pesan Lagi
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-1.5">Nama *</label>
                  <input type="text" required placeholder="Nama lengkap" className={inputCls}/>
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-1.5">Email *</label>
                  <input type="email" required placeholder="email@kamu.com" className={inputCls}/>
                </div>
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Topik</label>
                <select className={inputCls}>
                  <option value="">Pilih topik pesan</option>
                  <option>Pertanyaan Program</option>
                  <option>Kemitraan &amp; Sponsorship</option>
                  <option>Donasi</option>
                  <option>Media &amp; Pers</option>
                  <option>Lainnya</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Pesan *</label>
                <textarea required rows={6} placeholder="Tulis pesan Anda di sini..." className={`${inputCls} resize-none`}/>
              </div>
              <button type="submit" disabled={sending}
                className="w-full py-3 rounded-xl bg-gradient-to-r from-teal-600 to-emerald-600 text-white font-bold text-[15px] hover:from-teal-700 hover:to-emerald-700 transition-all disabled:opacity-60 shadow-md flex items-center justify-center gap-2">
                {sending
                  ? <><span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"/>Mengirim...</>
                  : 'Kirim Pesan 🌿'}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
