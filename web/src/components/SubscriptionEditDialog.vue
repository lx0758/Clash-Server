<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑订阅' : '添加订阅'"
    :width="isMobile ? '100%' : '500px'"
    :fullscreen="isMobile"
    @close="$emit('close')"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      :label-position="isMobile ? 'top' : 'right'"
      label-width="100px"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="订阅名称" />
      </el-form-item>

      <el-form-item v-if="!isEdit" label="类型" prop="source_type">
        <el-radio-group v-model="form.source_type">
          <el-radio value="remote">远程订阅</el-radio>
          <el-radio value="local">本地配置</el-radio>
        </el-radio-group>
      </el-form-item>

      <template v-if="form.source_type === 'remote'">
        <el-form-item label="URL" prop="url">
          <el-input v-model="form.url" placeholder="订阅地址" type="url" />
        </el-form-item>
        <el-form-item label="刷新间隔">
          <el-input-number
            v-model="form.interval"
            :min="1"
            placeholder="60"
          />
          <span style="margin-left: 8px; color: #a4b0be;">分钟</span>
        </el-form-item>

        <el-collapse v-model="activeCollapse">
          <el-collapse-item title="请求设置" name="advanced">
            <el-form-item>
              <el-switch v-model="form.use_proxy" active-text="使用代理更新订阅" />
            </el-form-item>
            <el-form-item>
              <el-switch v-model="form.use_user_agent" active-text="使用自定义 User-Agent" />
            </el-form-item>
            <el-form-item v-if="form.use_user_agent" label="User-Agent">
              <el-select v-model="form.user_agent" allow-create filterable placeholder="选择或输入">
                <el-option label="ClashForAndroid/2.5.12" value="ClashForAndroid/2.5.12" />
                <el-option label="ClashMeta" value="ClashMeta" />
                <el-option label="ClashX-Pro" value="ClashX-Pro" />
                <el-option label="clash-verge/v1.0.0" value="clash-verge/v1.0.0" />
                <el-option label="Mozilla/5.0" value="Mozilla/5.0" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-switch
                v-model="form.skip_cert"
                active-text="跳过证书校验"
                inactive-text=""
                style="--el-switch-on-color: #ffa502;"
              />
              <div style="color: #ffa502; font-size: 12px; margin-top: 4px;">
                不推荐启用
              </div>
            </el-form-item>
          </el-collapse-item>
        </el-collapse>
      </template>
    </el-form>

    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ isEdit ? '保存' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Subscription } from '@/types/api'
import { useBreakpoint } from '@/composables/useBreakpoint'

const props = defineProps<{
  subscription?: Subscription | null
  loading?: boolean
}>()

const emit = defineEmits<{
  close: []
  submit: [data: Partial<Subscription>]
}>()

const { isMobile } = useBreakpoint()

const visible = ref(true)
const formRef = ref<FormInstance>()
const isEdit = ref(!!props.subscription)
const activeCollapse = ref<string[]>([])

const form = reactive({
  name: '',
  source_type: 'remote' as 'remote' | 'local',
  url: '',
  interval: 60,
  use_proxy: false,
  use_user_agent: false,
  user_agent: '',
  skip_cert: false,
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入订阅名称', trigger: 'blur' },
  ],
  url: [
    { required: true, message: '请输入订阅地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' },
  ],
}

watch(() => props.subscription, (sub) => {
  if (sub) {
    isEdit.value = true
    form.name = sub.name
    form.source_type = sub.source_type
    form.url = sub.url
    form.interval = sub.interval || 60
    form.use_proxy = sub.use_proxy
    form.use_user_agent = !!sub.user_agent
    form.user_agent = sub.user_agent
    form.skip_cert = sub.skip_cert
    if (sub.use_proxy || sub.user_agent || sub.skip_cert) {
      activeCollapse.value = ['advanced']
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  const data: Partial<Subscription> = {
    name: form.name,
    source_type: form.source_type,
  }
  if (form.source_type === 'remote') {
    data.url = form.url
    data.interval = form.interval || 60
    data.use_proxy = form.use_proxy
    if (form.use_user_agent) {
      data.user_agent = form.user_agent
    }
    data.skip_cert = form.skip_cert
  }
  emit('submit', data)
}
</script>

<style scoped>
:deep(.el-collapse-item__header) {
  color: #a4b0be;
}

:deep(.el-collapse-item__content) {
  padding-top: 16px;
}
</style>