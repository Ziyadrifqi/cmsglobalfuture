import { useState, useRef } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { useParams, Link } from 'react-router-dom'
import { portalApi } from '../api'

const TYPE_LABEL: Record<string, string> = {
  regular:  'Reguler',
  event:    'Event',
  remote:   'Remote',
  training: 'Pelatihan',
}

export function VolunteerDetailPage() {
  const { slug }                          = useParams<{ slug: string }>()
  const [showForm, setShowForm]           = useState(false)
  const [success, setSuccess]             = useState(false)
  const formRef                           = useRef<HTMLFormElement>(null)

  const { data: volunteer, isLoading, isError } = useQuery({
    queryKey: ['volunteer-detail', slug],
    queryFn:  () => portalApi.volunteerDetail(slug!),
    enabled:  !!slug,
  })

  const applyMut = useMutation({
    mutationFn: (fd: FormData) => portalApi.apply(slug!, fd),
    onSuccess:  () => { setSuccess(true); setShowForm(false); formRef.current?.reset() },
  })

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    applyMut.mutate(new FormData(e.currentTarget))
  }

  if (isLoading) return (
    <div className="max-w-3xl mx-auto px-4 py-12 space-y-4 animate-pulse">
      <div className="flex gap-2">
        <div className="h-3 w-16 rounded-full bg-teal-100"/>
        <div className="h-3 w-3 rounded-full bg-teal-100"/>
        <div className="h-3 w-24 rounded-full bg-teal-100"/>
      </div>
      <div className="h-10 rounded-xl bg-slate-100 w-2/3"/>
      <div className="h-36 rounded-2xl bg-slate-50"/>
      <div className="h-48 rounded-2xl bg-slate-50"/>
      <div className="h-12 rounded-xl bg-teal-100"/>
    </div>
  )

  if (isError || !volunteer) return (
    <div className="max-w-3xl mx-auto px-4 py-24 text-center">
      <div className="text-6xl mb-5">🌱</div>
      <h1 className="text-xl font-bold text-teal-900 mb-3">Rekrutmen tidak ditemukan</h1>
      <Link to="/relawan" className="text-teal-600 hover:underline text-sm font-medium">
        ← Kembali ke rekrutmen
      </Link>
    </div>
  )

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">

      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-sm text-slate-400 mb-8">
        <Link to="/" className="hover:text-teal-600 transition-colors">Beranda</Link>
        <span>/</span>
        <Link to="/relawan" className="hover:text-teal-600 transition-colors">Relawan</Link>
        <span>/</span>
        <span className="text-slate-600 line-clamp-1 max-w-[200px]">{volunteer.Title}</span>
      </nav>

      {/* Header */}
      <div className="bg-white rounded-2xl border border-teal-100 shadow-sm p-7 mb-5">
        <div className="flex items-start gap-4 mb-4">
          <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-teal-100 to-emerald-200 flex items-center justify-center text-3xl shrink-0">
            🌿
          </div>
          <div>
            <h1 className="text-2xl md:text-3xl font-black text-teal-900 tracking-tight">
              {volunteer.Title}
            </h1>
          </div>
        </div>
        <div className="flex flex-wrap gap-2 mb-4">
          {volunteer.Division && (
            <span className="text-xs bg-slate-100 text-slate-600 px-2.5 py-1 rounded-full">
              {volunteer.Division}
            </span>
          )}
          {volunteer.Location && (
            <span className="text-xs bg-slate-100 text-slate-600 px-2.5 py-1 rounded-full">
              📍 {volunteer.Location}
            </span>
          )}
          <span className="text-xs bg-teal-100 text-teal-700 font-semibold px-2.5 py-1 rounded-full">
            {TYPE_LABEL[volunteer.Type] ?? volunteer.Type}
          </span>
        </div>
        <p className="text-xs text-slate-400">
          Dibuka: {new Date(volunteer.CreatedAt).toLocaleDateString('id-ID', {
            day: 'numeric', month: 'long', year: 'numeric',
          })}
        </p>
      </div>

      {/* Description */}
      <div className="bg-white rounded-2xl border border-teal-100 p-7 mb-5">
        <h2 className="font-bold text-teal-900 text-[17px] mb-4">Deskripsi Kegiatan</h2>
        <p className="text-slate-600 text-sm leading-[1.9] whitespace-pre-wrap">{volunteer.Description}</p>
      </div>

      {/* Requirements */}
      <div className="bg-white rounded-2xl border border-teal-100 p-7 mb-5">
        <h2 className="font-bold text-teal-900 text-[17px] mb-4">Kualifikasi &amp; Syarat</h2>
        <p className="text-slate-600 text-sm leading-[1.9] whitespace-pre-wrap">{volunteer.Requirements}</p>
      </div>

      {/* Benefits */}
      {volunteer.Benefits && (
        <div className="bg-teal-50 rounded-2xl border border-teal-200 p-7 mb-8">
          <h2 className="font-bold text-teal-900 text-[17px] mb-4">🎁 Yang Kamu Dapatkan</h2>
          <p className="text-slate-600 text-sm leading-[1.9] whitespace-pre-wrap">{volunteer.Benefits}</p>
        </div>
      )}

      {/* Success */}
      {success && (
        <div className="mb-6 p-5 bg-emerald-50 border border-emerald-200 rounded-2xl flex items-start gap-3">
          <span className="text-2xl">✅</span>
          <div>
            <p className="font-bold text-emerald-800 mb-0.5">Pendaftaran Berhasil!</p>
            <p className="text-sm text-emerald-700">
              Tim kami akan menghubungi Anda melalui email dalam 3–5 hari kerja. Terima kasih telah bergabung! 🌿
            </p>
          </div>
        </div>
      )}

      {/* Apply button */}
      {!success && !showForm && (
        <button
          onClick={() => setShowForm(true)}
          className="w-full py-4 rounded-2xl bg-gradient-to-r from-teal-600 to-emerald-600 text-white font-bold text-[15px] hover:from-teal-700 hover:to-emerald-700 transition-all shadow-lg shadow-teal-600/20 hover:-translate-y-0.5"
        >
          Daftar Jadi Relawan →
        </button>
      )}

      {/* Form */}
      {!success && showForm && (
        <div className="bg-white rounded-2xl border border-teal-100 shadow-sm p-7">
          <h2 className="font-bold text-teal-900 text-[17px] mb-6">Form Pendaftaran Relawan</h2>

          {applyMut.isError && (
            <div className="mb-5 p-4 bg-red-50 border border-red-200 rounded-xl flex items-start gap-2 text-red-600 text-sm">
              <span>⚠️</span> Gagal mengirim pendaftaran. Silakan coba lagi.
            </div>
          )}

          <form ref={formRef} onSubmit={handleSubmit} className="space-y-5">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Nama Lengkap *</label>
                <input
                  type="text" name="full_name" required
                  className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Email *</label>
                <input
                  type="email" name="email" required
                  className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Nomor HP *</label>
                <input
                  type="tel" name="phone" required
                  className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">Kota Domisili *</label>
                <input
                  type="text" name="city" required
                  className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1.5">Pekerjaan / Status</label>
              <input
                type="text" name="occupation"
                placeholder="Contoh: Mahasiswa, Karyawan, Freelancer..."
                className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                Unggah CV / Portofolio{' '}
                <span className="text-slate-400 font-normal">(PDF, maks 5MB)</span>
              </label>
              <div className="border-2 border-dashed border-teal-200 rounded-xl p-5 text-center bg-teal-50/40 hover:bg-teal-50 transition-colors cursor-pointer relative">
                <div className="text-2xl mb-1.5">📎</div>
                <p className="text-sm text-slate-500">Klik untuk upload atau drag &amp; drop</p>
                <p className="text-xs text-slate-400 mt-1">PDF, DOC, DOCX (maks 5MB)</p>
                <input
                  type="file" name="cv" accept=".pdf,.doc,.docx"
                  className="absolute inset-0 opacity-0 cursor-pointer"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1.5">Motivasi Menjadi Relawan *</label>
              <textarea
                name="motivation" rows={5} required
                placeholder="Ceritakan mengapa Anda ingin bergabung dan apa yang bisa Anda kontribusikan..."
                className="w-full px-3.5 py-2.5 rounded-xl border-2 border-teal-100 text-sm focus:outline-none focus:border-teal-500 transition-colors resize-y"
              />
            </div>

            <div className="flex gap-3 pt-1">
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="flex-1 py-2.5 rounded-xl border-2 border-teal-100 text-sm font-semibold text-slate-600 hover:bg-slate-50 transition-colors"
              >
                Batal
              </button>
              <button
                type="submit"
                disabled={applyMut.isPending}
                className="flex-[2] py-2.5 rounded-xl bg-teal-600 text-white text-sm font-bold hover:bg-teal-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
              >
                {applyMut.isPending ? (
                  <>
                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin inline-block" />
                    Mengirim...
                  </>
                ) : 'Kirim Pendaftaran 🌿'}
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  )
}