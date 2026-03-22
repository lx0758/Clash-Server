<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h1>ClashServer</h1>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        :label-position="isMobile ? 'top' : 'right'"
        @submit.prevent="handleSubmit"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            :prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            native-type="submit"
            :loading="loading"
            size="large"
            class="submit-btn"
          >
            {{ isInit ? '初始化' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useBreakpoint } from '@/composables/useBreakpoint'

const router = useRouter()
const authStore = useAuthStore()
const { isMobile } = useBreakpoint()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isInit = ref(false)

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度为 2-50 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 100, message: '密码长度为 6-100 个字符', trigger: 'blur' },
  ],
}

onMounted(async () => {
  try {
    const initialized = await authStore.checkInit()
    isInit.value = !initialized
  } catch {
    ElMessage.error('无法连接服务器，请检查网络连接')
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    if (isInit.value) {
      await authStore.initPassword(form.username, form.password)
    }
    await authStore.login(form.username, form.password)
    router.push('/')
  } catch (e: unknown) {
    if (e instanceof Error) {
      ElMessage.error(e.message || '操作失败')
    } else {
      ElMessage.error('操作失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  padding: 20px;
}

.login-card {
  width: 400px;
  max-width: 100%;
}

.login-card :deep(.el-card__header) {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid var(--el-border-color);
}

h1 {
  color: var(--el-color-primary);
  margin: 0;
  font-size: 28px;
  font-weight: 600;
}

.login-card :deep(.el-card__body) {
  padding: 30px 20px;
}

.submit-btn {
  width: 100%;
}

@media (max-width: 768px) {
  .login-card {
    width: 100%;
    max-width: none;
  }

  .login-card :deep(.el-card__body) {
    padding: 20px;
  }
}
</style>