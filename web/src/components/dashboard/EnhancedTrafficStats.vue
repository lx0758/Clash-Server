<template>
  <div class="enhanced-traffic-stats">
    <TrafficGraph :compact="isMobile" />

    <el-row :gutter="gutter">
      <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="8">
        <StatCard
          icon="Top"
          title="上行速度"
          :value="traffic.up"
          unit="/s"
          color="secondary"
          :compact="isMobile"
        />
      </el-col>

      <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="8">
        <StatCard
          icon="Bottom"
          title="下行速度"
          :value="traffic.down"
          unit="/s"
          color="primary"
          :compact="isMobile"
        />
      </el-col>

      <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="8">
        <StatCard
          icon="Link"
          title="活跃连接"
          :value="connections.connections.length"
          color="success"
          :compact="isMobile"
        />
      </el-col>

      <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="8">
        <el-card>
          <template #header>
            <span class="card-title">流量统计</span>
          </template>
          <div class="traffic-cards">
            <StatCard
              icon="CloudUpload"
              title="上传量"
              :value="connections.uploadTotal"
              color="secondary"
              :compact="isMobile"
            />
            <StatCard
              icon="CloudDownload"
              title="下载量"
              :value="connections.downloadTotal"
              color="primary"
              :compact="isMobile"
            />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="12" :lg="8" :xl="8">
        <StatCard
          icon="MemoryStick"
          title="内存使用"
          :value="memory.inuse"
          color="error"
          :compact="isMobile"
        />
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useWebSocket } from '@/composables/useWebSocket'
import TrafficGraph from './TrafficGraph.vue'
import StatCard from './StatCard.vue'

const { isMobile } = useBreakpoint()
const { traffic, connections, memory } = useWebSocket()

const gutter = computed(() => {
  return isMobile.value ? 8 : 16
})
</script>

<style scoped>
.enhanced-traffic-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.enhanced-traffic-stats > .el-row:first-child {
  margin-top: 0;
}

.traffic-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
</style>
