<template>
  <div
    class="traffic-indicator"
    @click="handleClick"
    tabindex="0"
    role="button"
    :aria-label="`实时流量：上行 ${uploadValue}，下行 ${downloadValue}`"
  >
    <div class="traffic-item upload">
      <el-icon><Top /></el-icon>
      <span class="traffic-value">{{ uploadValue }}</span>
    </div>
    <div class="traffic-item download">
      <el-icon><Bottom /></el-icon>
      <span class="traffic-value">{{ downloadValue }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Top, Bottom } from '@element-plus/icons-vue'
import { useWebSocket } from '@/composables/useWebSocket'
import { formatTraffic } from '@/utils/format'

const router = useRouter()
const { traffic } = useWebSocket()

const uploadValue = computed(() => formatTraffic(traffic.value.up, true))
const downloadValue = computed(() => formatTraffic(traffic.value.down, true))

const handleClick = () => {
  router.push('/')
}
</script>

<style scoped>
.traffic-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  background: #f1f5f9;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  border: none;
  outline: none;
}

.traffic-indicator:hover {
  background: #e2e8f0;
}

.traffic-indicator:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

.traffic-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #64748b;
}

.traffic-item.upload {
  color: #f59e0b;
}

.traffic-item.download {
  color: #3b82f6;
}

.traffic-value {
  font-weight: 600;
}

@media (max-width: 360px) {
  .traffic-indicator {
    gap: 8px;
    padding: 2px 6px;
  }

  .traffic-value {
    font-size: 11px;
  }
}
</style>
