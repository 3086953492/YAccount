<template>
  <FormCard title="欢迎回来" subtitle="登录您的YAccount账户">
    <ErrorMessage v-if="errorMessage" :message="errorMessage" type="error" @dismiss="errorMessage = ''" />

    <form @submit.prevent="handleLogin" class="login-form">
      <FormInput id="username" v-model="loginForm.username" label="用户名" placeholder="请输入用户名" required
        :disabled="loading" />

      <FormInput id="password" v-model="loginForm.password" label="密码" type="password" placeholder="请输入密码" required
        :disabled="loading" />

      <div class="form-actions">
        <FormButton type="submit" text="登录" :loading="loading" loading-text="登录中..." />
      </div>

      <FormFooter text="还没有账号？" link-text="立即注册" link-to="/register" />
    </form>
  </FormCard>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login } from '@/api/auth'
import ErrorMessage from '@/components/ui/ErrorMessage.vue'
import FormCard from '@/components/form/FormCard.vue'
import FormInput from '@/components/form/FormInput.vue'
import FormButton from '@/components/form/FormButton.vue'
import FormFooter from '@/components/form/FormFooter.vue'

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
.login-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-actions {
  margin-top: 8px;
}

@media (max-width: 640px) {
  .login-form {
    gap: 20px;
  }
}
</style>
