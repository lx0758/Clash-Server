import request from '@/utils/request'
import type { ApiResponse, CoreConfig } from '@/types/api'

interface ConfigResponse {
  server: {
    host: string
    port: number
    database: string
  }
  core: CoreConfig
}

interface ConfigUpdateRequest {
  core: Partial<CoreConfig>
}

interface ConfigUpdateResponse {
  core: CoreConfig
  core_error?: string
}

export const configApi = {
  get: () => request.get<ApiResponse<ConfigResponse>>('/config'),
  update: (data: ConfigUpdateRequest) => request.put<ApiResponse<ConfigUpdateResponse>>('/config', data),
}
