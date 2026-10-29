import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'
import { useLoadingStore } from '@/stores/loading'
import { useToastStore } from '@/stores/toast'

interface RequestConfig extends AxiosRequestConfig {
  showLoading?: boolean
  showErrorToast?: boolean
}

/**
 * 业务错误。
 *
 * 携带后端的业务状态码与响应数据, 使调用方能区分不同失败语义:
 *   - 完全失败: code 非 0 且 data.partial_success 不为 true
 *   - 半成功:   code = 1007 且 data.partial_success = true
 *     表示数据已持久化但核心未生效, 应以警告而非错误语义提示
 */
export class ApiError extends Error {
  /** 后端业务状态码 */
  readonly code: number
  /** 后端响应体中的 data 部分 */
  readonly data: unknown

  constructor(message: string, code: number, data?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.data = data
  }

  /**
   * 是否为"已保存但核心未生效"的半成功状态。
   */
  get isPartialSuccess(): boolean {
    if (this.code !== 1007) return false
    if (!this.data || typeof this.data !== 'object') return false
    return (this.data as Record<string, unknown>).partial_success === true
  }

  /**
   * 核心未生效的原因, 仅半成功状态下有值。
   */
  get coreError(): string {
    if (!this.data || typeof this.data !== 'object') return ''
    const v = (this.data as Record<string, unknown>).core_error
    return typeof v === 'string' ? v : ''
  }
}

const instance: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
})

instance.interceptors.request.use(
  (config) => {
    const loadingStore = useLoadingStore()
    const customConfig = config as RequestConfig
    if (customConfig.showLoading !== false) {
      loadingStore.show()
    }
    return config
  },
  (error) => {
    const loadingStore = useLoadingStore()
    loadingStore.hide()
    return Promise.reject(error)
  }
)

instance.interceptors.response.use(
  (response) => {
    const loadingStore = useLoadingStore()
    const toastStore = useToastStore()
    loadingStore.hide()

    const data = response.data as ApiResponse
    const customConfig = response.config as RequestConfig
    if (data.code !== 0) {
      const payload = data.data as unknown
      const apiError = new ApiError(data.message || '请求失败', data.code, payload)

      // 业务码 401: 会话已失效, 跳转登录页且不弹冗余提示
      // (后端以 HTTP 200 + code 401 表示未登录, 不走下面的 error 分支)
      if (data.code === 401 && !customConfig.url?.includes('/session')) {
        window.location.href = '/login'
        return Promise.reject(apiError)
      }

      // 半成功(数据已保存但核心未生效)不弹错误 toast,
      // 交由业务层按警告语义呈现, 避免把警告渲染成错误
      if (!apiError.isPartialSuccess && customConfig.showErrorToast !== false) {
        toastStore.error(data.message || '请求失败')
      }
      return Promise.reject(apiError)
    }
    return response
  },
  (error) => {
    const loadingStore = useLoadingStore()
    const toastStore = useToastStore()
    loadingStore.hide()

    const customConfig = error.config as RequestConfig | undefined
    if (error.response?.status === 401) {
      if (!error.config?.url?.includes('/session')) {
        window.location.href = '/login'
      }
      return Promise.reject(error)
    }

    if (customConfig?.showErrorToast !== false) {
      let message = '网络错误，请稍后重试'
      if (error.code === 'ECONNABORTED') {
        message = '请求超时，请稍后重试'
      } else if (error.response?.data?.message) {
        message = error.response.data.message
      } else if (error.message) {
        message = error.message
      }
      toastStore.error(message)
    }
    return Promise.reject(error)
  }
)

export type { RequestConfig }
export default instance
