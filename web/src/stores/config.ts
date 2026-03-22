import { defineStore } from 'pinia'
import { ref } from 'vue'
import { configApi } from '@/api/config'
import type { Config, CoreConfig } from '@/types/api'

export const useConfigStore = defineStore('config', () => {
  const config = ref<Config | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchConfig = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await configApi.get()
      config.value = res.data.data as Config
    } catch (err) {
      console.error('Failed to fetch config:', err)
      error.value = err instanceof Error ? err.message : '加载配置失败'
    } finally {
      loading.value = false
    }
  }

  const updateCoreConfig = async (core: Partial<CoreConfig>) => {
    const res = await configApi.update({ core })
    if (config.value && res.data.data?.core) {
      config.value.core = res.data.data.core
    }
    return res.data.data
  }

  return {
    config,
    loading,
    error,
    fetchConfig,
    updateCoreConfig,
  }
})
