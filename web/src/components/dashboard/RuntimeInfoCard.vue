<template>
  <div class="runtime-info">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    <template v-else>
      <div class="info-grid">
        <div class="info-item">
          <el-icon class="info-icon" :class="{ running: coreRunning }">
            <CircleCheckFilled v-if="coreRunning" />
            <CircleCloseFilled v-else />
          </el-icon>
          <div class="info-content">
            <span class="info-label">核心状态</span>
            <span class="info-value" :class="{ running: coreRunning }">
              {{ coreRunning ? (systemStore.coreStatus.version || '运行中') : '已停止' }}
            </span>
          </div>
        </div>

        <div class="info-item">
          <el-icon class="info-icon"><Coin /></el-icon>
          <div class="info-content">
            <span class="info-label">核心内存</span>
            <span class="info-value">{{ formatBytes(memory.inuse) }}</span>
          </div>
        </div>

        <div class="info-item" v-if="coreConfig.mixed_port">
          <el-icon class="info-icon"><Connection /></el-icon>
          <div class="info-content">
            <span class="info-label">混合端口</span>
            <span class="info-value">{{ coreConfig.mixed_port }}</span>
          </div>
        </div>

        <div class="info-item" v-if="coreConfig.log_level">
          <el-icon class="info-icon"><Document /></el-icon>
          <div class="info-content">
            <span class="info-label">日志级别</span>
            <span class="info-value">{{ coreConfig.log_level }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useWebSocket } from '@/composables/useWebSocket'
import request from '@/utils/request'
import { formatBytes } from '@/utils/format'
import { 
  Loading, 
  CircleCheckFilled, 
  CircleCloseFilled, 
  Coin, 
  Connection, 
  Document 
} from '@element-plus/icons-vue'
import type { CoreConfig } from '@/types/api'

const systemStore = useSystemStore()
const { memory, subscribe } = useWebSocket()

const loading = ref(true)
const coreConfig = ref<CoreConfig>({
  api_host: '',
  api_port: 0,
  api_secret: '',
  mixed_port: 0,
  allow_lan: false,
  mode: '',
  log_level: '',
  ipv6: false,
})

const coreRunning = computed(() => systemStore.coreStatus.running)

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get('/config')
    if (res.data?.data?.core) {
      coreConfig.value = res.data.data.core
    }
  } catch (error) {
    console.error('Failed to fetch config:', error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchConfig()
  await systemStore.fetchSystemInfo()
  subscribe(['memory'])
})
</script>

<style scoped>
.runtime-info {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px;
  color: #94a3b8;
  font-size: 13px;
}

.loading-state .el-icon {
  font-size: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  width: 100%;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.info-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  border-radius: 8px;
  color: #64748b;
  font-size: 18px;
  flex-shrink: 0;
}

.info-icon.running {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  overflow: hidden;
}

.info-label {
  font-size: 12px;
  color: #94a3b8;
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
}

.info-value.running {
  color: #10b981;
}

@media (max-width: 992px) {
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
}

@media (max-width: 576px) {
  .info-icon {
    width: 32px;
    height: 32px;
    font-size: 16px;
  }

  .info-value {
    font-size: 13px;
  }
}
</style>
