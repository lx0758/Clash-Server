<template>
  <div class="rules">
    <div class="rules-header">
      <h2>规则</h2>
      <el-tag size="default">{{ filteredRules.length }} 条</el-tag>
    </div>

    <el-input
      v-model="searchQuery"
      placeholder="搜索规则..."
      clearable
      class="search-bar"
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
    </el-input>

    <el-skeleton :loading="loading" animated>
      <template #default>
        <el-empty v-if="!coreRunning" description="核心未运行" />
        <el-empty v-else-if="filteredRules.length === 0" description="无匹配规则" />

        <div v-else class="rule-list">
          <div
            v-for="rule in filteredRules"
            :key="rule.index"
            class="rule-item"
          >
            <div class="rule-main">
              <el-tag size="small" type="info">{{ rule.type }}</el-tag>
              <span class="rule-proxy">{{ rule.proxy }}</span>
            </div>
            <div class="rule-footer">
              <span class="rule-payload">{{ rule.payload }}</span>
              <el-tag v-if="rule.extra?.hitCount" size="small" type="warning">
                {{ rule.extra.hitCount }} 次
              </el-tag>
            </div>
          </div>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { coreApi } from '@/api/core'
import { useSystemStore } from '@/stores/system'
import type { CoreRule } from '@/types/api'

const systemStore = useSystemStore()

const rules = ref<CoreRule[]>([])
const loading = ref(true)
const searchQuery = ref('')

const coreRunning = computed(() => systemStore.coreStatus?.running ?? false)

const filteredRules = computed(() => {
  if (!searchQuery.value) return rules.value
  const query = searchQuery.value.toLowerCase()
  return rules.value.filter(rule =>
    rule.type.toLowerCase().includes(query) ||
    rule.payload.toLowerCase().includes(query) ||
    rule.proxy.toLowerCase().includes(query)
  )
})

onMounted(async () => {
  await systemStore.fetchSystemInfo()
  if (coreRunning.value) {
    await fetchRules()
  }
  loading.value = false
})

watch(coreRunning, async (running) => {
  if (running && rules.value.length === 0) {
    loading.value = true
    await fetchRules()
    loading.value = false
  }
})

const fetchRules = async () => {
  try {
    const res = await coreApi.getRules()
    rules.value = res.data.data.rules || []
  } catch {
    rules.value = []
  }
}
</script>

<style scoped>
.rules {
  padding: 0;
  display: flex;
  flex-direction: column;
}

.rules-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

h2 {
  margin: 0;
  color: #3b82f6;
}

.search-bar {
  margin-bottom: 12px;
}

.rule-list {
  overflow: visible;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  align-content: start;
}

.rule-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 16px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: default;
  transition: all 0.2s ease;
}

.rule-item:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.rule-main {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rule-proxy {
  font-size: 14px;
  font-weight: 600;
  color: #10b981;
}

.rule-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.rule-payload {
  flex: 1;
  color: #64748b;
  font-size: 13px;
  word-break: break-all;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 1400px) {
  .rule-list {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 1000px) {
  .rule-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .rule-list {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .rule-item {
    padding: 10px 12px;
  }
}
</style>