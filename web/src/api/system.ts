import request from '@/utils/request'
import type { ApiResponse, SystemInfo } from '@/types/api'

export const systemApi = {
  getSystemInfo: () => request.get<ApiResponse<SystemInfo>>('/system/info'),
}
