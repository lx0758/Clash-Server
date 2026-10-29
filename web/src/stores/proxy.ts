import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { proxyApi, type ProxyMode } from '@/api/proxy'
import type { Proxy, ProxyGroup, MihomoProxy, MihomoProxyGroup } from '@/types/api'
import type { DelayResult } from '@/composables/useWebSocket'

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

  /**
   * 拉取代理列表。
   * silent 为 true 时不设置全局 loading, 用于测延迟完成后的静默同步,
   * 避免整页塌陷为骨架屏。
   */
  const fetchProxies = async (silent = false) => {
    if (!silent) {
      loading.value = true
    }
    try {
      const res = await proxyApi.list()
      rawProxies.value = res.data.data.proxies
    } finally {
      if (!silent) {
        loading.value = false
      }
    }
  }

  const fetchMode = async () => {
    try {
      const res = await proxyApi.getMode()
      mode.value = res.data.data.mode
    } catch {
      mode.value = 'rule'
    }
  }

  const setMode = async (newMode: ProxyMode) => {
    await proxyApi.setMode(newMode)
    mode.value = newMode
  }

  /**
   * 切换代理组选中的节点。
   * 采用乐观更新: 立即反映用户选择, 服务端拒绝时回滚。
   */
  const selectProxy = async (group: string, name: string) => {
    const g = rawProxies.value[group] as MihomoProxyGroup | undefined
    const previous = g?.now
    if (g) {
      rawProxies.value = {
        ...rawProxies.value,
        [group]: { ...g, now: name }
      }
    }
    try {
      await proxyApi.select(group, name)
    } catch (err) {
      // 回滚到切换前的节点
      const cur = rawProxies.value[group] as MihomoProxyGroup | undefined
      if (cur && previous !== undefined) {
        rawProxies.value = {
          ...rawProxies.value,
          [group]: { ...cur, now: previous }
        }
      }
      throw err
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
    const res = await proxyApi.checkDelay(proxyName)
    return res.data.data.delay
  }

  /**
   * 启动组内逐节点延迟测试。
   *
   * 触发后立即返回, 各节点的延迟结果通过 WebSocket 的 delay 消息陆续送达,
   * 由 handleDelayResult 逐个更新, 实现"快的先显示"。
   */
  const startGroupDelayTest = async (groupName: string) => {
    return proxyApi.startGroupDelayTest(groupName)
  }

  /**
   * 处理单个节点的延迟结果。
   * 只更新目标节点的 delay 字段, 不触碰列表整体。
   */
  const handleDelayResult = (result: DelayResult) => {
    if (!result.name) return
    const target = rawProxies.value[result.name] as MihomoProxy | undefined
    if (!target) return

    const delay = result.error ? 0 : (result.delay ?? 0)
    const history = [...(target.history || []), { delay, time: '' }]
    rawProxies.value = {
      ...rawProxies.value,
      [result.name]: { ...target, history },
    }
  }

  /**
   * 测延迟全部完成后的静默同步。
   * 单节点测速会把延迟写入核心的代理历史, 需要重新拉取以同步。
   * 使用静默模式, 不触发整页骨架屏。
   */
  const syncAfterDelayTest = async () => {
    await fetchProxies(true)
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
    startGroupDelayTest,
    handleDelayResult,
    syncAfterDelayTest,
  }
})