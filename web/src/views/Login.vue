<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>用户登录</h1>
        <p>欢迎使用YAccount系统</p>
      </div>

      <ErrorMessage v-if="errorMessage" :message="errorMessage" type="error" @dismiss="errorMessage = ''" />

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input id="username" v-model="loginForm.username" type="text" placeholder="请输入用户名" required
            :disabled="loading" />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="loginForm.password" type="password" placeholder="请输入密码" required
            :disabled="loading" />
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="loading" class="login-btn">
            <span v-if="loading">登录中...</span>
            <span v-else>登录</span>
          </button>
        </div>

        <div class="form-footer">
          <p>还没有账号？ <router-link to="/register">立即注册</router-link></p>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login } from '@/api/auth'
import ErrorMessage from '@/components/ErrorMessage.vue'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const errorMessage = ref('')

const loginForm = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    errorMessage.value = '请填写完整的登录信息'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await login(loginForm)

    if (response.data.success) {
      // 保存用户信息和token
      authStore.setUser(response.data.data.user)
      authStore.setToken(response.data.data.token)

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
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  color: #333;
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 600;
}

.login-header p {
  color: #666;
  margin: 0;
  font-size: 16px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: #333;
  font-weight: 500;
  font-size: 14px;
}

.form-group input {
  padding: 12px 16px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.form-actions {
  margin-top: 10px;
}

.login-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.form-footer {
  text-align: center;
  margin-top: 20px;
}

.form-footer p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.form-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  text-decoration: underline;
}
</style>
