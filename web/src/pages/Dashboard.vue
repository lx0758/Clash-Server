<template>
  <div class="dashboard">
    <el-alert
      v-if="systemStore.coreStatus.error"
      type="error"
      :title="'Core 错误: ' + systemStore.coreStatus.error"
      :closable="false"
      show-icon
      class="error-banner"
    />

    <el-row :gutter="16" class="row-one">
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="mode-card">
          <template #header>
            <div class="card-header">
              <el-icon class="card-icon"><Operation /></el-icon>
              <span class="card-title">代理模式</span>
            </div>
          </template>
          <ClashModeCard />
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card class="subscription-card">
          <template #header>
            <div class="card-header">
              <el-icon class="card-icon"><FolderOpened /></el-icon>
              <span class="card-title">订阅信息</span>
            </div>
          </template>
          <ActiveSubscriptionCard />
        </el-card>
      </el-col>
    </el-row>

    <el-card class="traffic-card">
      <template #header>
        <div class="card-header">
          <el-icon class="card-icon"><TrendCharts /></el-icon>
          <span class="card-title">流量信息</span>
        </div>
      </template>
      <div class="traffic-content">
        <div class="chart-section">
          <TrafficGraph :compact="isMobile" />
        </div>
        <el-row :gutter="12" class="stats-row">
          <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
            <StatCard icon="Top" title="上行速度" :value="traffic.up" unit="/s" color="secondary" :compact="isMobile" />
          </el-col>
          <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
            <StatCard icon="Bottom" title="下行速度" :value="traffic.down" unit="/s" color="primary" :compact="isMobile" />
          </el-col>
          <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
            <StatCard icon="Upload" title="上传流量" :value="connections.uploadTotal" color="secondary" :compact="isMobile" />
          </el-col>
          <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
            <StatCard icon="Download" title="下载流量" :value="connections.downloadTotal" color="primary" :compact="isMobile" />
          </el-col>
          <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6">
            <StatCard icon="Link" title="活跃连接" :value="connections.connections?.length ?? 0" color="success" :compact="isMobile" raw unit="个" />
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card class="runtime-card">
      <template #header>
        <div class="card-header">
          <el-icon class="card-icon"><Monitor /></el-icon>
          <span class="card-title">运行信息</span>
        </div>
      </template>
      <RuntimeInfoCard />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'
import { useSystemStore } from '@/stores/system'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { Operation, FolderOpened, TrendCharts, Monitor } from '@element-plus/icons-vue'
import ClashModeCard from '@/components/dashboard/ClashModeCard.vue'
import TrafficGraph from '@/components/dashboard/TrafficGraph.vue'
import StatCard from '@/components/dashboard/StatCard.vue'
import ActiveSubscriptionCard from '@/components/dashboard/ActiveSubscriptionCard.vue'
import RuntimeInfoCard from '@/components/dashboard/RuntimeInfoCard.vue'

const systemStore = useSystemStore()
const { traffic, connections } = useWebSocket()
const { isMobile } = useBreakpoint()

onMounted(() => {
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.error-banner {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-icon {
  font-size: 18px;
  color: #3b82f6;
}

.card-title {
  font-weight: 600;
  color: #1e293b;
}

.row-one .el-card {
  height: 100%;
}

.mode-card :deep(.el-card__body),
.subscription-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 100px;
}

@media (max-width: 991px) {
  .row-one .el-card {
    height: auto;
  }
}

.row-one .el-col {
  margin-bottom: 16px;
}

@media (min-width: 992px) {
  .row-one .el-col {
    margin-bottom: 0;
  }
}

.mode-card :deep(.el-card__body),
.subscription-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 100px;
}

.traffic-card {
  width: 100%;
}

.traffic-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chart-section {
  width: 100%;
}

.stats-row {
  width: 100%;
}

.stats-row .el-col {
  margin-bottom: 12px;
}

.runtime-card :deep(.el-card__body) {
  padding: 16px;
}

@media (max-width: 768px) {
  .dashboard {
    gap: 12px;
  }

  .row-one .el-col:last-child {
    margin-bottom: 0;
  }
}
</style>
