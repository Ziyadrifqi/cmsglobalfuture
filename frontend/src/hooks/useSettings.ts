import { useQuery } from '@tanstack/react-query'
import { portalApi } from '../api'

export function useSettings() {
  return useQuery({
    queryKey: ['site-settings'],
    queryFn: portalApi.settings,
    staleTime: 5 * 60 * 1000,
  })
}
