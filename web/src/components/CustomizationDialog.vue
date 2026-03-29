<template>
  <el-dialog
    :model-value="true"
    :title="'自定义配置 - ' + subscriptionName"
    :width="isMobile ? '100%' : '800px'"
    :fullscreen="isMobile"
    :close-on-click-modal="false"
    :class="{ 'mobile-dialog': isMobile }"
    @close="$emit('close')"
  >
    <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="false" show-icon class="error-alert" />

    <el-tabs v-model="activeTab">
      <el-tab-pane label="代理" name="proxy">
        <div class="section">
          <div class="section-title">插入节点 <span class="hint">在订阅节点列表头部添加</span></div>
          <el-input v-model="form.proxy_insert" type="textarea" autosize placeholder="- name: my-ss
  type: ss
  server: 1.2.3.4
  port: 443
  cipher: aes-256-gcm
  password: your-password" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">追加节点 <span class="hint">在订阅节点列表尾部添加</span></div>
          <el-input v-model="form.proxy_append" type="textarea" autosize placeholder="- name: my-vmess
  type: vmess
  server: 5.6.7.8
  port: 443
  uuid: your-uuid
  alterId: 0
  cipher: auto" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">移除节点 <span class="hint">从订阅节点列表中移除指定节点</span></div>
          <el-input v-model="form.proxy_remove" type="textarea" autosize placeholder="- node-hk-01
- node-hk-02
- node-us-01" class="yaml-editor" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="代理组" name="proxyGroup">
        <div class="section">
          <div class="section-title">插入代理组 <span class="hint">在订阅代理组列表头部添加</span></div>
          <el-input v-model="form.proxy_group_insert" type="textarea" autosize placeholder="- name: 我的策略组
  type: select
  proxies:
    - node-hk-01
    - node-us-01
    - DIRECT" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">追加代理组 <span class="hint">在订阅代理组列表尾部添加</span></div>
          <el-input v-model="form.proxy_group_append" type="textarea" autosize placeholder="- name: 自动选择
  type: url-test
  proxies:
    - node-hk-01
    - node-hk-02
  url: http://www.gstatic.com/generate_204
  interval: 300" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">移除代理组 <span class="hint">从订阅代理组列表中移除指定代理组</span></div>
          <el-input v-model="form.proxy_group_remove" type="textarea" autosize placeholder="- 机场策略组
- 自动选择" class="yaml-editor" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="规则" name="rule">
        <div class="section">
          <div class="section-title">插入规则 <span class="hint">在订阅规则列表头部添加</span></div>
          <el-input v-model="form.rule_insert" type="textarea" autosize placeholder="- DOMAIN-SUFFIX,google.com,Proxy
- DOMAIN-KEYWORD,facebook,Proxy
- GEOIP,CN,DIRECT" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">追加规则 <span class="hint">在订阅规则列表尾部添加</span></div>
          <el-input v-model="form.rule_append" type="textarea" autosize placeholder="- DOMAIN-SUFFIX,ads.com,REJECT
- DOMAIN-KEYWORD,tracker,REJECT
- MATCH,Proxy" class="yaml-editor" />
        </div>
        <div class="section">
          <div class="section-title">移除规则 <span class="hint">从订阅规则列表中移除指定规则（完整匹配）</span></div>
          <el-input v-model="form.rule_remove" type="textarea" autosize placeholder="- DOMAIN-SUFFIX,ads.com,REJECT
- DOMAIN-KEYWORD,tracker,REJECT" class="yaml-editor" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="全局" name="global">
        <div class="section-hint">覆盖订阅中的全局配置字段（如 dns、ipv6 等）</div>
        <el-input v-model="form.global_override" type="textarea" autosize placeholder="dns:
  enable: true
  enhanced-mode: fake-ip
  nameserver:
    - 114.114.114.114
    - 8.8.8.8

ipv6: true
allow-lan: true" class="yaml-editor" />
      </el-tab-pane>
      <el-tab-pane label="脚本" name="script">
        <div class="section-hint">后处理脚本，接收全局变量 <code>config</code>，修改后返回即可生效</div>
        <el-input v-model="form.script" type="textarea" autosize placeholder="function main(config) {
  config.dns = {
    enable: true,
    'enhanced-mode': 'fake-ip'
  };
  return config;
}" class="yaml-editor" />
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { customizationApi } from '@/api/customization'
import { useBreakpoint } from '@/composables/useBreakpoint'
import * as yaml from 'js-yaml'

const { isMobile } = useBreakpoint()

const props = defineProps<{
  subscriptionId: number
  subscriptionName: string
}>()

const emit = defineEmits<{
  close: []
}>()

const loading = ref(false)
const activeTab = ref('proxy')
const errorMessage = ref('')
const form = reactive({
  proxy_insert: '',
  proxy_append: '',
  proxy_remove: '',
  proxy_group_insert: '',
  proxy_group_append: '',
  proxy_group_remove: '',
  rule_insert: '',
  rule_append: '',
  rule_remove: '',
  global_override: '',
  script: '',
})

const fetchCustomization = async () => {
  try {
    const res = await customizationApi.get(props.subscriptionId)
    const c = res.data.data.customization
    if (c) {
      Object.assign(form, {
        proxy_insert: c.proxy_insert || '',
        proxy_append: c.proxy_append || '',
        proxy_remove: c.proxy_remove || '',
        proxy_group_insert: c.proxy_group_insert || '',
        proxy_group_append: c.proxy_group_append || '',
        proxy_group_remove: c.proxy_group_remove || '',
        rule_insert: c.rule_insert || '',
        rule_append: c.rule_append || '',
        rule_remove: c.rule_remove || '',
        global_override: c.global_override || '',
        script: c.script || '',
      })
    }
  } catch {
    // ignore
  }
}

onMounted(fetchCustomization)
watch(() => props.subscriptionId, fetchCustomization)

const validateYAML = (value: string, fieldName: string): string | null => {
  if (!value.trim()) return null
  try {
    yaml.load(value)
    return null
  } catch (e) {
    return `${fieldName} YAML 格式错误: ${(e as Error).message}`
  }
}

const validate = (): string | null => {
  const yamlFields = [
    { key: 'proxy_insert', name: '插入节点' },
    { key: 'proxy_append', name: '追加节点' },
    { key: 'proxy_remove', name: '移除节点' },
    { key: 'proxy_group_insert', name: '插入代理组' },
    { key: 'proxy_group_append', name: '追加代理组' },
    { key: 'proxy_group_remove', name: '移除代理组' },
    { key: 'rule_insert', name: '插入规则' },
    { key: 'rule_append', name: '追加规则' },
    { key: 'rule_remove', name: '移除规则' },
    { key: 'global_override', name: '全局配置' },
  ]

  for (const field of yamlFields) {
    const err = validateYAML(form[field.key as keyof typeof form], field.name)
    if (err) return err
  }

  return null
}

const handleSave = async () => {
  errorMessage.value = ''

  const validationError = validate()
  if (validationError) {
    errorMessage.value = validationError
    return
  }

  loading.value = true
  try {
    const res = await customizationApi.update(props.subscriptionId, { ...form })
    if (res.data.code !== 0) {
      errorMessage.value = res.data.message || '保存失败'
      return
    }
    ElMessage.success('保存成功')
    emit('close')
  } catch (e: any) {
    const msg = e.response?.data?.message || e.message || '保存失败'
    errorMessage.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.error-alert {
  margin-bottom: 16px;
}

.section {
  margin-bottom: 16px;
}

.section-title {
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 8px;
}

.section-title .hint {
  font-weight: normal;
  font-size: 12px;
  color: #64748b;
  margin-left: 8px;
}

.section-hint {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 12px;
}

.section-hint code {
  color: #10b981;
  background: #ecfdf5;
  padding: 2px 6px;
  border-radius: 3px;
}

:deep(.yaml-editor textarea) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  min-height: 120px;
}

:deep(.el-tabs__content) {
  max-height: 60vh;
  overflow-y: auto;
}
</style>

<style>
.mobile-dialog .el-dialog__body {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 54px);
  overflow: hidden;
  padding: 12px;
}

.mobile-dialog .el-dialog__body > .el-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.mobile-dialog .el-dialog__body > .el-tabs > .el-tabs__header {
  flex-shrink: 0;
}

.mobile-dialog .el-dialog__body > .el-tabs > .el-tabs__content {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: 60px;
  max-height: none;
}

.mobile-dialog .el-dialog__footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1px solid #e2e8f0;
  padding: 12px 16px;
  z-index: 100;
}
</style>
