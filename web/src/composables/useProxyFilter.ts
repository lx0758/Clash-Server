import type { Proxy } from '@/types/api'
import type { SortType } from '@/stores/proxy'

export function useProxyFilter() {
  const filterSort = (
    proxies: Proxy[],
    filterText: string,
    sortType: SortType
  ): Proxy[] => {
    let result = [...proxies]

    if (filterText) {
      const lowerFilter = filterText.toLowerCase()
      result = result.filter(p =>
        p.name.toLowerCase().includes(lowerFilter)
      )
    }

    switch (sortType) {
      case 'delay-asc':
        result.sort((a, b) => {
          const aDelay = a.delay || 999999
          const bDelay = b.delay || 999999
          return aDelay - bDelay
        })
        break
      case 'delay-desc':
        result.sort((a, b) => {
          const aDelay = a.delay || 0
          const bDelay = b.delay || 0
          return bDelay - aDelay
        })
        break
      case 'name':
        result.sort((a, b) => a.name.localeCompare(b.name))
        break
    }

    return result
  }

  const getDelayColor = (delay: number): string => {
    if (delay === 0) return '#e94560'
    if (delay < 100) return '#2ed573'
    if (delay < 300) return '#7bed9f'
    if (delay < 500) return '#ffa502'
    return '#e94560'
  }

  const formatDelay = (delay: number): string => {
    if (delay === 0) return '超时'
    return `${delay}ms`
  }

  return {
    filterSort,
    getDelayColor,
    formatDelay,
  }
}