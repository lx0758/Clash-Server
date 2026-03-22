import request from '@/utils/request'
import type { ApiResponse, Rule } from '@/types/api'

export const ruleApi = {
  list: (subscriptionId: number) => request.get<ApiResponse<{ rules: Rule[] }>>(`/subscriptions/${subscriptionId}/rules`),
  get: (subscriptionId: number, ruleId: number) => request.get<ApiResponse<{ rule: Rule }>>(`/subscriptions/${subscriptionId}/rules/${ruleId}`),
  create: (subscriptionId: number, data: Partial<Rule>) => request.post<ApiResponse<{ rule: Rule }>>(`/subscriptions/${subscriptionId}/rules`, data),
  update: (subscriptionId: number, ruleId: number, data: Partial<Rule>) => request.put<ApiResponse<{ rule: Rule }>>(`/subscriptions/${subscriptionId}/rules/${ruleId}`, data),
  delete: (subscriptionId: number, ruleId: number) => request.delete<ApiResponse>(`/subscriptions/${subscriptionId}/rules/${ruleId}`),
}