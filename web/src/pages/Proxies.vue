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
                @update-filter="onGlobalFilterChange"
                @update-sort="onGlobalSortChange"
                @check-all="onCheckAllGlobal"
              />
              <ProxyGrid
                v-if="proxyStore.globalGroup"
                :proxies="filteredGlobalProxies"
                :selected-proxy="proxyStore.globalGroup?.now"
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
                    @update-filter="(text: string) => proxyStore.setGroupFilter(group.name, text)"
                    @update-sort="(type: string) => proxyStore.setGroupSort(group.name, type as any)"
                    @check-all="onCheckAll(group.name)"
                  />
                  <ProxyGrid
                    :proxies="getFilteredProxies(group)"
                    :selected-proxy="group.now"
                    @select="(proxy: Proxy) => proxyStore.selectProxy(group.name, proxy.name)"
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
import { onMounted, ref, computed } from 'vue'
import { useProxyStore } from '@/stores/proxy'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useProxyFilter } from '@/composables/useProxyFilter'
import type { SortType } from '@/stores/proxy'
import type { Proxy, ProxyGroup } from '@/types/api'
import ModeSwitcher from '@/components/proxy/ModeSwitcher.vue'
import ProxyHead from '@/components/proxy/ProxyHead.vue'
import ProxyGrid from '@/components/proxy/ProxyGrid.vue'

const proxyStore = useProxyStore()
const { isMobile } = useBreakpoint()
const { filterSort } = useProxyFilter()

const activeGroups = ref<string[]>([])

const globalFilterText = ref('')
const globalSortType = ref<SortType>('default')

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

const onModeChange = async () => {
  await proxyStore.setMode(proxyStore.mode)
}

const onCheckAll = async (groupName: string) => {
  try {
    await proxyStore.checkGroupDelay(groupName)
  } catch (e) {
    console.error('Failed to check delay:', e)
  }
}

const onGlobalFilterChange = (text: string) => {
  globalFilterText.value = text
}

const onGlobalSortChange = (type: SortType) => {
  globalSortType.value = type
}

const onCheckAllGlobal = async () => {
  try {
    await proxyStore.checkGroupDelay('GLOBAL')
  } catch (e) {
    console.error('Failed to check delay:', e)
  }
}

const onSelectGlobalProxy = async (proxy: Proxy) => {
  await proxyStore.selectProxy('GLOBAL', proxy.name)
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