<template>
  <div class="register-container">
    <el-card class="register-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <h2 class="title">创建账户</h2>
          <p class="subtitle">加入{{ systemStore.getConfig('system_name', 'YAccount') }}，开启您的数字生活</p>
        </div>
      </template>

      <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="true" @close="errorMessage = ''"
        show-icon class="error-alert" />

      <el-alert v-if="successMessage" :title="successMessage" type="success" :closable="false" show-icon
        class="success-alert" />

      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" @submit.prevent="handleRegister"
        class="register-form" label-position="top" size="large">
        <div class="form-row">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="registerForm.username" placeholder="3-15位字符" :prefix-icon="User" :disabled="loading"
              clearable />
            <template #label>
              <span>用户名</span>
              <el-tooltip content="支持字母、数字、下划线" placement="top">
                <el-icon class="hint-icon">
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </template>
          </el-form-item>

          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="registerForm.nickname" placeholder="请输入昵称" :prefix-icon="UserFilled" :disabled="loading"
              clearable />
            <template #label>
              <span>昵称</span>
              <el-tooltip content="将显示在个人资料中" placement="top">
                <el-icon class="hint-icon">
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </template>
          </el-form-item>
        </div>

        <div class="form-row">
          <el-form-item label="密码" prop="password">
            <el-input v-model="registerForm.password" type="password" placeholder="至少6位密码" :prefix-icon="Lock"
              :disabled="loading" show-password clearable />
            <template #label>
              <span>密码</span>
              <el-tooltip content="建议包含字母、数字和特殊字符" placement="top">
                <el-icon class="hint-icon">
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </template>
          </el-form-item>

          <el-form-item label="确认密码" prop="confirm_password">
            <el-input v-model="registerForm.confirm_password" type="password" placeholder="请再次输入密码" :prefix-icon="Lock"
              :disabled="loading" show-password clearable />
            <template #label>
              <span>确认密码</span>
              <el-tooltip content="确保两次输入一致" placement="top">
                <el-icon class="hint-icon">
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </template>
          </el-form-item>
        </div>

        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" class="register-button" size="large">
            {{ loading ? '创建中...' : '创建账户' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="form-footer">
        <span class="footer-text">已有账号？</span>
        <el-link type="primary" @click="$router.push('/login')" :underline="false">
          立即登录
        </el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { register } from '@/api/auth'
import { User, UserFilled, Lock, QuestionFilled } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const systemStore = useSystemStore()
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const registerFormRef = ref<FormInstance>()

const registerForm = reactive({
  username: '',
  nickname: '',
  password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 15, message: '用户名长度必须在3-15位之间', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度必须在2-20位之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return

  try {
    await registerFormRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await register(registerForm)

    if (response.data.success) {
      // 显示成功消息并跳转
      successMessage.value = '注册成功！正在跳转到登录页面...'
      ElMessage.success('注册成功！')
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    } else {
      errorMessage.value = response.data.message || '注册失败'
    }
  } catch (error: any) {
    console.error('注册错误:', error)
    errorMessage.value = error.response?.data?.message || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // 获取系统信息
  await systemStore.fetchSystemInfo()
})
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-bg-color-page);
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 600px;
  border-radius: 16px;
}

.card-header {
  text-align: center;
}

.title {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.error-alert,
.success-alert {
  margin-bottom: 20px;
}

.register-form {
  margin-bottom: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.register-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.form-footer {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}

.footer-text {
  color: var(--el-text-color-regular);
  margin-right: 8px;
}

.hint-icon {
  margin-left: 4px;
  color: var(--el-text-color-secondary);
  cursor: help;
}

@media (max-width: 640px) {
  .register-container {
    padding: 15px;
  }

  .register-card {
    max-width: 100%;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>
