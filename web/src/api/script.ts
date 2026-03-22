import request from '@/utils/request'
import type { ApiResponse, Script } from '@/types/api'

export const scriptApi = {
  list: (subscriptionId: number) => request.get<ApiResponse<{ scripts: Script[] }>>(`/subscriptions/${subscriptionId}/scripts`),
  get: (subscriptionId: number, scriptId: number) => request.get<ApiResponse<{ script: Script }>>(`/subscriptions/${subscriptionId}/scripts/${scriptId}`),
  create: (subscriptionId: number, data: Partial<Script>) => request.post<ApiResponse<{ script: Script }>>(`/subscriptions/${subscriptionId}/scripts`, data),
  update: (subscriptionId: number, scriptId: number, data: Partial<Script>) => request.put<ApiResponse<{ script: Script }>>(`/subscriptions/${subscriptionId}/scripts/${scriptId}`, data),
  delete: (subscriptionId: number, scriptId: number) => request.delete<ApiResponse>(`/subscriptions/${subscriptionId}/scripts/${scriptId}`),
  test: (subscriptionId: number, scriptId: number, config: Record<string, unknown>) =>
    request.post<ApiResponse<{ result: Record<string, unknown> }>>(`/subscriptions/${subscriptionId}/scripts/${scriptId}/test`, { config }),
}