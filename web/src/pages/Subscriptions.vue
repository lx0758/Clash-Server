<template>
  <div class="subscriptions-page">
    <el-alert
      v-if="coreError"
      type="error"
      :title="'Core 错误: ' + coreError"
      :closable="true"
      @close="coreError = ''"
      show-icon
      class="error-banner"
    />

    <div class="header">
      <h2>订阅管理</h2>
      <el-button type="primary" @click="openAddDialog">
        <el-icon><Plus /></el-icon>
        添加订阅
      </el-button>
    </div>

    <el-skeleton :loading="initialLoading" animated>
      <template #default>
        <el-empty v-if="subscriptions.length === 0" description="暂无订阅">
          <el-button type="primary" @click="openAddDialog">添加第一个订阅</el-button>
        </el-empty>

        <div v-else class="subscription-list">
          <SubscriptionCard
            v-for="sub in subscriptions"
            :key="sub.id"
            :subscription="sub"
            :rule-count="ruleCounts[sub.id] || 0"
            :script-count="scriptCounts[sub.id] || 0"
            :loading="loading"
            @activate="handleActivate(sub.id)"
            @refresh="handleRefresh(sub.id)"
            @edit="openEditDialog(sub)"
            @delete="handleDelete(sub.id)"
            @rules="openRuleEditor(sub)"
            @scripts="openScriptDialog(sub)"
            @merged-config="openMergedConfig(sub)"
            @content="openContentEditor(sub)"
          />
        </div>
      </template>
    </el-skeleton>

    <SubscriptionEditDialog
      v-if="showSubDialog"
      :subscription="editingSub"
      :loading="loading"
      @close="closeSubDialog"
      @submit="handleSaveSub"
    />

    <MergedConfigDialog
      v-if="showMergedConfig"
      :subscription-name="mergedConfigName"
      :yaml="mergedConfigYaml"
      @close="closeMergedConfig"
    />

    <RuleEditorDialog
      v-if="showRuleDialog"
      :subscription-id="currentSubscriptionId"
      :subscription-name="currentSubscriptionName"
      @close="closeRuleDialog"
    />

    <ScriptEditorDialog
      v-if="showScriptDialog"
      :subscription-id="currentSubscriptionId"
      :subscription-name="currentSubscriptionName"
      @close="closeScriptDialog"
    />

    <SubscriptionContentEditor
      v-if="showContentEditor && contentEditorSub"
      :subscription="contentEditorSub"
      @close="closeContentEditor"
      @saved="handleContentSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { subscriptionApi } from '@/api/subscription'
import type { Subscription } from '@/types/api'
import SubscriptionCard from '@/components/SubscriptionCard.vue'
import SubscriptionEditDialog from '@/components/SubscriptionEditDialog.vue'
import MergedConfigDialog from '@/components/MergedConfigDialog.vue'
import RuleEditorDialog from '@/components/RuleEditorDialog.vue'
import ScriptEditorDialog from '@/components/ScriptEditorDialog.vue'
import SubscriptionContentEditor from '@/components/SubscriptionContentEditor.vue'

const loading = ref(false)
const initialLoading = ref(true)
const coreError = ref('')
const subscriptions = ref<Subscription[]>([])
const ruleCounts = ref<Record<number, number>>({})
const scriptCounts = ref<Record<number, number>>({})

const showSubDialog = ref(false)
const editingSub = ref<Subscription | null>(null)

const showMergedConfig = ref(false)
const mergedConfigName = ref('')
const mergedConfigYaml = ref('')

const showRuleDialog = ref(false)
const showScriptDialog = ref(false)
const showContentEditor = ref(false)
const currentSubscriptionId = ref(0)
const currentSubscriptionName = ref('')
const contentEditorSub = ref<Subscription | null>(null)

onMounted(() => {
  fetchSubs()
})

const fetchSubs = async () => {
  try {
    const res = await subscriptionApi.list()
    const items = res.data.data.subscriptions || []
    subscriptions.value = items.map((item: any) => ({
      ...item.subscription,
      rule_count: item.rule_count,
      script_count: item.script_count,
    }))
    for (const sub of subscriptions.value) {
      ruleCounts.value[sub.id] = sub.rule_count || 0
      scriptCounts.value[sub.id] = sub.script_count || 0
    }
  } finally {
    initialLoading.value = false
  }
}

const fetchCounts = async (id: number) => {
  try {
    const res = await subscriptionApi.get(id)
    const data = res.data.data
    ruleCounts.value[id] = data.rules?.length || 0
    scriptCounts.value[id] = data.scripts?.length || 0
  } catch {
    ruleCounts.value[id] = 0
    scriptCounts.value[id] = 0
  }
}

const checkCoreError = (data: { core_error?: string }) => {
  if (data?.core_error) coreError.value = data.core_error
}

const openAddDialog = () => {
  editingSub.value = null
  showSubDialog.value = true
}

const openEditDialog = (sub: Subscription) => {
  editingSub.value = sub
  showSubDialog.value = true
}

const closeSubDialog = () => {
  showSubDialog.value = false
  editingSub.value = null
}

const handleSaveSub = async (data: Partial<Subscription>) => {
  loading.value = true
  try {
    let res
    if (editingSub.value) {
      res = await subscriptionApi.update(editingSub.value.id, data)
    } else {
      res = await subscriptionApi.create(data)
    }
    checkCoreError(res.data)
    closeSubDialog()
    await fetchSubs()
    ElMessage.success(editingSub.value ? '订阅已更新' : '订阅已创建')
  } finally {
    loading.value = false
  }
}

const handleActivate = async (id: number) => {
  loading.value = true
  try {
    const res = await subscriptionApi.activate(id)
    checkCoreError(res.data)
    await fetchSubs()
    ElMessage.success('订阅已激活')
  } finally {
    loading.value = false
  }
}

const handleRefresh = async (id: number) => {
  loading.value = true
  try {
    const res = await subscriptionApi.refresh(id)
    checkCoreError(res.data)
    await fetchSubs()
    ElMessage.success('订阅已刷新')
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除此订阅吗？关联的规则和脚本也会被删除。',
      '删除确认',
      { type: 'warning' }
    )
    loading.value = true
    const res = await subscriptionApi.delete(id)
    checkCoreError(res.data)
    await fetchSubs()
    ElMessage.success('订阅已删除')
  } catch {
  } finally {
    loading.value = false
  }
}

const openMergedConfig = async (sub: Subscription) => {
  try {
    const res = await subscriptionApi.getMergedConfig(sub.id)
    mergedConfigName.value = sub.name
    mergedConfigYaml.value = res.data.data.yaml
    showMergedConfig.value = true
  } catch {
    ElMessage.error('获取配置失败')
  }
}

const closeMergedConfig = () => {
  showMergedConfig.value = false
}

const openRuleEditor = (sub: Subscription) => {
  currentSubscriptionId.value = sub.id
  currentSubscriptionName.value = sub.name
  showRuleDialog.value = true
}

const closeRuleDialog = () => {
  showRuleDialog.value = false
  fetchCounts(currentSubscriptionId.value)
}

const openScriptDialog = (sub: Subscription) => {
  currentSubscriptionId.value = sub.id
  currentSubscriptionName.value = sub.name
  showScriptDialog.value = true
}

const closeScriptDialog = () => {
  showScriptDialog.value = false
  fetchCounts(currentSubscriptionId.value)
}

const openContentEditor = (sub: Subscription) => {
  contentEditorSub.value = sub
  showContentEditor.value = true
}

const closeContentEditor = () => {
  showContentEditor.value = false
  contentEditorSub.value = null
}

const handleContentSaved = async () => {
  await fetchSubs()
  ElMessage.success('内容已保存')
}
</script>

<style scoped>
.subscriptions-page {
  padding: 0;
}

.error-banner {
  margin-bottom: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h2 {
  margin: 0;
  color: #3b82f6;
}

.subscription-list {
  display: flex;
  flex-direction: column;
}
</style>