import request from '@/utils/request'
import type { ApiResponse, Customization } from '@/types/api'

export const customizationApi = {
  get: (subscriptionId: number) => 
    request.get<ApiResponse<{ customization: Customization | null }>>(`/subscriptions/${subscriptionId}/customization`),
  update: (subscriptionId: number, data: Partial<Customization>) => 
    request.put<ApiResponse<{ customization: Customization }>>(`/subscriptions/${subscriptionId}/customization`, data),
}
