<template>
  <div class="proxy-navigator">
    <div
      v-for="group in groups"
      :key="group"
      class="nav-item"
      :class="{ active: activeGroup === group }"
      @click="$emit('select', group)"
    >
      {{ group }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  groups: string[]
  activeGroup?: string
}>()

defineEmits<{
  select: [group: string]
}>()
</script>

<style scoped>
.proxy-navigator {
  display: flex;
  gap: 8px;
  padding: 12px 20px;
  overflow-x: auto;
  background: #1a1a2e;
  border-bottom: 1px solid #16213e;
}

.proxy-navigator::-webkit-scrollbar {
  display: none;
}

.nav-item {
  padding: 6px 16px;
  background: #16213e;
  border-radius: 16px;
  font-size: 13px;
  color: #8b9dc3;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: #0f3460;
  color: #fff;
}

.nav-item.active {
  background: #e94560;
  color: #fff;
}
</style>