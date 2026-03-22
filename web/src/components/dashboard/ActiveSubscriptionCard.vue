<template>
  <div class="active-subscription">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    <div v-else-if="!activeSubscription" class="empty-state">
      <el-icon><FolderOpened /></el-icon>
      <span>暂无激活的订阅</span>
    </div>
    <template v-else>
      <div class="sub-header">
        <el-icon class="status-icon"><CircleCheckFilled /></el-icon>
        <span class="sub-name">{{ activeSubscription.name }}</span>
      </div>

      <div class="sub-info-grid">
        <div class="info-item">
          <span class="info-value">{{ activeSubscription.node_count }}</span>
          <span class="info-label">节点</span>
        </div>
        <div v-if="activeSubscription.expire_at" class="info-item">
          <span class="info-value" :class="expireClass">{{ formatExpireTime(activeSubscription.expire_at) }}</span>
          <span class="info-label">到期</span>
        </div>
      </div>

      <div v-if="hasTraffic" class="traffic-section">
        <div class="traffic-header">
          <span>流量使用</span>
          <span class="traffic-text">{{ formatBytes(usedTraffic) }} / {{ formatBytes(activeSubscription.total_transfer) }}</span>
        </div>
        <el-progress 
          :percentage="trafficPercentage" 
          :stroke-width="8"
          :show-text="false"
          :color="progressColor"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { subscriptionApi } from '@/api/subscription'
import { formatBytes } from '@/utils/format'
import { Loading, FolderOpened, CircleCheckFilled } from '@element-plus/icons-vue'
import type { Subscription, SubscriptionWithCounts } from '@/types/api'

const loading = ref(true)
const activeSubscription = ref<Subscription | null>(null)

const hasTraffic = computed(() => {
  if (!activeSubscription.value) return false
  const sub = activeSubscription.value
  return sub.total_transfer > 0
})

const usedTraffic = computed(() => {
  if (!activeSubscription.value) return 0
  return activeSubscription.value.upload_used + activeSubscription.value.download_used
})

const trafficPercentage = computed(() => {
  if (!activeSubscription.value || !hasTraffic.value) return 0
  const total = activeSubscription.value.total_transfer
  if (total <= 0) return 0
  return Math.min(Math.round((usedTraffic.value / total) * 100), 100)
})

const progressColor = computed(() => {
  const pct = trafficPercentage.value
  if (pct >= 90) return '#ef4444'
  if (pct >= 70) return '#f59e0b'
  return '#3b82f6'
})

const expireClass = computed(() => {
  if (!activeSubscription.value?.expire_at) return ''
  const expire = new Date(activeSubscription.value.expire_at)
  const now = new Date()
  const diff = expire.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  
  if (diff < 0) return 'expired'
  if (days <= 7) return 'expiring'
  return ''
})

const formatExpireTime = (expireAt: string): string => {
  const expire = new Date(expireAt)
  const now = new Date()
  const diff = expire.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  
  if (diff < 0) return '已过期'
  if (days <= 7) return `${days}天`
  return expire.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

const fetchActiveSubscription = async () => {
  loading.value = true
  try {
    const res = await subscriptionApi.list()
    const items: SubscriptionWithCounts[] = res.data.data?.subscriptions || []
    const activeItem = items.find((item) => item.subscription?.active)
    activeSubscription.value = activeItem?.subscription || null
  } catch (error) {
    console.error('Failed to fetch subscriptions:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchActiveSubscription()
})
</script>

<style scoped>
.active-subscription {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-state,
.empty-state {
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

.empty-state .el-icon {
  font-size: 24px;
}

.sub-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-icon {
  font-size: 18px;
  color: #10b981;
  flex-shrink: 0;
}

.sub-name {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sub-info-grid {
  display: flex;
  gap: 24px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-value {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}

.info-value.expired {
  color: #ef4444;
}

.info-value.expiring {
  color: #f59e0b;
}

.info-label {
  font-size: 12px;
  color: #94a3b8;
}

.traffic-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-top: 8px;
  border-top: 1px solid #f1f5f9;
}

.traffic-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #64748b;
}

.traffic-text {
  color: #1e293b;
  font-weight: 500;
}
</style>
