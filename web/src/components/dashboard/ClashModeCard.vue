<template>
  <div class="clash-mode-card">
    <div class="mode-buttons">
      <div
        v-for="mode in modes"
        :key="mode.key"
        class="mode-button"
        :class="{ active: currentMode === mode.key }"
        @click="handleModeChange(mode.key)"
      >
        <el-icon><component :is="mode.icon" /></el-icon>
        <span>{{ mode.label }}</span>
      </div>
    </div>
    <div class="mode-description">
      {{ currentModeDescription }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useProxyStore } from '@/stores/proxy'
import { proxyApi } from '@/api/proxy'
import { ElMessage } from 'element-plus'
import { Operation, Promotion, Aim } from '@element-plus/icons-vue'

const proxyStore = useProxyStore()

const modes = [
  { key: 'rule' as const, label: '规则', icon: Operation, description: '根据规则自动选择代理' },
  { key: 'global' as const, label: '全局', icon: Promotion, description: '所有流量使用同一代理' },
  { key: 'direct' as const, label: '直连', icon: Aim, description: '直连网络，不使用代理' },
]

const currentMode = computed(() => proxyStore.mode)

const currentModeDescription = computed(() => {
  const mode = modes.find(m => m.key === currentMode.value)
  return mode?.description || ''
})

const handleModeChange = async (newMode: 'rule' | 'global' | 'direct') => {
  if (newMode === currentMode.value) return

  try {
    await proxyApi.setMode(newMode)
    proxyStore.setMode(newMode)
    ElMessage.success('代理模式已切换')
  } catch (error) {
    console.error('Failed to change mode:', error)
    ElMessage.error('切换模式失败，请稍后重试')
  }
}
</script>

<style scoped>
.clash-mode-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mode-buttons {
  display: flex;
  gap: 8px;
}

.mode-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
  font-size: 13px;
  font-weight: 500;
}

.mode-button:hover:not(.active) {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.mode-button.active {
  background: #3b82f6;
  color: #ffffff;
  border-color: #3b82f6;
}

.mode-button :deep(.el-icon) {
  font-size: 16px;
}

.mode-description {
  text-align: center;
  color: #64748b;
  font-size: 12px;
  padding: 8px 12px;
  background: #f8fafc;
  border-radius: 6px;
}

@media (max-width: 640px) {
  .mode-button {
    padding: 8px 10px;
    font-size: 12px;
  }
}
</style>
