<template>
  <div class="traffic-stats">
    <div class="chart-section">
      <TrafficGraph :compact="isMobile" />
    </div>

    <div class="cards-section">
      <el-row :gutter="12">
        <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
          <StatCard
            icon="Top"
            title="上行速度"
            :value="traffic.up"
            unit="/s"
            color="secondary"
            :compact="isMobile"
          />
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
          <StatCard
            icon="Bottom"
            title="下行速度"
            :value="traffic.down"
            unit="/s"
            color="primary"
            :compact="isMobile"
          />
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
          <StatCard
            icon="CloudUpload"
            title="上传流量"
            :value="connections.uploadTotal"
            color="secondary"
            :compact="isMobile"
          />
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6" :xl="6">
          <StatCard
            icon="CloudDownload"
            title="下载流量"
            :value="connections.downloadTotal"
            color="primary"
            :compact="isMobile"
          />
        </el-col>

        <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="6">
          <StatCard
            icon="Link"
            title="活跃连接"
            :value="connections.connections.length"
            color="success"
            :compact="isMobile"
          />
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useWebSocket } from '@/composables/useWebSocket'
import TrafficGraph from './TrafficGraph.vue'
import StatCard from './StatCard.vue'

const { isMobile } = useBreakpoint()
const { traffic, connections } = useWebSocket()
</script>

<style scoped>
.traffic-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
}

.chart-section {
  flex: 0 0 1;
}

.cards-section {
  flex: 1;
  overflow: hidden;
}

@media (max-width: 768px) {
  .cards-section {
    display: flex;
    flex-direction: column;
  }

  .cards-section > .el-row {
    flex: 0;
  }
}
</style>
