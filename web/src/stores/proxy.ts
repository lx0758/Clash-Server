import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { proxyApi, type ProxyMode } from '@/api/proxy'
import type { Proxy, ProxyGroup, MihomoProxy, MihomoProxyGroup } from '@/types/api'

export type SortType = 'default' | 'delay-asc' | 'delay-desc' | 'name'

export interface ProxyGroupState {
  open: boolean
  filterText: string
  sortType: SortType
}

const STORAGE_KEY = 'proxy-group-states'

function loadGroupStates(): Record<string, ProxyGroupState> {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) return JSON.parse(saved)
  } catch {}
  return {}
}

function saveGroupStates(states: Record<string, ProxyGroupState>) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(states))
  } catch {}
}

export const useProxyStore = defineStore('proxy', () => {
  const rawProxies = ref<Record<string, MihomoProxy | MihomoProxyGroup>>({})
  const groupStates = ref<Record<string, ProxyGroupState>>(loadGroupStates())
  const loading = ref(false)
  const mode = ref<ProxyMode>('rule')

  const getProxyInfo = (name: string): Proxy => {
    const p = rawProxies.value[name]
    if (!p) {
      return { name, type: 'Unknown', alive: false, delay: 0 }
    }
    const history = (p as MihomoProxy).history || []
    const lastDelay = history.length > 0 ? history[history.length - 1]?.delay ?? 0 : 0
    return {
      name: p.name,
      type: p.type,
      alive: (p as MihomoProxy).alive ?? false,
      delay: lastDelay,
    }
  }

  const proxyGroups = computed(() => {
    const groups = Object.entries(rawProxies.value)
      .filter(([name, p]) => {
        if (name === 'GLOBAL') return false
        const type = p.type.toLowerCase()
        return type === 'selector' || type === 'urltest' || type === 'fallback'
      })
      .map(([name, group]) => {
        const g = group as MihomoProxyGroup
        const allProxies = (g.all || []).map(proxyName => getProxyInfo(proxyName))
        return {
          name,
          type: g.type,
          all: allProxies,
          now: g.now || '',
        } as ProxyGroup
      })
    
    const typeOrder: Record<string, number> = {
      selector: 0,
      urltest: 1,
      fallback: 2,
    }
    groups.sort((a, b) => {
      const aOrder = typeOrder[a.type.toLowerCase()] ?? 99
      const bOrder = typeOrder[b.type.toLowerCase()] ?? 99
      return aOrder - bOrder
    })
    return groups
  })

  const globalGroup = computed(() => {
    const g = rawProxies.value['GLOBAL'] as MihomoProxyGroup | undefined
    if (!g) return null
    const allProxies = (g.all || []).map(proxyName => getProxyInfo(proxyName))
    return {
      name: 'GLOBAL',
      type: g.type,
      all: allProxies,
      now: g.now || '',
    } as ProxyGroup
  })

  const globalProxy = computed(() => rawProxies.value['GLOBAL'])

  const fetchProxies = async () => {
    loading.value = true
    try {
      const res = await proxyApi.list()
      rawProxies.value = res.data.data.proxies
    } finally {
      loading.value = false
    }
  }

  const fetchMode = async () => {
    try {
      const result = await proxyApi.getMode()
      mode.value = result.mode
    } catch {
      mode.value = 'rule'
    }
  }

  const setMode = async (newMode: ProxyMode) => {
    await proxyApi.setMode(newMode)
    mode.value = newMode
  }

  const selectProxy = async (group: string, name: string) => {
    await proxyApi.select(group, name)
    const g = rawProxies.value[group] as MihomoProxyGroup | undefined
    if (g) {
      rawProxies.value = {
        ...rawProxies.value,
        [group]: { ...g, now: name }
      }
    }
  }

  const getGroupState = (groupName: string): ProxyGroupState => {
    if (!groupStates.value[groupName]) {
      groupStates.value[groupName] = {
        open: false,
        filterText: '',
        sortType: 'default',
      }
    }
    return groupStates.value[groupName]
  }

  const toggleGroup = (groupName: string) => {
    const state = getGroupState(groupName)
    state.open = !state.open
    saveGroupStates(groupStates.value)
  }

  const setGroupFilter = (groupName: string, filterText: string) => {
    const state = getGroupState(groupName)
    state.filterText = filterText
    saveGroupStates(groupStates.value)
  }

  const setGroupSort = (groupName: string, sortType: SortType) => {
    const state = getGroupState(groupName)
    state.sortType = sortType
    saveGroupStates(groupStates.value)
  }

  const checkDelay = async (proxyName: string): Promise<number> => {
    const result = await proxyApi.checkDelay(proxyName)
    return result.delay
  }

  const checkGroupDelay = async (groupName: string) => {
    await proxyApi.checkGroupDelay(groupName)
    await fetchProxies()
  }

  return {
    rawProxies,
    proxyGroups,
    globalGroup,
    globalProxy,
    groupStates,
    loading,
    mode,
    fetchProxies,
    fetchMode,
    setMode,
    selectProxy,
    getGroupState,
    toggleGroup,
    setGroupFilter,
    setGroupSort,
    checkDelay,
    checkGroupDelay,
  }
})