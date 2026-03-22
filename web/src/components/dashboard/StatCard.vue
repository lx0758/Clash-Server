<template>
  <el-card
    class="stat-card"
    :class="{ 'compact': compact }"
    shadow="hover"
    @click="handleClick"
  >
    <div class="stat-content">
      <div class="stat-icon" :style="{ backgroundColor: iconColor }">
        <el-icon>
          <component :is="iconComponent" />
        </el-icon>
      </div>
      <div class="stat-info">
        <div class="stat-label">{{ title }}</div>
        <div class="stat-value-container">
          <span class="stat-value">{{ formattedValue }}</span>
          <span class="stat-unit">{{ formattedUnit }}</span>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import * as ElementPlusIcons from '@element-plus/icons-vue'
import { formatBytes } from '@/utils/format'

interface Props {
  icon: string
  title: string
  value: number
  unit?: string
  color?: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info'
  compact?: boolean
  raw?: boolean
  onClick?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  unit: '',
  color: 'primary',
  compact: false,
  raw: false,
  onClick: undefined,
})

const iconComponent = computed(() => {
  return (ElementPlusIcons as any)[props.icon] || ElementPlusIcons.Document
})

const iconColor = computed(() => {
  const colors = {
    primary: 'rgba(59, 130, 246, 0.1)',
    secondary: 'rgba(245, 158, 11, 0.1)',
    success: 'rgba(16, 185, 129, 0.1)',
    error: 'rgba(220, 38, 38, 0.1)',
    warning: 'rgba(217, 119, 6, 0.1)',
    info: 'rgba(6, 182, 212, 0.1)',
  }
  return colors[props.color]
})

const iconSolidColor = computed(() => {
  const colors = {
    primary: '#3b82f6',
    secondary: '#f59e0b',
    success: '#10b981',
    error: '#dc2626',
    warning: '#d97706',
    info: '#06b6d4',
  }
  return colors[props.color]
})

const formattedValue = computed(() => {
  if (props.raw) {
    return props.value.toLocaleString()
  }
  const result = formatBytes(props.value)
  const parts = result.split(' ')
  return parts[0] || '0'
})

const formattedUnit = computed(() => {
  if (props.raw) {
    return props.unit || ''
  }
  const result = formatBytes(props.value)
  const parts = result.split(' ')
  const bytesUnit = parts[1] || 'B'
  return props.unit ? `${bytesUnit}${props.unit}` : bytesUnit
})

const handleClick = () => {
  if (props.onClick) {
    props.onClick()
  }
}
</script>

<style scoped>
.stat-card {
  height: 100%;
  min-height: 72px;
  transition: all 0.2s ease-in-out;
  border: 1px solid #e2e8f0;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.stat-card:deep(.el-card__body) {
  padding: 14px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  flex-shrink: 0;
}

.stat-icon :deep(.el-icon) {
  font-size: 22px;
  color: v-bind(iconSolidColor);
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 4px;
}

.stat-value-container {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-value {
  font-size: 22px;
  font-weight: 600;
  color: #1e293b;
  line-height: 1;
}

.stat-unit {
  font-size: 12px;
  color: #64748b;
  font-weight: 400;
}

.stat-card.compact {
  height: auto;
}

.stat-card.compact:deep(.el-card__body) {
  padding: 8px;
}

.stat-card.compact .stat-icon {
  width: 32px;
  height: 32px;
}

.stat-card.compact .stat-icon :deep(.el-icon) {
  font-size: 16px;
}

.stat-card.compact .stat-value {
  font-size: 16px;
}

@media (max-width: 640px) {
  .stat-card {
    height: 100%;
    min-height: 64px;
  }

  .stat-card:deep(.el-card__body) {
    padding: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .stat-content {
    width: 100%;
  }

  .stat-icon {
    width: 36px;
    height: 36px;
  }

  .stat-icon :deep(.el-icon) {
    font-size: 18px;
  }

  .stat-value {
    font-size: 18px;
  }

  .stat-label {
    font-size: 11px;
  }
}
</style>
