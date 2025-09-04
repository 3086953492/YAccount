<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <h2 class="title">欢迎回来</h2>
          <p class="subtitle">登录您的{{ systemStore.getConfig('system_name', 'YAccount') }}账户</p>
        </div>
      </template>

      <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="true" @close="errorMessage = ''"
        show-icon class="error-alert" />

      <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" @submit.prevent="handleLogin"
        class="login-form" label-position="top" size="large">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" :prefix-icon="User" :disabled="loading"
            clearable />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" :prefix-icon="Lock"
            :disabled="loading" show-password clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" class="login-button" size="large">
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="form-footer">
        <span class="footer-text">还没有账号？</span>
        <el-link type="primary" @click="$router.push('/register')" :underline="false">
          立即注册
        </el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSystemStore } from '@/stores/system'
import { login } from '@/api/auth'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const systemStore = useSystemStore()

const loading = ref(false)
const errorMessage = ref('')
const loginFormRef = ref<FormInstance>()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await login(loginForm)

    if (response.data.success) {
      // 保存用户信息和OAuth令牌数据
      authStore.setUser(response.data.data.user)
      authStore.setTokenData(response.data.data.token)

      // 显示成功消息
      ElMessage.success('登录成功')

      // 跳转到首页
      router.push('/')
    } else {
      errorMessage.value = response.data.message || '登录失败'
    }
  } catch (error: any) {
    console.error('登录错误:', error)
    errorMessage.value = error.response?.data?.message || '登录失败，请稍后重试'
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
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-bg-color-page);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
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

.error-alert {
  margin-bottom: 20px;
}

.login-form {
  margin-bottom: 20px;
}

.login-button {
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

@media (max-width: 640px) {
  .login-container {
    padding: 15px;
  }

  .login-card {
    max-width: 100%;
  }
}
</style>
