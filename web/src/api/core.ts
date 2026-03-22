import request from '@/utils/request'
import type { ApiResponse, CoreRulesData } from '@/types/api'

export const coreApi = {
  getRules: () => request.get<ApiResponse<CoreRulesData>>('/rules'),
}