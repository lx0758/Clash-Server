import { defineStore } from 'pinia'
import { ref } from 'vue'
import { systemApi } from '@/api/system'
import type { SystemInfo, CoreStatus } from '@/types/api'

export const useSystemStore = defineStore('system', () => {
  const systemInfo = ref<SystemInfo | null>(null)
  const coreStatus = ref<CoreStatus>({ running: false })
  const loading = ref(false)

  const fetchSystemInfo = async () => {
    loading.value = true
    try {
      const res = await systemApi.getSystemInfo()
      systemInfo.value = res.data.data
      coreStatus.value = res.data.data.core
    } finally {
      loading.value = false
    }
  }

  const updateCoreStatus = (status: CoreStatus) => {
    coreStatus.value = status
    if (systemInfo.value) {
      systemInfo.value.core = status
    }
  }

  return {
    systemInfo,
    coreStatus,
    loading,
    fetchSystemInfo,
    updateCoreStatus,
  }
})
