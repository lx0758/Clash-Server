<template>
  <div class="proxies">
    <div class="proxies-header">
      <h2>代理</h2>
      <ModeSwitcher v-model="proxyStore.mode" @update:model-value="onModeChange" />
    </div>

    <el-skeleton :loading="proxyStore.loading" animated>
      <template #default>
        <el-empty v-if="proxyStore.mode === 'direct'" description="当前为直连模式" />

        <template v-else-if="proxyStore.mode === 'global'">
          <div class="global-group">
            <div class="global-header">
              <div class="group-left">
                <span class="group-name">GLOBAL</span>
                <el-tag size="small" type="info">Selector</el-tag>
                <span class="group-count">{{ proxyStore.globalGroup?.all.length || 0 }}</span>
              </div>
              <div class="group-right">
                <span class="group-current">{{ proxyStore.globalGroup?.now }}</span>
              </div>
            </div>
            <div class="group-content global">
              <ProxyHead
                group-name="GLOBAL"
                :initial-filter-text="globalFilterText"
                :initial-sort-type="globalSortType"
                :is-mobile="isMobile"
                :checking="isGroupChecking('GLOBAL')"
                @update-filter="onGlobalFilterChange"
                @update-sort="onGlobalSortChange"
                @check-all="onCheckAllGlobal"
              />
              <ProxyGrid
                v-if="proxyStore.globalGroup"
                :proxies="filteredGlobalProxies"
                :selected-proxy="proxyStore.globalGroup?.now"
                :testing-names="getTestingNames('GLOBAL')"
                @select="onSelectGlobalProxy"
              />
            </div>
          </div>
        </template>

        <template v-else>
          <div class="proxy-groups">
            <el-collapse v-model="activeGroups">
              <el-collapse-item
                v-for="group in proxyStore.proxyGroups"
                :key="group.name"
                :name="group.name"
              >
                <template #title>
                  <div class="group-title">
                    <div class="group-left">
                      <span class="group-name">{{ group.name }}</span>
                      <el-tag size="small" type="info">{{ group.type }}</el-tag>
                      <span class="group-count">{{ group.all.length }}</span>
                    </div>
                    <div class="group-right">
                      <span class="group-current">{{ group.now }}</span>
                    </div>
                  </div>
                </template>
                <div class="group-content">
                  <ProxyHead
                    :group-name="group.name"
                    :initial-filter-text="proxyStore.getGroupState(group.name).filterText"
                    :initial-sort-type="proxyStore.getGroupState(group.name).sortType"
                    :is-mobile="isMobile"
                    :checking="isGroupChecking(group.name)"
                    @update-filter="(text: string) => proxyStore.setGroupFilter(group.name, text)"
                    @update-sort="(type: string) => proxyStore.setGroupSort(group.name, type as any)"
                    @check-all="onCheckAll(group.name)"
                  />
                  <ProxyGrid
                    :proxies="getFilteredProxies(group)"
                    :selected-proxy="group.now"
                    :testing-names="getTestingNames(group.name)"
                    @select="(proxy: Proxy) => onSelectProxy(group.name, proxy.name)"
                  />
                </div>
              </el-collapse-item>
            </el-collapse>
          </div>
        </template>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { useProxyStore } from '@/stores/proxy'
import { usePendingStore } from '@/stores/pending'
import { useToastStore } from '@/stores/toast'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useProxyFilter } from '@/composables/useProxyFilter'
import { useWebSocket } from '@/composables/useWebSocket'
import type { SortType } from '@/stores/proxy'
import type { Proxy, ProxyGroup } from '@/types/api'
import ModeSwitcher from '@/components/proxy/ModeSwitcher.vue'
import ProxyHead from '@/components/proxy/ProxyHead.vue'
import ProxyGrid from '@/components/proxy/ProxyGrid.vue'

const proxyStore = useProxyStore()
const pendingStore = usePendingStore()
const toastStore = useToastStore()
const { isMobile } = useBreakpoint()
const { filterSort } = useProxyFilter()
const { onDelay } = useWebSocket()

const activeGroups = ref<string[]>([])

const globalFilterText = ref('')
const globalSortType = ref<SortType>('default')

// 正在测延迟的节点名称, 按组归集
const testingNodes = ref<Record<string, string[]>>({})

// 测延迟超时兜底: 超过该时间未收到完成信号, 强制收敛状态
const DELAY_TIMEOUT_MS = 30000
const delayTimers = new Map<string, ReturnType<typeof setTimeout>>()

const groupKey = (groupName: string) => `delay:${groupName}`

const isGroupChecking = (groupName: string) => pendingStore.isPending(groupKey(groupName))

const getTestingNames = (groupName: string) => testingNodes.value[groupName] || []

const clearGroupTesting = (groupName: string) => {
  testingNodes.value = { ...testingNodes.value, [groupName]: [] }
  const timer = delayTimers.get(groupName)
  if (timer) {
    clearTimeout(timer)
    delayTimers.delete(groupName)
  }
}

/** 启动某组的逐节点测延迟 */
const startDelayTest = async (groupName: string) => {
  if (pendingStore.isPending(groupKey(groupName))) return

  try {
    const res = await proxyStore.startGroupDelayTest(groupName)
    const nodes = res.data.data.nodes || []
    // 立即标记该组全部节点为测速中, 各节点结果到达后逐个收敛
    testingNodes.value = { ...testingNodes.value, [groupName]: [...nodes] }
    pendingStore.start(groupKey(groupName))

    // 超时兜底: 未收到完成信号时强制收敛, 避免节点永久转圈
    const timer = setTimeout(() => {
      if (pendingStore.isPending(groupKey(groupName))) {
        clearGroupTesting(groupName)
        pendingStore.finish(groupKey(groupName))
        toastStore.warning('测延迟超时，部分节点未返回结果')
      }
    }, DELAY_TIMEOUT_MS)
    delayTimers.set(groupName, timer)
  } catch (e) {
    toastStore.error(e instanceof Error ? e.message : '启动测延迟失败')
  }
}

const filteredGlobalProxies = computed(() => {
  const group = proxyStore.globalGroup
  if (!group) return []
  return filterSort(group.all, globalFilterText.value, globalSortType.value)
})

const getFilteredProxies = (group: ProxyGroup & { name: string }) => {
  const state = proxyStore.getGroupState(group.name)
  return filterSort(group.all, state.filterText, state.sortType)
}

onMounted(async () => {
  await Promise.all([
    proxyStore.fetchProxies(),
    proxyStore.fetchMode(),
  ])

  const groups = proxyStore.proxyGroups
  if (groups.length > 0 && groups[0]) {
    activeGroups.value = [groups[0].name]
  }
})

// 订阅测延迟结果, 各节点到达即更新
const unsubscribeDelay = onDelay(async (result) => {
  if (result.done) {
    // 该组全部完成
    clearGroupTesting(result.group)
    pendingStore.finish(groupKey(result.group))
    // 静默同步核心侧的延迟历史, 不触发整页骨架屏
    try {
      await proxyStore.syncAfterDelayTest()
    } catch {
      // 同步失败不打断用户
    }
    return
  }

  if (!result.name) return
  proxyStore.handleDelayResult(result)
  // 该节点完成, 从进行中列表移除
  const group = result.group
  const remaining = (testingNodes.value[group] || []).filter(n => n !== result.name)
  testingNodes.value = { ...testingNodes.value, [group]: remaining }
})

onUnmounted(() => {
  unsubscribeDelay()
  delayTimers.forEach(timer => clearTimeout(timer))
  delayTimers.clear()
})

const onModeChange = async () => {
  await proxyStore.setMode(proxyStore.mode)
}

const onCheckAll = async (groupName: string) => {
  await startDelayTest(groupName)
}

const onGlobalFilterChange = (text: string) => {
  globalFilterText.value = text
}

const onGlobalSortChange = (type: SortType) => {
  globalSortType.value = type
}

const onCheckAllGlobal = async () => {
  await startDelayTest('GLOBAL')
}

const onSelectGlobalProxy = async (proxy: Proxy) => {
  await onSelectProxy('GLOBAL', proxy.name)
}

/** 切换节点, 失败时回滚并提示 */
const onSelectProxy = async (groupName: string, proxyName: string) => {
  try {
    await proxyStore.selectProxy(groupName, proxyName)
  } catch (e) {
    toastStore.error(e instanceof Error ? e.message : '切换节点失败')
  }
}
</script>

<style scoped>
.proxies {
  padding: 0;
}

.proxies-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

h2 {
  color: #3b82f6;
  margin: 0;
}

.proxy-groups {
  margin-top: 8px;
}

.group-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 12px;
}

.group-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-name {
  font-weight: 600;
  color: #1e293b;
}

.group-count {
  font-size: 13px;
  color: #64748b;
}

.group-current {
  font-size: 12px;
  font-weight: 500;
  color: #3b82f6;
  background: #eff6ff;
  padding: 4px 10px;
  border-radius: 4px;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 8px;
}

:deep(.el-collapse-item) {
  margin-bottom: 12px;
}

:deep(.el-collapse-item__header) {
  background: #ffffff;
  border-radius: 6px;
  padding: 0 16px;
  border: 1px solid #e2e8f0;
  height: 52px;
  transition: border-radius 0.2s ease;
}

:deep(.el-collapse-item__header.is-active) {
  border-radius: 6px 6px 0 0;
  border-bottom: none;
}

:deep(.el-collapse-item__wrap) {
  background: transparent;
  border: none;
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

.group-content {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-top: none;
  border-radius: 0 0 6px 6px;
  padding: 12px;
}

.group-content.global {
  border-radius: 0 0 6px 6px;
  border-top: none;
  margin-bottom: 12px;
}

.global-group {
  margin-bottom: 12px;
}

.global-header {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-bottom: none;
  border-radius: 6px 6px 0 0;
  padding: 0 16px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

@media (max-width: 640px) {
  .group-title {
    gap: 8px;
  }

  .group-left {
    flex-wrap: wrap;
    gap: 6px;
  }

  .group-count {
    display: none;
  }

  .group-current {
    max-width: 80px;
    font-size: 12px;
  }
}
</style>