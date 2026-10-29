<template>
  <div class="subscriptions-page">
    <el-alert
      v-if="coreError"
      :type="coreErrorType"
      :title="coreErrorTitle"
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
            :loading="isSubBusy(sub.id)"
            :pending-action="getSubPendingAction(sub.id)"
            @activate="handleActivate(sub.id)"
            @refresh="handleRefresh(sub.id)"
            @edit="openEditDialog(sub)"
            @delete="handleDelete(sub.id)"
            @customize="openCustomizeDialog(sub)"
            @merged-config="openMergedConfig(sub)"
            @content="openContentEditor(sub)"
          />
        </div>
      </template>
    </el-skeleton>

    <SubscriptionEditDialog
      v-if="showSubDialog"
      :subscription="editingSub"
      :loading="isPending('save')"
      @close="closeSubDialog"
      @submit="handleSaveSub"
    />

    <MergedConfigDialog
      v-if="showMergedConfig"
      :subscription-name="mergedConfigName"
      :yaml="mergedConfigYaml"
      @close="closeMergedConfig"
    />

    <CustomizationDialog
      v-if="showCustomizeDialog"
      :subscription-id="currentSubscriptionId"
      :subscription-name="currentSubscriptionName"
      @close="closeCustomizeDialog"
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
import { ref, computed, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { subscriptionApi } from '@/api/subscription'
import { usePendingStore } from '@/stores/pending'
import { useToastStore } from '@/stores/toast'
import { ApiError } from '@/utils/request'
import type { Subscription, PendingAction } from '@/types/api'
import SubscriptionCard from '@/components/SubscriptionCard.vue'
import SubscriptionEditDialog from '@/components/SubscriptionEditDialog.vue'
import MergedConfigDialog from '@/components/MergedConfigDialog.vue'
import CustomizationDialog from '@/components/CustomizationDialog.vue'
import SubscriptionContentEditor from '@/components/SubscriptionContentEditor.vue'

const pendingStore = usePendingStore()
const toastStore = useToastStore()

const initialLoading = ref(true)
/** 核心未生效的原因, 非空时在页面顶部展示横幅 */
const coreError = ref('')
/** 横幅类型: 半成功为警告, 完全失败为错误 */
const coreErrorType = ref<'error' | 'warning'>('error')
const subscriptions = ref<Subscription[]>([])

const showSubDialog = ref(false)
const editingSub = ref<Subscription | null>(null)

const showMergedConfig = ref(false)
const mergedConfigName = ref('')
const mergedConfigYaml = ref('')

const showCustomizeDialog = ref(false)
const currentSubscriptionId = ref(0)
const currentSubscriptionName = ref('')
const showContentEditor = ref(false)
const contentEditorSub = ref<Subscription | null>(null)

onMounted(() => {
  fetchSubs()
})

const fetchSubs = async () => {
  try {
    const res = await subscriptionApi.list()
    subscriptions.value = res.data.data.subscriptions?.map((item: any) => item.subscription) || []
  } finally {
    initialLoading.value = false
  }
}

/**
 * pending key 约定: `<动作>:<订阅ID>`, 全局动作用 `save`。
 */
const pendingKey = (action: PendingAction, id?: number) =>
  id === undefined ? action : `${action}:${id}`

const isPending = (key: string) => pendingStore.isPending(key)

/** 某订阅当前是否有操作在进行 */
const isSubBusy = (id: number) =>
  (['refresh', 'activate', 'delete'] as PendingAction[]).some(a =>
    pendingStore.isPending(pendingKey(a, id))
  )

/** 某订阅当前进行中的操作类型, 用于在卡片上显示具体状态 */
const getSubPendingAction = (id: number): PendingAction | undefined => {
  for (const action of ['refresh', 'activate', 'delete'] as PendingAction[]) {
    if (pendingStore.isPending(pendingKey(action, id))) return action
  }
  return undefined
}

/**
 * 同一订阅的操作互斥守卫。
 *
 * 某个订阅已有操作进行中时, 拒绝对同一订阅发起新操作,
 * 避免并发提交产生重复副作用。操作结束后自动恢复可用。
 */
const guardSubExclusive = (id: number): boolean => {
  if (isSubBusy(id)) {
    toastStore.info('该订阅正在处理中，请稍候')
    return false
  }
  return true
}

/**
 * 判断异常是否为"已保存但核心未生效"的半成功状态。
 */
const isPartialSuccess = (e: unknown): e is ApiError =>
  e instanceof ApiError && e.isPartialSuccess

/**
 * 统一的操作结果处理, 区分三种语义:
 *   - 完全成功: 绿色提示
 *   - 半成功:   数据已保存但核心未生效, 黄色警告
 *   - 完全失败: 红色错误
 */
const reportOutcome = (e: unknown, successMessage: string, failedMessage: string) => {
  if (isPartialSuccess(e)) {
    const reason = e.coreError || '核心未生效'
    coreError.value = reason
    coreErrorType.value = 'warning'
    toastStore.warning(`${successMessage}，但核心未生效：${reason}`)
    return
  }
  // 完全失败时错误提示已由请求拦截器弹出, 这里补充页面级上下文
  coreError.value = e instanceof Error ? e.message : failedMessage
  coreErrorType.value = 'error'
}

const coreErrorTitle = computed(() =>
  coreErrorType.value === 'warning'
    ? `已保存但未生效：${coreError.value}`
    : `Core 错误: ${coreError.value}`
)

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
  const isEdit = !!editingSub.value
  await pendingStore.withPending('save', async () => {
    try {
      if (isEdit) {
        await subscriptionApi.update(editingSub.value!.id, data)
      } else {
        await subscriptionApi.create(data)
      }
      closeSubDialog()
      await fetchSubs()
      toastStore.success(isEdit ? '订阅已更新' : '订阅已创建')
    } catch (e) {
      reportOutcome(e, isEdit ? '订阅已更新' : '订阅已创建', '保存订阅失败')
      // 半成功时数据已落库, 刷新列表以反映已保存状态
      if (isPartialSuccess(e)) {
        await fetchSubs()
      }
    }
  })
}

const handleActivate = async (id: number) => {
  if (!guardSubExclusive(id)) return
  await pendingStore.withPending(pendingKey('activate', id), async () => {
    try {
      await subscriptionApi.activate(id)
      await fetchSubs()
      toastStore.success('订阅已激活')
    } catch (e) {
      reportOutcome(e, '订阅已激活', '激活订阅失败')
      if (isPartialSuccess(e)) {
        await fetchSubs()
      }
    }
  })
}

const handleRefresh = async (id: number) => {
  if (!guardSubExclusive(id)) return
  await pendingStore.withPending(pendingKey('refresh', id), async () => {
    try {
      await subscriptionApi.refresh(id)
      await fetchSubs()
      toastStore.success('订阅已刷新')
    } catch (e) {
      reportOutcome(e, '订阅已刷新', '刷新订阅失败')
      if (isPartialSuccess(e)) {
        await fetchSubs()
      }
    }
  })
}

const handleDelete = async (id: number) => {
  if (!guardSubExclusive(id)) return
  try {
    await ElMessageBox.confirm(
      '确定要删除此订阅吗？关联的规则和脚本也会被删除。',
      '删除确认',
      { type: 'warning' }
    )
  } catch {
    // 用户取消
    return
  }

  await pendingStore.withPending(pendingKey('delete', id), async () => {
    try {
      await subscriptionApi.delete(id)
      await fetchSubs()
      toastStore.success('订阅已删除')
    } catch (e) {
      reportOutcome(e, '订阅已删除', '删除订阅失败')
      if (isPartialSuccess(e)) {
        await fetchSubs()
      }
    }
  })
}

const openMergedConfig = async (sub: Subscription) => {
  try {
    const res = await subscriptionApi.getMerged(sub.id)
    mergedConfigName.value = sub.name
    mergedConfigYaml.value = res.data.data.yaml
    showMergedConfig.value = true
  } catch {
    // 错误提示已由请求拦截器统一弹出
  }
}

const closeMergedConfig = () => {
  showMergedConfig.value = false
}

const openCustomizeDialog = (sub: Subscription) => {
  currentSubscriptionId.value = sub.id
  currentSubscriptionName.value = sub.name
  showCustomizeDialog.value = true
}

const closeCustomizeDialog = () => {
  showCustomizeDialog.value = false
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
  toastStore.success('内容已保存')
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