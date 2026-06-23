import { useEffect, useState } from 'react'
import { Link, NavLink } from 'react-router-dom'
import { useSettings } from '../../hooks/useSettings'
import { getText } from '../../lib/settings'

const links = [
  { to: '/',         label: 'Beranda' },
  { to: '/tentang',  label: 'Tentang' },
  { to: '/berita',   label: 'Berita' },
  { to: '/galeri',   label: 'Galeri' },
  { to: '/relawan',  label: 'Relawan' },
  { to: '/kontak',   label: 'Kontak' },
]

export function Navbar() {
  const [open, setOpen]         = useState(false)
  const [scrolled, setScrolled] = useState(false)
  const { data: settings } = useSettings()

  const siteName    = getText(settings, 'site_name', 'Green Future')
  const siteNameSub = getText(settings, 'site_name_sub', 'Indonesia')
  const logoImage   = settings?.site_logo_image

  useEffect(() => {
    const fn = () => setScrolled(window.scrollY > 10)
    window.addEventListener('scroll', fn)
    return () => window.removeEventListener('scroll', fn)
  }, [])

  return (
    <header
      className={`sticky top-0 z-50 bg-white/95 backdrop-blur-md transition-shadow duration-300 ${
        scrolled ? 'shadow-[0_2px_20px_rgba(15,118,110,0.12)]' : 'border-b border-teal-100'
      }`}
    >
      <div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">

        {/* Logo */}
        <Link to="/" className="flex items-center gap-2.5">
          {logoImage ? (
            <img src={logoImage} alt={siteName} className="w-9 h-9 rounded-xl object-cover shadow-sm" />
          ) : (
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-teal-600 to-teal-900 flex items-center justify-center text-lg shadow-sm">
              🌿
            </div>
          )}
          <div className="leading-tight">
            <span className="font-extrabold text-base text-teal-900 tracking-tight block leading-none">
              {siteName}
            </span>
            <span className="text-[10px] text-teal-600 font-semibold tracking-widest uppercase">{siteNameSub}</span>
          </div>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden md:flex items-center gap-1">
          {links.map(l => (
            <NavLink
              key={l.to} to={l.to} end={l.to === '/'}
              className={({ isActive }) =>
                `px-3 py-2 rounded-lg text-sm font-medium transition-all duration-150 ${
                  isActive
                    ? 'bg-teal-50 text-teal-700 font-semibold'
                    : 'text-slate-600 hover:text-teal-700 hover:bg-teal-50/60'
                }`
              }
            >
              {l.label}
            </NavLink>
          ))}
        </nav>

        {/* CTA */}
        <Link
          to="/relawan"
          className="hidden md:inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-teal-600 text-white text-sm font-semibold hover:bg-teal-700 transition-colors shadow-sm"
        >
          Jadi Relawan →
        </Link>

        {/* Mobile hamburger */}
        <button
          onClick={() => setOpen(!open)}
          className="md:hidden w-9 h-9 flex items-center justify-center rounded-lg text-teal-800 hover:bg-teal-50 transition-colors"
          aria-label="Toggle menu"
        >
          {open ? '✕' : '☰'}
        </button>
      </div>

      {open && (
        <div className="md:hidden border-t border-teal-100 px-4 py-3 space-y-1 bg-white">
          {links.map(l => (
            <NavLink
              key={l.to} to={l.to} end={l.to === '/'}
              onClick={() => setOpen(false)}
              className={({ isActive }) =>
                `block py-2.5 px-3 rounded-lg text-sm font-medium transition-colors ${
                  isActive ? 'bg-teal-50 text-teal-700' : 'text-slate-600 hover:bg-slate-50'
                }`
              }
            >
              {l.label}
            </NavLink>
          ))}
          <Link
            to="/relawan" onClick={() => setOpen(false)}
            className="block mt-2 py-2.5 px-3 rounded-lg text-sm font-semibold text-center bg-teal-600 text-white"
          >
            Jadi Relawan →
          </Link>
        </div>
      )}
    </header>
  )
}
