<template>
  <div class="proxy-head" :class="{ mobile: isMobile }">
    <el-input
      v-model="filterText"
      placeholder="筛选节点..."
      clearable
      size="default"
      class="filter-input"
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
    </el-input>
    <el-select v-model="sortType" size="default" class="sort-select">
      <el-option label="默认" value="default" />
      <el-option label="延迟升序" value="delay-asc" />
      <el-option label="延迟降序" value="delay-desc" />
      <el-option label="名称" value="name" />
    </el-select>
    <el-button type="primary" size="default" @click="$emit('checkAll')">
      测延迟
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'
import type { SortType } from '@/stores/proxy'

const props = defineProps<{
  groupName: string
  initialFilterText?: string
  initialSortType?: SortType
  isMobile?: boolean
}>()

const emit = defineEmits<{
  updateFilter: [text: string]
  updateSort: [type: SortType]
  locate: []
  checkAll: []
}>()

const filterText = ref(props.initialFilterText || '')
const sortType = ref<SortType>(props.initialSortType || 'default')

watch(() => props.initialFilterText, (v) => {
  filterText.value = v || ''
})

watch(() => props.initialSortType, (v) => {
  sortType.value = v || 'default'
})

watch(filterText, (v) => {
  emit('updateFilter', v)
})

watch(sortType, (v) => {
  emit('updateSort', v)
})
</script>

<style scoped>
.proxy-head {
  display: flex;
  gap: 8px;
  padding: 8px 0;
}

.filter-input {
  flex: 1;
  min-width: 0;
}

.sort-select {
  width: 120px;
  flex-shrink: 0;
}

.proxy-head.mobile {
  flex-wrap: wrap;
}

.proxy-head.mobile .filter-input {
  width: 100%;
}

.proxy-head.mobile .sort-select {
  flex: 1;
}
</style>