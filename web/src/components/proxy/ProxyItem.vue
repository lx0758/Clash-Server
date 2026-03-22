<template>
  <div
    class="proxy-item"
    :class="{ selected }"
    @click="$emit('click')"
  >
    <div v-if="selected" class="selected-badge"></div>
    <div class="proxy-name">{{ proxy.name }}</div>
    <div class="proxy-meta">
      <span class="proxy-type">{{ proxy.type }}</span>
      <span class="proxy-delay" :class="delayClass">
        {{ delayDisplay }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Proxy } from '@/types/api'
import { useProxyFilter } from '@/composables/useProxyFilter'

const props = defineProps<{
  proxy: Proxy
  selected?: boolean
}>()

defineEmits<{
  click: []
}>()

const { formatDelay } = useProxyFilter()

const delayText = computed(() => formatDelay(props.proxy.delay))

const delayDisplay = computed(() => {
  if (props.proxy.delay === 0) return '超时'
  return delayText.value
})

const delayClass = computed(() => {
  const delay = props.proxy.delay
  if (delay === 0) return 'timeout'
  if (delay <= 400) return 'fast'
  if (delay <= 800) return 'medium'
  if (delay <= 1500) return 'slow'
  if (delay <= 3000) return 'very-slow'
  return 'error'
})
</script>

<style scoped>
.proxy-item {
  padding: 12px 16px;
  background: #f8fafc;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e2e8f0;
  min-height: 76px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-sizing: border-box;
  overflow: visible;
  position: relative;
  width: 100%;
}

.proxy-item:hover {
  background: #ffffff;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.proxy-item.selected {
  background: #ffffff;
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.15);
}

.selected-badge {
  position: absolute;
  top: -1px;
  left: -1px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 28px 28px 0 0;
  border-color: #60a5fa transparent transparent transparent;
  z-index: 1;
  filter: drop-shadow(1px 1px 2px rgba(0, 0, 0, 0.15));
}

.proxy-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
  margin-bottom: 8px;
}

.proxy-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.proxy-type {
  font-size: 12px;
  color: #64748b;
  background: #f1f5f9;
  padding: 2px 8px;
  border-radius: 3px;
  flex-shrink: 0;
}

.proxy-delay {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 3px;
  flex-shrink: 0;
}

.proxy-delay.fast {
  background: #f0fdf4;
  color: #22c55e;
}

.proxy-delay.medium {
  background: #f7fee7;
  color: #84cc16;
}

.proxy-delay.slow {
  background: #fefce8;
  color: #eab308;
}

.proxy-delay.very-slow {
  background: #fff7ed;
  color: #f97316;
}

.proxy-delay.error {
  background: #fef2f2;
  color: #ef4444;
}

.proxy-delay.timeout {
  background: #f1f5f9;
  color: #94a3b8;
}

@media (max-width: 768px) {
  .proxy-item {
    padding: 10px 12px;
    min-height: 64px;
  }

  .proxy-name {
    font-size: 13px;
    margin-bottom: 6px;
  }

  .proxy-meta {
    gap: 6px;
  }

  .proxy-type {
    font-size: 11px;
    padding: 2px 6px;
  }

  .proxy-delay {
    font-size: 11px;
    padding: 2px 6px;
  }
}
</style>