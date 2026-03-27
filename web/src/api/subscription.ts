import request, { type RequestConfig } from '@/utils/request'
import type { ApiResponse, Subscription, SubscriptionWithCounts, SubscriptionDetail, RefreshResult, MergedConfig } from '@/types/api'

const noTimeout: RequestConfig = { timeout: 0 }

export const subscriptionApi = {
  list: () => request.get<ApiResponse<{ subscriptions: SubscriptionWithCounts[] }>>('/subscriptions'),
  get: (id: number) => request.get<ApiResponse<SubscriptionDetail>>(`/subscriptions/${id}`),
  create: (data: Partial<Subscription>) => request.post<ApiResponse<{ subscription: Subscription }>>('/subscriptions', data, noTimeout),
  update: (id: number, data: Partial<Subscription>) => request.put<ApiResponse<{ subscription: Subscription }>>(`/subscriptions/${id}`, data, noTimeout),
  delete: (id: number) => request.delete<ApiResponse>(`/subscriptions/${id}`),
  refresh: (id: number) => request.post<ApiResponse<RefreshResult>>(`/subscriptions/${id}/refresh`, undefined, noTimeout),
  activate: (id: number) => request.put<ApiResponse>(`/subscriptions/${id}/activate`),
  getMergedConfig: (id: number) => request.get<ApiResponse<MergedConfig>>(`/subscriptions/${id}/merged-config`),
  getContent: (id: number) => request.get<ApiResponse<{ content: string }>>(`/subscriptions/${id}/content`),
  updateContent: (id: number, content: string) => request.put<ApiResponse>(`/subscriptions/${id}/content`, { content }),
}