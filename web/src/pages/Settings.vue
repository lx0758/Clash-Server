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
            :type="saveResult.core_error ? 'error' : 'success'"
            :title="saveResult.core_error ? '配置已保存，但 Core 启动失败: ' + saveResult.core_error : '配置已保存，Core 重启成功'"
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
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { useBreakpoint } from '@/composables/useBreakpoint'
import type { CoreConfig } from '@/types/api'

const configStore = useConfigStore()
const { isMobile } = useBreakpoint()

const saving = ref(false)
const saveResult = ref<{ core_error?: string } | null>(null)

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

onMounted(async () => {
  await configStore.fetchConfig()
  if (configStore.config?.core) {
    Object.assign(coreConfig, configStore.config.core)
  }
})

const saveConfig = async () => {
  saving.value = true
  saveResult.value = null
  try {
    const result = await configStore.updateCoreConfig(coreConfig)
    saveResult.value = result
    ElMessage.success('配置已保存')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
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