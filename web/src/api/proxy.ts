import request, { get, put, patch } from '@/utils/request'
import type { ApiResponse, MihomoProxyGroup, MihomoProxy } from '@/types/api'

export type ProxyMode = 'rule' | 'global' | 'direct'

export const proxyApi = {
  list: () => request.get<ApiResponse<{ proxies: Record<string, MihomoProxy | MihomoProxyGroup> }>>('/proxies'),
  get: (name: string) => request.get<ApiResponse<{ proxy: MihomoProxy | MihomoProxyGroup }>>(`/proxies/${encodeURIComponent(name)}`),
  select: (group: string, name: string) => put<void>(`/proxies/${encodeURIComponent(group)}`, { name }),
  checkDelay: (name: string, testUrl?: string, timeout?: number) =>
    get<{ delay: number }>(`/proxies/${encodeURIComponent(name)}/delay`, { params: { url: testUrl, timeout } }),
  checkGroupDelay: (group: string, testUrl?: string, timeout?: number) =>
    get<void>(`/proxies/group/${encodeURIComponent(group)}/delay`, { params: { url: testUrl, timeout } }),
  getMode: () => get<{ mode: ProxyMode }>('/proxies/mode'),
  setMode: (mode: ProxyMode) => patch<void>('/proxies/mode', { mode }),
}