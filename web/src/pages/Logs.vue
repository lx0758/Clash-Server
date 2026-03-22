<template>
  <div class="logs">
    <div class="logs-header">
      <h2>日志</h2>
      <div class="header-actions">
        <el-tag size="default">{{ filteredLogs.length }} 条</el-tag>
        <el-button type="danger" size="default" @click="handleClearLogs">清空</el-button>
      </div>
    </div>

    <el-input
      v-model="searchQuery"
      placeholder="搜索日志..."
      clearable
      class="search-bar"
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
    </el-input>

    <el-empty v-if="filteredLogs.length === 0" description="暂无日志" />

    <div v-else class="log-list">
      <div
        v-for="(log, index) in filteredLogs"
        :key="index"
        class="log-item"
        :class="getLogLevelClass(log.type)"
      >
        <span class="log-time-box">{{ log.time }}</span>
        <span class="log-level-box" :class="getLogLevelClass(log.type)">{{ log.type.toUpperCase() }}</span>
        <span class="log-payload">{{ log.payload }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { useWebSocket } from '@/composables/useWebSocket'

const { logs, connect, disconnect, subscribe, unsubscribe, clearLogs } = useWebSocket()

const searchQuery = ref('')

interface LogItem {
  type: string
  payload: string
  time: string
}

const filteredLogs = computed(() => {
  const list = logs.value || []
  const now = new Date()
  const items: LogItem[] = list.map((log, i) => ({
    type: log.type,
    payload: log.payload,
    time: formatTime(new Date(now.getTime() - i * 100))
  }))
  if (!searchQuery.value) return items.slice(0, 500)
  const query = searchQuery.value.toLowerCase()
  return items.filter(log =>
    log.type.toLowerCase().includes(query) ||
    log.payload.toLowerCase().includes(query)
  ).slice(0, 500)
})

const formatTime = (date: Date) => {
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  const ms = String(date.getMilliseconds()).padStart(3, '0')
  return `${hour}:${minute}:${second}.${ms}`
}

const getLogLevelClass = (type: string) => {
  const t = type.toLowerCase()
  if (t === 'error') return 'level-error'
  if (t === 'warning' || t === 'warn') return 'level-warning'
  if (t === 'info') return 'level-info'
  if (t === 'debug') return 'level-debug'
  return ''
}

const handleClearLogs = () => {
  clearLogs()
}

onMounted(() => {
  connect()
  subscribe(['logs'])
})

onUnmounted(() => {
  unsubscribe(['logs'])
  disconnect()
})
</script>

<style scoped>
.logs {
  padding: 0;
  display: flex;
  flex-direction: column;
}

.logs-header {
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

.log-list {
  overflow: visible;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
}

.log-item {
  padding: 8px 16px;
  border-bottom: 1px solid #f1f5f9;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.6;
  transition: background 0.2s;
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 6px;
}

.log-item:last-child {
  border-bottom: none;
}

.log-item:hover {
  background: #f8fafc;
}

.log-item.level-error {
  background: rgba(239, 68, 68, 0.05);
}

.log-item.level-error:hover {
  background: rgba(239, 68, 68, 0.1);
}

.log-item.level-warning {
  background: rgba(245, 158, 11, 0.05);
}

.log-item.level-warning:hover {
  background: rgba(245, 158, 11, 0.1);
}

.log-time-box,
.log-level-box {
  display: inline-block;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  border: 1px solid;
}

.log-time-box {
  color: #64748b;
  background: #f8fafc;
  border-color: #e2e8f0;
}

.log-level-box {
  font-weight: 500;
}

.log-level-box.level-error {
  color: #ef4444;
  background: #fef2f2;
  border-color: #fecaca;
}

.log-level-box.level-warning {
  color: #f59e0b;
  background: #fffbeb;
  border-color: #fed7aa;
}

.log-level-box.level-info {
  color: #3b82f6;
  background: #eff6ff;
  border-color: #bfdbfe;
}

.log-level-box.level-debug {
  color: #8b5cf6;
  background: #f5f3ff;
  border-color: #ddd6fe;
}

.log-payload {
  color: #475569;
  word-break: break-word;
  flex: 1;
  min-width: 0;
}

@media (max-width: 640px) {
  .log-item {
    padding: 6px 12px;
    font-size: 12px;
    flex-wrap: wrap;
    gap: 4px;
  }

  .log-time-box,
  .log-level-box {
    flex-shrink: 0;
  }

  .log-payload {
    flex: none;
    width: 100%;
    margin-top: 2px;
  }
}
</style>