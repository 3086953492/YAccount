<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h1>用户注册</h1>
        <p>创建您的YAccount账户</p>
      </div>

      <ErrorMessage v-if="errorMessage" :message="errorMessage" type="error" @dismiss="errorMessage = ''" />

      <SuccessMessage v-if="successMessage" :message="successMessage" :dismissible="false" />

      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input id="username" v-model="registerForm.username" type="text" placeholder="请输入用户名（3-15位）" required
            :disabled="loading" />
        </div>

        <div class="form-group">
          <label for="email">邮箱</label>
          <input id="email" v-model="registerForm.email" type="email" placeholder="请输入邮箱地址" required
            :disabled="loading" />
        </div>

        <div class="form-group">
          <label for="nickname">昵称</label>
          <input id="nickname" v-model="registerForm.nickname" type="text" placeholder="请输入昵称" required
            :disabled="loading" />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="registerForm.password" type="password" placeholder="请输入密码（至少6位）" required
            :disabled="loading" />
        </div>

        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input id="confirmPassword" v-model="registerForm.confirm_password" type="password" placeholder="请再次输入密码"
            required :disabled="loading" />
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="loading" class="register-btn">
            <span v-if="loading">注册中...</span>
            <span v-else>注册</span>
          </button>
        </div>

        <div class="form-footer">
          <p>已有账号？ <router-link to="/login">立即登录</router-link></p>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/auth'
import ErrorMessage from '@/components/ErrorMessage.vue'
import SuccessMessage from '@/components/SuccessMessage.vue'

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const registerForm = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirm_password: ''
})

const handleRegister = async () => {
  // 表单验证
  if (!registerForm.username || !registerForm.email || !registerForm.nickname || !registerForm.password || !registerForm.confirm_password) {
    errorMessage.value = '请填写完整的注册信息'
    return
  }

  if (registerForm.username.length < 3 || registerForm.username.length > 15) {
    errorMessage.value = '用户名长度必须在3-15位之间'
    return
  }

  if (registerForm.password.length < 6) {
    errorMessage.value = '密码长度至少6位'
    return
  }

  if (registerForm.password !== registerForm.confirm_password) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await register(registerForm)

    if (response.data.success) {
      // 显示成功消息并跳转
      successMessage.value = '注册成功！正在跳转到登录页面...'
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
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 450px;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h1 {
  color: #333;
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 600;
}

.register-header p {
  color: #666;
  margin: 0;
  font-size: 16px;
}

.register-form {
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

.register-btn {
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

.register-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.register-btn:disabled {
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
