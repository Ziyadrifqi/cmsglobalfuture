import typography from '@tailwindcss/typography'

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: '#0F766E',
          secondary: '#134E4A',
          accent: '#F59E0B',
          bg: '#F8FAFC',
        },
      },
      keyframes: {
        'fade-in': {
          from: { opacity: '0', transform: 'translateY(10px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '200% 0' },
          '100%': { backgroundPosition: '-200% 0' },
        },
      },
      animation: {
        'fade-in': 'fade-in 0.3s ease both',
        shimmer: 'shimmer 1.4s ease-in-out infinite',
      },
      boxShadow: {
        'teal-sm': '0 2px 12px rgba(15,118,110,0.08)',
        'teal-md': '0 4px 24px rgba(15,118,110,0.14)',
        'teal-lg': '0 12px 40px rgba(15,118,110,0.18)',
      },
    },
  },
  plugins: [typography],
}