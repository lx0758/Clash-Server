<template>
  <div class="proxy-grid">
    <ProxyItem
      v-for="proxy in proxies"
      :key="proxy.name"
      :proxy="proxy"
      :selected="proxy.name === selectedProxy"
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