<template>
  <div class="proxy-grid">
    <ProxyItem
      v-for="proxy in proxies"
      :key="proxy.name"
      :proxy="proxy"
      :selected="proxy.name === selectedProxy"
      :testing="testingNames?.includes(proxy.name)"
      @click="$emit('select', proxy)"
    />
  </div>
</template>

<script setup lang="ts">
import ProxyItem from './ProxyItem.vue'
import type { Proxy } from '@/types/api'

defineProps<{
  proxies: Proxy[]
  columns?: number
  selectedProxy?: string
  /** 正在测延迟的节点名称列表 */
  testingNames?: string[]
}>()

defineEmits<{
  select: [proxy: Proxy]
}>()
</script>

<style scoped>
.proxy-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.proxy-grid > * {
  min-width: 0;
}

@media (max-width: 1200px) {
  .proxy-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 900px) {
  .proxy-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .proxy-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}
</style>