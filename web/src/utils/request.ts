import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'
import { useLoadingStore } from '@/stores/loading'
import { useToastStore } from '@/stores/toast'

interface RequestConfig extends AxiosRequestConfig {
  showLoading?: boolean
  showErrorToast?: boolean
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
      if (customConfig.showErrorToast !== false) {
        toastStore.error(data.message || '请求失败')
      }
      return Promise.reject(new Error(data.message))
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

export async function get<T>(url: string, config?: RequestConfig) {
  const res = await instance.get<ApiResponse<T>>(url, config)
  return res.data.data
}

export async function post<T>(url: string, data?: unknown, config?: RequestConfig) {
  const res = await instance.post<ApiResponse<T>>(url, data, config)
  return res.data.data
}

export async function put<T>(url: string, data?: unknown, config?: RequestConfig) {
  const res = await instance.put<ApiResponse<T>>(url, data, config)
  return res.data.data
}

export async function patch<T>(url: string, data?: unknown, config?: RequestConfig) {
  const res = await instance.patch<ApiResponse<T>>(url, data, config)
  return res.data.data
}

export async function del<T>(url: string, config?: RequestConfig) {
  const res = await instance.delete<ApiResponse<T>>(url, config)
  return res.data.data
}

export type { RequestConfig }
export default instance
