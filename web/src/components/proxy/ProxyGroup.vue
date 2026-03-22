<template>
  <el-card class="proxy-group">
    <template #header>
      <div class="group-header" @click="$emit('toggle')">
        <div class="group-info">
          <span class="group-name">{{ group.name }}</span>
          <el-tag size="small">{{ group.type }}</el-tag>
        </div>
        <div class="group-meta">
          <el-tag type="info" size="small">{{ group.all.length }} 节点</el-tag>
          <span class="group-current">{{ group.now }}</span>
          <el-icon class="expand-icon" :class="{ expanded: open }">
            <ArrowDown />
          </el-icon>
        </div>
      </div>
    </template>
    
    <el-collapse-transition>
      <div v-show="open" class="group-content">
        <ProxyHead
          :group-name="group.name"
          :initial-filter-text="filterText"
          :initial-sort-type="sortType"
          @update-filter="onFilterChange"
          @update-sort="onSortChange"
          @locate="$emit('locate')"
          @check-all="$emit('checkAll')"
        />
        <ProxyGrid
          :proxies="filteredProxies"
          :columns="columns"
          :selected-proxy="group.now"
          @select="onSelectProxy"
        />
      </div>
    </el-collapse-transition>
  </el-card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import type { ProxyGroup, Proxy } from '@/types/api'
import type { SortType } from '@/stores/proxy'
import ProxyHead from './ProxyHead.vue'
import ProxyGrid from './ProxyGrid.vue'
import { useProxyFilter } from '@/composables/useProxyFilter'

const props = defineProps<{
  group: ProxyGroup & { name: string }
  open: boolean
  columns: number
  filterText?: string
  sortType?: SortType
}>()

const emit = defineEmits<{
  toggle: []
  locate: []
  checkAll: []
  updateFilter: [text: string]
  updateSort: [type: SortType]
  selectProxy: [name: string]
}>()

const { filterSort } = useProxyFilter()

const currentFilterText = ref(props.filterText || '')
const currentSortType = ref<SortType>(props.sortType || 'default')

const filteredProxies = computed(() => {
  return filterSort(
    props.group.all,
    currentFilterText.value,
    currentSortType.value
  )
})

const onSelectProxy = (proxy: Proxy) => {
  emit('selectProxy', proxy.name)
}

const onFilterChange = (text: string) => {
  currentFilterText.value = text
  emit('updateFilter', text)
}

const onSortChange = (type: SortType) => {
  currentSortType.value = type
  emit('updateSort', type)
}

const groupRefs = ref<Record<string, HTMLElement | null>>({})

defineExpose({
  scrollTo: () => {
    if (groupRefs.value[props.group.name]) {
      groupRefs.value[props.group.name]?.scrollIntoView({ behavior: 'smooth' })
    }
  }
})
</script>

<style scoped>
.proxy-group {
  margin-bottom: 16px;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-name {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.group-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-current {
  font-size: 13px;
  color: #64748b;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.expand-icon {
  color: #64748b;
  transition: transform 0.2s ease;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

.group-content {
  margin-top: 12px;
  padding: 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
}
</style>