<template>
  <el-dialog
    v-model="visible"
    :title="readonly ? '查看源文件' : '编辑内容'"
    :width="isMobile ? '100%' : '700px'"
    :fullscreen="isMobile"
    :close-on-click-modal="false"
    @close="$emit('close')"
  >
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    <template v-else>
      <div class="header-actions">
        <el-button size="small" @click="copyContent">复制</el-button>
        <el-button size="small" @click="downloadContent">下载</el-button>
        <span class="content-info">
          <span class="char-count">{{ content.length }} 字符</span>
          <span v-if="readonly" class="readonly-hint">只读</span>
        </span>
      </div>
      <el-input
        v-model="content"
        type="textarea"
        :rows="isMobile ? undefined : 20"
        :readonly="readonly"
        placeholder="暂无内容"
        class="content-editor"
      />
    </template>

    <template #footer>
      <el-button v-if="!readonly" type="primary" :loading="saving" @click="handleSave">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { subscriptionApi } from '@/api/subscription'
import { useBreakpoint } from '@/composables/useBreakpoint'
import type { Subscription } from '@/types/api'

const props = defineProps<{
  subscription: Subscription
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const { isMobile } = useBreakpoint()

const visible = ref(true)
const loading = ref(true)
const saving = ref(false)
const content = ref('')

const readonly = props.subscription.source_type === 'remote'

const fetchContent = async () => {
  loading.value = true
  try {
    const res = await subscriptionApi.getContent(props.subscription.id)
    content.value = res.data?.data?.content || ''
  } catch {
    content.value = ''
  } finally {
    loading.value = false
  }
}

const copyContent = async () => {
  try {
    await navigator.clipboard.writeText(content.value)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const downloadContent = () => {
  const blob = new Blob([content.value], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.subscription.name}.yaml`
  a.click()
  URL.revokeObjectURL(url)
}

const handleSave = async () => {
  saving.value = true
  try {
    await subscriptionApi.updateContent(props.subscription.id, content.value)
    emit('saved')
    emit('close')
  } finally {
    saving.value = false
  }
}

watch(() => props.subscription, () => {
  fetchContent()
}, { immediate: true })
</script>

<style scoped>
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 60px 0;
  color: #94a3b8;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.content-info {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: #94a3b8;
}

.readonly-hint {
  color: #f59e0b;
}

.content-editor :deep(textarea) {
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre;
  overflow-x: auto;
}
</style>

<style>
@media (max-width: 768px) {
  .el-dialog.is-fullscreen .content-editor {
    flex: 1;
    min-height: 0;
  }

  .el-dialog.is-fullscreen .content-editor .el-textarea__inner {
    height: 100% !important;
    min-height: 100% !important;
    resize: none !important;
    border-radius: 0 !important;
    background: #f8fafc !important;
    white-space: pre !important;
    overflow-x: auto !important;
  }
}
</style>
