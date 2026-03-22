import instance, { type RequestConfig } from '@/utils/request'
import type { ApiResponse, User } from '@/types/api'

const noLoadingNoToast: RequestConfig = { showLoading: false, showErrorToast: false }

export const authApi = {
  checkInit: () => instance.get<ApiResponse<{ initialized: boolean }>>('/init', noLoadingNoToast),
  initPassword: (username: string, password: string) =>
    instance.post<ApiResponse>('/init', { username, password }, noLoadingNoToast),
  login: (username: string, password: string) =>
    instance.post<ApiResponse<{ user: User }>>('/session', { username, password }, noLoadingNoToast),
  logout: () => instance.delete<ApiResponse>('/session'),
  getCurrentUser: () => instance.get<ApiResponse<{ user: User }>>('/users/me', noLoadingNoToast),
  changePassword: (old_password: string, new_password: string) =>
    instance.put<ApiResponse>('/users/me/password', { old_password, new_password }),
}
