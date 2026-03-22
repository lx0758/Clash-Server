<template>
  <div class="traffic-info">
    <div class="traffic-header">
      <span class="traffic-label">流量:</span>
      <span class="traffic-value">
        {{ formatBytes(uploadUsed + downloadUsed) }} / {{ formatBytes(totalTransfer) }}
        <span class="percent">({{ percent }}%)</span>
      </span>
    </div>
    <el-progress
      :percentage="percent"
      :status="progressStatus"
      :stroke-width="6"
      :show-text="false"
    />
    <div class="traffic-details">
      <span class="detail-item">
        <span class="arrow up">↑</span>
        {{ formatBytes(uploadUsed) }}
      </span>
      <span class="detail-item">
        <span class="arrow down">↓</span>
        {{ formatBytes(downloadUsed) }}
      </span>
      <span v-if="expireAt" class="detail-item expire" :class="expireClass">
        过期: {{ formatDate(expireAt) }}
        <span v-if="daysLeft !== null && daysLeft <= 7 && daysLeft > 0" class="warning">
          ({{ daysLeft }}天)
        </span>
        <span v-if="isExpired" class="expired">(已过期)</span>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatBytes, formatDaysLeft, isExpired as checkExpired } from '@/utils/format'

const props = defineProps<{
  uploadUsed: number
  downloadUsed: number
  totalTransfer: number
  expireAt: string | null
}>()

const percent = computed(() => {
  const used = props.uploadUsed + props.downloadUsed
  if (props.totalTransfer === 0) return 0
  return Math.round((used / props.totalTransfer) * 100)
})

const daysLeft = computed(() => formatDaysLeft(props.expireAt))

const isExpired = computed(() => checkExpired(props.expireAt))

const progressStatus = computed(() => {
  if (percent.value >= 90) return 'exception'
  if (percent.value >= 70) return 'warning'
  return 'success'
})

const expireClass = computed(() => {
  if (isExpired.value) return 'expired'
  if (daysLeft.value !== null && daysLeft.value <= 7) return 'warning'
  return ''
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.traffic-info {
  background: #f8fafc;
  padding: 12px;
  border-radius: 6px;
}

.traffic-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.traffic-label {
  color: #64748b;
}

.traffic-value {
  color: #1e293b;
}

.traffic-value .percent {
  color: #64748b;
  margin-left: 4px;
}

.traffic-details {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  font-size: 12px;
  color: #64748b;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.arrow {
  font-size: 10px;
}

.arrow.up {
  color: #f59e0b;
}

.arrow.down {
  color: #10b981;
}

.detail-item.expire {
  margin-left: auto;
}

.detail-item.expire.warning {
  color: #d97706;
}

.detail-item.expire.expired {
  color: #dc2626;
}

.warning {
  color: #d97706;
}

.expired {
  color: #dc2626;
}
</style>