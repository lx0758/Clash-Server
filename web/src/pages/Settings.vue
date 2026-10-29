<template>
  <div class="settings">
    <h2>设置</h2>

    <el-skeleton :loading="configStore.loading" animated>
      <template #default>
        <el-alert
          v-if="!configStore.config"
          type="error"
          title="无法加载配置"
          :closable="false"
        />

        <div v-else class="settings-content">
          <el-form
            ref="formRef"
            :model="coreConfig"
            :label-position="isMobile ? 'top' : 'left'"
            label-width="100px"
            class="settings-form"
          >
            <el-card class="form-section">
              <template #header>
                <div class="section-header">
                  <span>Server 配置</span>
                  <el-tag type="info" size="small">只读</el-tag>
                </div>
              </template>
              <el-form-item label="监听地址">
                <el-input :model-value="configStore.config.server.host" disabled />
              </el-form-item>
              <el-form-item label="监听端口">
                <el-input :model-value="configStore.config.server.port" disabled />
              </el-form-item>
              <el-form-item label="数据库路径">
                <el-input :model-value="configStore.config.server.database" disabled />
              </el-form-item>
            </el-card>

            <el-card class="form-section">
              <template #header>
                <span>Core 配置</span>
              </template>
              <el-form-item label="API 主机">
                <el-input v-model="coreConfig.api_host" />
              </el-form-item>
              <el-form-item label="API 端口">
                <el-input-number v-model="coreConfig.api_port" :min="1" :max="65535" />
              </el-form-item>
              <el-form-item label="API 密钥">
                <el-input v-model="coreConfig.api_secret" type="password" show-password />
              </el-form-item>
              <el-form-item label="混合端口">
                <el-input-number v-model="coreConfig.mixed_port" :min="1" :max="65535" />
              </el-form-item>
              <el-form-item label="允许局域网">
                <el-switch v-model="coreConfig.allow_lan" />
              </el-form-item>
              <el-form-item label="日志级别">
                <el-select v-model="coreConfig.log_level">
                  <el-option label="静默" value="silent" />
                  <el-option label="错误" value="error" />
                  <el-option label="警告" value="warning" />
                  <el-option label="信息" value="info" />
                  <el-option label="调试" value="debug" />
                </el-select>
              </el-form-item>
              <el-form-item label="启用 IPv6">
                <el-switch v-model="coreConfig.ipv6" />
              </el-form-item>
            </el-card>
          </el-form>

          <el-alert
            v-if="saveResult"
            :type="saveResult.partial ? 'warning' : saveResult.core_error ? 'error' : 'success'"
            :title="saveResultTitle"
            :closable="false"
            show-icon
            class="save-alert"
          />

          <div class="form-footer">
            <el-button type="primary" :loading="saving" @click="saveConfig">
              {{ saving ? '保存中...' : '保存配置' }}
            </el-button>
          </div>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { useToastStore } from '@/stores/toast'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { ApiError } from '@/utils/request'
import type { CoreConfig } from '@/types/api'

const configStore = useConfigStore()
const toastStore = useToastStore()
const { isMobile } = useBreakpoint()

const saving = ref(false)
const saveResult = ref<{ core_error?: string; partial?: boolean } | null>(null)

/**
 * 需要重启核心才能生效的配置项。
 *
 * 变更其中任意一项都会导致核心重启, 当前连接会被断开,
 * 保存前需明确告知用户。
 */
const RESTART_REQUIRED_KEYS = ['mixed_port', 'api_secret', 'allow_lan', 'ipv6'] as const

/** 保存前的配置快照, 用于判断哪些键发生了变化 */
let originalConfig: CoreConfig

const coreConfig = reactive<CoreConfig>({
  api_host: '',
  api_port: 9090,
  api_secret: '',
  mixed_port: 7890,
  allow_lan: true,
  mode: 'rule',
  log_level: 'info',
  ipv6: false,
})

originalConfig = { ...coreConfig }

onMounted(async () => {
  await configStore.fetchConfig()
  if (configStore.config?.core) {
    Object.assign(coreConfig, configStore.config.core)
  }
  originalConfig = { ...coreConfig }
})

/**
 * 判断本次变更是否涉及需要重启核心的配置项。
 */
const needsRestartConfirm = (): boolean => {
  return RESTART_REQUIRED_KEYS.some(
    key => originalConfig[key] !== coreConfig[key]
  )
}

/** 保存结果的展示文案, 区分成功 / 半成功 / 完全失败三种语义 */
const saveResultTitle = computed(() => {
  if (!saveResult.value) return ''
  if (saveResult.value.partial) {
    return `配置已保存，但核心未生效：${saveResult.value.core_error || '核心应用失败'}。请检查核心状态或手动重启核心。`
  }
  if (saveResult.value.core_error) {
    return `配置保存失败：${saveResult.value.core_error}`
  }
  return '配置已保存，核心已生效'
})

const saveConfig = async () => {
  // 端口类变更会重启核心并断开当前连接, 需先确认
  if (needsRestartConfirm()) {
    try {
      await ElMessageBox.confirm(
        '此操作将重启核心并断开当前连接, 确定继续吗?',
        '确认重启核心',
        { type: 'warning', confirmButtonText: '继续', cancelButtonText: '取消' }
      )
    } catch {
      // 用户取消
      return
    }
  }

  saving.value = true
  saveResult.value = null
  try {
    const result = await configStore.updateCoreConfig(coreConfig)
    saveResult.value = result
    originalConfig = { ...coreConfig }
    toastStore.success('配置已保存并生效')
  } catch (e) {
    // 三态结果: 完全成功 / 已保存但核心未生效 / 完全失败
    if (e instanceof ApiError && e.isPartialSuccess) {
      saveResult.value = { core_error: e.coreError, partial: true }
      originalConfig = { ...coreConfig }
      toastStore.warning(
        `配置已保存但核心未生效：${e.coreError || '核心应用失败'}。请检查核心状态或手动重启核心。`
      )
    } else {
      // 完全失败时错误提示已由请求拦截器弹出
      saveResult.value = { core_error: e instanceof Error ? e.message : '保存失败' }
    }
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.settings {
  padding: 0;
}

h2 {
  margin-bottom: 20px;
  color: #3b82f6;
}

.settings-content {
  width: 100%;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-section {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.save-alert {
  margin-top: 16px;
}

.form-footer {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

:deep(.el-form-item__label) {
  color: #64748b;
}

@media (max-width: 768px) {
  h2 {
    font-size: 20px;
  }
}
</style>