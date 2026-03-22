<template>
  <el-radio-group v-model="selectedMode" size="default">
    <el-radio-button
      v-for="m in modes"
      :key="m.value"
      :value="m.value"
    >
      {{ m.label }}
    </el-radio-button>
  </el-radio-group>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ProxyMode } from '@/api/proxy'

const props = defineProps<{
  modelValue: ProxyMode
}>()

const emit = defineEmits<{
  'update:modelValue': [mode: ProxyMode]
}>()

const selectedMode = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const modes: Array<{ value: ProxyMode; label: string }> = [
  { value: 'rule', label: '规则' },
  { value: 'global', label: '全局' },
  { value: 'direct', label: '直连' },
]
</script>

<style scoped>
:deep(.el-radio-button__inner) {
  padding: 10px 20px;
  background: #f1f5f9;
  border-color: #e2e8f0;
  color: #64748b;
  font-weight: 500;
}

:deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: #3b82f6;
  border-color: #3b82f6;
  color: #fff;
}

@media (max-width: 640px) {
  :deep(.el-radio-button__inner) {
    padding: 8px 14px;
    font-size: 13px;
  }
}
</style>