import request from '@/utils/request'
import type { ApiResponse, Connection } from '@/types/api'

export const connectionApi = {
  list: () => request.get<ApiResponse<{ connections: Connection[] }>>('/connections'),
  close: (id: string) => request.delete<ApiResponse>(`/connections/${id}`),
  closeAll: () => request.delete<ApiResponse>('/connections'),
}
