import { Routes, Route, Link } from 'react-router-dom'
import { Navbar }               from './components/layout/Navbar'
import { Footer }               from './components/layout/Footer'
import { HomePage }             from './pages/HomePage'
import { NewsListPage }         from './pages/NewsListPage'
import { NewsDetailPage }       from './pages/NewsDetailPage'
import { VolunteerListPage }    from './pages/VolunteerListPage'
import { VolunteerDetailPage }  from './pages/VolunteerDetailPage'
import { GalleryPage }          from './pages/GalleryPage'
import { AboutPage }            from './pages/AboutPage'
import { ContactPage }          from './pages/ContactPage'

function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center py-28 gap-4 px-4 text-center">
      <div className="text-8xl font-black text-teal-100 leading-none">404</div>
      <div className="text-5xl">🌱</div>
      <h1 className="text-2xl font-black text-teal-900">Halaman Tidak Ditemukan</h1>
      <p className="text-slate-500 max-w-xs">Sepertinya halaman ini sudah kembali ke alam. Yuk balik ke beranda!</p>
      <Link to="/" className="mt-2 px-6 py-3 rounded-xl bg-teal-600 text-white text-sm font-bold hover:bg-teal-700 transition-colors">
        ← Kembali ke Beranda
      </Link>
    </div>
  )
}

export default function App() {
  return (
    <div className="min-h-screen flex flex-col bg-[#F8FAFC]">
      <Navbar />
      <main className="flex-1 animate-fade-in">
        <Routes>
          <Route path="/"                element={<HomePage />} />
          <Route path="/berita"          element={<NewsListPage />} />
          <Route path="/berita/:slug"    element={<NewsDetailPage />} />
          <Route path="/relawan"         element={<VolunteerListPage />} />
          <Route path="/relawan/:slug"   element={<VolunteerDetailPage />} />
          <Route path="/galeri"          element={<GalleryPage />} />
          <Route path="/tentang"         element={<AboutPage />} />
          <Route path="/kontak"          element={<ContactPage />} />
          <Route path="*"                element={<NotFoundPage />} />
        </Routes>
      </main>
      <Footer />
    </div>
  )
}