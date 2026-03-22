<template>
  <el-dialog
    :model-value="true"
    :title="'脚本管理 - ' + subscriptionName"
    :width="isMobile ? '100%' : '700px'"
    :fullscreen="isMobile"
    :class="{ 'mobile-dialog': isMobile }"
    @close="$emit('close')"
  >
    <div class="toolbar">
      <el-button type="primary" @click="openAdd">+ 添加脚本</el-button>
    </div>

    <div class="table-wrapper">
      <el-table :data="scripts" style="width: 100%">
        <el-table-column label="名称" prop="name" />
        <el-table-column label="描述" prop="description">
          <template #default="{ row }">
            <span class="desc">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="toggleEnabled(row)">
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="deleteScript(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="scripts.length === 0" description="暂无脚本" />
    </div>

    <template #footer>
      <el-button @click="$emit('close')">关闭</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="showEditor"
    :title="editingScript ? '编辑脚本' : '添加脚本'"
    :width="isMobile ? '100%' : '600px'"
    :fullscreen="isMobile"
    :class="{ 'mobile-dialog': isMobile }"
  >
    <el-form :model="scriptForm" label-position="top">
      <el-form-item label="名称">
        <el-input v-model="scriptForm.name" placeholder="脚本名称" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="scriptForm.description" placeholder="脚本描述" />
      </el-form-item>
      <el-form-item>
        <template #label>
          <span>脚本内容</span>
          <div class="hint">接收全局变量 <code>config</code>，修改后返回即可生效</div>
        </template>
        <el-input
          v-model="scriptForm.content"
          type="textarea"
          :rows="16"
          placeholder="// 示例
function main(config) {
  config.dns = {
    enable: true,
    'enhanced-mode': 'fake-ip'
  };
  return config;
}"
          class="editor"
        />
      </el-form-item>
      <el-form-item>
        <el-switch v-model="scriptForm.enabled" active-text="启用此脚本" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="showEditor = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="saveScript">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { scriptApi } from '@/api/script'
import { useBreakpoint } from '@/composables/useBreakpoint'
import type { Script } from '@/types/api'

const { isMobile } = useBreakpoint()

const props = defineProps<{
  subscriptionId: number
  subscriptionName: string
}>()

const emit = defineEmits<{
  close: []
}>()

const loading = ref(false)
const scripts = ref<Script[]>([])
const showEditor = ref(false)
const editingScript = ref<Script | null>(null)
const scriptForm = reactive({ name: '', description: '', content: '', enabled: true })

const fetchScripts = async () => {
  const res = await scriptApi.list(props.subscriptionId)
  scripts.value = res.data.data.scripts || []
}

onMounted(fetchScripts)
watch(() => props.subscriptionId, fetchScripts)

const openAdd = () => {
  editingScript.value = null
  scriptForm.name = ''
  scriptForm.description = ''
  scriptForm.content = ''
  scriptForm.enabled = true
  showEditor.value = true
}

const openEdit = (script: Script) => {
  editingScript.value = script
  scriptForm.name = script.name
  scriptForm.description = script.description
  scriptForm.content = script.content
  scriptForm.enabled = script.enabled
  showEditor.value = true
}

const toggleEnabled = async (script: Script) => {
  await scriptApi.update(props.subscriptionId, script.id, { enabled: !script.enabled })
  await fetchScripts()
}

const deleteScript = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除此脚本吗？', '删除确认', { type: 'warning' })
    await scriptApi.delete(props.subscriptionId, id)
    await fetchScripts()
  } catch {}
}

const saveScript = async () => {
  loading.value = true
  try {
    if (editingScript.value) {
      await scriptApi.update(props.subscriptionId, editingScript.value.id, scriptForm)
    } else {
      await scriptApi.create(props.subscriptionId, scriptForm)
    }
    showEditor.value = false
    await fetchScripts()
    ElMessage.success('保存成功')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}

.desc {
  color: #64748b;
}

.hint {
  font-size: 12px;
  color: #64748b;
  font-weight: normal;
  margin-top: 4px;
}

.hint code {
  color: #10b981;
  background: #ecfdf5;
  padding: 2px 6px;
  border-radius: 3px;
}

:deep(.editor textarea) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

@media (max-width: 768px) {
  :deep(.el-table) {
    font-size: 12px;
  }

  :deep(.editor textarea) {
    font-size: 12px;
  }
}
</style>

<style>
.mobile-dialog .el-dialog__body {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  overflow: hidden;
}

.mobile-dialog .el-dialog__body > .toolbar {
  flex-shrink: 0;
}

.mobile-dialog .el-dialog__body > .table-wrapper {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.mobile-dialog .el-dialog__body > .el-form {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.mobile-dialog .el-dialog__footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1px solid #e2e8f0;
  padding: 16px 20px;
  z-index: 100;
}
</style>