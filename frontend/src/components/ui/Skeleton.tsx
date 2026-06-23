export function Skeleton({ className = '' }: { className?: string }) {
  return (
    <div
      className={`animate-pulse rounded-xl bg-gradient-to-r from-teal-50 via-white to-teal-50 bg-[length:200%_100%] ${className}`}
      style={{ animation: 'shimmer 1.4s ease-in-out infinite, pulse 2s ease-in-out infinite' }}
    />
  )
}