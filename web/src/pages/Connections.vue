<template>
  <div class="connections">
    <div class="connections-header">
      <h2>连接</h2>
      <div class="header-actions">
        <el-tag size="default" type="warning">
          <el-icon class="el-icon--left"><Link /></el-icon>
          {{ activeCount }} 个
        </el-tag>
        <el-button type="danger" size="default" @click="closeAll">关闭所有</el-button>
      </div>
    </div>

    <el-input
      v-model="searchQuery"
      placeholder="搜索连接..."
      clearable
      class="search-bar"
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
    </el-input>

    <el-empty v-if="filteredConnections.length === 0" description="暂无连接" />

    <div v-else class="connection-list">
      <div
        v-for="conn in filteredConnections"
        :key="conn.id"
        class="connection-item"
        :class="{ 'is-disconnected': conn.disconnected }"
      >
        <div class="conn-main">
          <div class="conn-host">
            <span class="conn-time" v-if="conn.start">{{ formatStartTime(conn.start) }}</span><span class="conn-separator"> - </span>{{ conn.metadata?.host || conn.metadata?.process || '未知' }}
            <el-tag v-if="conn.disconnected" type="info" size="small">已断开</el-tag>
          </div>
          <div class="conn-info-row primary">
            {{ conn.metadata?.network || 'TCP' }} / {{ conn.metadata?.type || '-' }} / {{ conn.metadata?.destinationIP || conn.metadata?.host || '-' }}:{{ conn.metadata?.destinationPort || '-' }}
          </div>
          <div class="conn-info-row secondary" v-if="conn.rule || (conn.chains && conn.chains.length > 0)">
            {{ conn.rule || '-' }}<span v-if="conn.rulePayload"> / {{ conn.rulePayload }}</span> => {{ (conn.chains || []).slice().reverse().join(' → ') || '-' }}
          </div>
        </div>
        <div class="conn-actions">
          <div class="conn-stats">
            <span class="stat-item upload">
              <span class="stat-label">↑</span>
              {{ formatBytes(conn.upload || 0) }}
            </span>
            <span class="stat-item download">
              <span class="stat-label">↓</span>
              {{ formatBytes(conn.download || 0) }}
            </span>
          </div>
          <el-button
            v-if="!conn.disconnected"
            type="danger"
            size="small"
            plain
            class="close-btn"
            @click="closeConnection(conn.id)"
          >
            关闭
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Search, Link } from '@element-plus/icons-vue'
import { connectionApi } from '@/api/connection'
import { useWebSocket } from '@/composables/useWebSocket'

const { connections } = useWebSocket()

const searchQuery = ref('')

const activeCount = computed(() => {
  return connections.value?.connections?.length || 0
})

const filteredConnections = computed(() => {
  const active = connections.value?.connections || []
  const history = connections.value?.history || []
  const combined = [...active, ...history]
  if (!searchQuery.value) return combined
  const query = searchQuery.value.toLowerCase()
  return combined.filter(conn => {
    const host = conn.metadata?.host || conn.metadata?.process || ''
    const chain = (conn.chains || []).join(' ')
    const rule = conn.rule || ''
    return host.toLowerCase().includes(query) ||
      chain.toLowerCase().includes(query) ||
      rule.toLowerCase().includes(query)
  })
})

onMounted(() => {
})

const closeConnection = async (id: string) => {
  await connectionApi.close(id)
}

const closeAll = async () => {
  await connectionApi.closeAll()
}

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatStartTime = (isoString: string) => {
  const date = new Date(isoString)
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  const ms = String(date.getMilliseconds()).padStart(3, '0')
  return `${hour}:${minute}:${second}.${ms}`
}
</script>

<style scoped>
.connections {
  padding: 0;
  display: flex;
  flex-direction: column;
}

.connections-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

h2 {
  margin: 0;
  color: #3b82f6;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-bar {
  margin-bottom: 12px;
}

.connection-list {
  overflow: visible;
}

.connection-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  margin-bottom: 8px;
  gap: 16px;
}

.connection-item:hover {
  background: #f8fafc;
}

.connection-item.is-disconnected {
  opacity: 0.7;
  background: #f8fafc;
}

.connection-item.is-disconnected:hover {
  background: #f1f5f9;
}

.conn-main {
  flex: 1;
  min-width: 0;
}

.conn-host {
  font-weight: 500;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.conn-host {
  font-weight: 400;
  color: #1e293b;
  font-size: 14px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 8px;
  display: block;
}

.conn-time {
  font-size: 12px;
  font-weight: 400;
  color: #64748b;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
}

.conn-separator {
  color: #64748b;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
}

.conn-host > .el-tag {
  vertical-align: baseline;
  margin-left: 8px;
}

.conn-info-row {
  display: inline-block;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  margin-bottom: 4px;
}

.conn-info-row.secondary {
  display: block;
  width: fit-content;
}

.conn-info-row.primary {
  background: #f1f5f9;
  color: #475569;
}

.conn-info-row.secondary {
  background: #ecfdf5;
  color: #047857;
}

.conn-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.conn-stats {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-item {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 3px;
}

.stat-item.upload {
  color: #d97706;
  background: #fffbeb;
}

.stat-item.download {
  color: #2563eb;
  background: #eff6ff;
}

.stat-label {
  margin-right: 2px;
}

.close-btn {
  flex-shrink: 0;
}

@media (max-width: 640px) {
  .connection-item {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .conn-actions {
    justify-content: space-between;
  }

  .close-btn {
    padding: 4px 10px;
    font-size: 12px;
    height: 28px;
  }
}
</style>