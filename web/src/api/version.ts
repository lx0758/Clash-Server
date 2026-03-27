import request from '@/utils/request'

export interface VersionResponse {
  version: string
}

export const versionApi = {
  getVersion: () => request.get<{ data: VersionResponse }>('/version'),
}
