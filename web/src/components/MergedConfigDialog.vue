<template>
  <el-dialog
    :model-value="true"
    :title="'合并后配置 - ' + subscriptionName"
    :width="isMobile ? '100%' : '900px'"
    :fullscreen="isMobile"
    @close="$emit('close')"
  >
    <div class="header-actions">
      <el-button @click="copyYaml">复制 YAML</el-button>
      <el-button @click="downloadYaml">下载文件</el-button>
    </div>
    <div class="config-viewer">
      <pre><code>{{ yaml }}</code></pre>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { useBreakpoint } from '@/composables/useBreakpoint'

const props = defineProps<{
  subscriptionName: string
  yaml: string
}>()

const emit = defineEmits<{
  close: []
}>()

const { isMobile } = useBreakpoint()

const copyYaml = async () => {
  try {
    await navigator.clipboard.writeText(props.yaml)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const downloadYaml = () => {
  const blob = new Blob([props.yaml], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'config.yaml'
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.config-viewer {
  max-height: 60vh;
  overflow: auto;
  background: #f8fafc;
  border-radius: 6px;
  padding: 16px;
}

.config-viewer pre {
  margin: 0;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre;
}

.config-viewer code {
  color: #334155;
}
</style>

<style>
@media (max-width: 768px) {
  .el-dialog.is-fullscreen {
    height: 100vh !important;
    max-height: 100vh !important;
    margin: 0 !important;
    display: flex !important;
    flex-direction: column !important;
    overflow: hidden !important;
  }

  .el-dialog.is-fullscreen .el-dialog__header {
    flex-shrink: 0 !important;
  }

  .el-dialog.is-fullscreen .el-dialog__body {
    flex: 1 !important;
    display: flex !important;
    flex-direction: column !important;
    padding: 12px 0 0 0 !important;
    min-height: 0 !important;
    overflow: hidden !important;
  }

  .el-dialog.is-fullscreen .header-actions {
    flex-shrink: 0 !important;
    margin-bottom: 12px !important;
    padding: 0 5px !important;
  }

  .el-dialog.is-fullscreen .config-viewer {
    flex: 1 !important;
    min-height: 0 !important;
    max-height: none !important;
    overflow: auto !important;
    padding: 5px !important;
    border-radius: 0 !important;
    background: #f8fafc !important;
  }
}
</style>