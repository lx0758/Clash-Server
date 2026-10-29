import request from '@/utils/request'
import type { ApiResponse, MihomoProxyGroup, MihomoProxy } from '@/types/api'

export type ProxyMode = 'rule' | 'global' | 'direct'

export const proxyApi = {
  list: () => request.get<ApiResponse<{ proxies: Record<string, MihomoProxy | MihomoProxyGroup> }>>('/proxies'),
  get: (name: string) => request.get<ApiResponse<{ proxy: MihomoProxy | MihomoProxyGroup }>>(`/proxies/${encodeURIComponent(name)}`),
  select: (group: string, name: string) => request.put<ApiResponse<void>>(`/proxies/${encodeURIComponent(group)}`, { name }),
  checkDelay: (name: string, testUrl?: string, timeout?: number) =>
    request.get<ApiResponse<{ delay: number }>>(`/proxies/${encodeURIComponent(name)}/delay`, { params: { url: testUrl, timeout } }),
  checkGroupDelay: (group: string, testUrl?: string, timeout?: number) =>
    request.get<ApiResponse<void>>(`/proxies/group/${encodeURIComponent(group)}/delay`, { params: { url: testUrl, timeout } }),
  /**
   * 启动组内逐节点延迟测试。
   * 接口立即返回, 各节点的延迟结果经 WebSocket 的 delay 消息陆续送达。
   */
  startGroupDelayTest: (group: string, testUrl?: string, timeout?: number) =>
    request.post<ApiResponse<{ group: string; nodes: string[]; total: number }>>(
      `/proxies/group/${encodeURIComponent(group)}/delay-test`,
      undefined,
      { params: { url: testUrl, timeout } }
    ),
  getMode: () => request.get<ApiResponse<{ mode: ProxyMode }>>('/proxies/mode'),
  setMode: (mode: ProxyMode) => request.patch<ApiResponse<void>>('/proxies/mode', { mode }),
}
