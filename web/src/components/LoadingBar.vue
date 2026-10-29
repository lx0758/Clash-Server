<script setup lang="ts">
import { computed } from 'vue'
import { useLoadingStore } from '@/stores/loading'

const loadingStore = useLoadingStore()

// 仅在存在进行中的请求时显示
const visible = computed(() => loadingStore.count > 0)
</script>

<template>
  <div v-if="visible" class="loading-bar" role="progressbar" aria-label="加载中"></div>
</template>

<style scoped>
.loading-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  z-index: 4000;
  background: linear-gradient(90deg, #3b82f6, #60a5fa, #3b82f6);
  background-size: 200% 100%;
  animation: loading-slide 1.2s linear infinite;
  pointer-events: none;
}

@keyframes loading-slide {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
