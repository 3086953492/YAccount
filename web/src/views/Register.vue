<template>
  <FormCard title="创建账户" subtitle="加入YAccount，开启您的数字生活">
    <ErrorMessage v-if="errorMessage" :message="errorMessage" type="error" @dismiss="errorMessage = ''" />

    <SuccessMessage v-if="successMessage" :message="successMessage" :dismissible="false" />

    <form @submit.prevent="handleRegister" class="register-form">
      <div class="form-row">
        <FormInput id="username" v-model="registerForm.username" label="用户名" placeholder="3-15位字符" 
          hint="支持字母、数字、下划线" required :disabled="loading" />
        
        <FormInput id="nickname" v-model="registerForm.nickname" label="昵称" placeholder="请输入昵称" 
          hint="将显示在个人资料中" required :disabled="loading" />
      </div>

      <div class="form-row">
        <FormInput id="password" v-model="registerForm.password" label="密码" type="password" placeholder="至少6位密码"
          hint="建议包含字母、数字和特殊字符" required :disabled="loading" />

        <FormInput id="confirmPassword" v-model="registerForm.confirm_password" label="确认密码" type="password"
          placeholder="请再次输入密码" hint="确保两次输入一致" required :disabled="loading" />
      </div>

      <div class="form-actions">
        <FormButton type="submit" text="创建账户" :loading="loading" loading-text="创建中..." />
      </div>

      <FormFooter text="已有账号？" link-text="立即登录" link-to="/login" />
    </form>
  </FormCard>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/auth'
import ErrorMessage from '@/components/ui/ErrorMessage.vue'
import SuccessMessage from '@/components/ui/SuccessMessage.vue'
import FormCard from '@/components/form/FormCard.vue'
import FormInput from '@/components/form/FormInput.vue'
import FormButton from '@/components/form/FormButton.vue'
import FormFooter from '@/components/form/FormFooter.vue'

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const registerForm = reactive({
  username: '',
  nickname: '',
  password: '',
  confirm_password: ''
})

const handleRegister = async () => {
  // 表单验证
  if (!registerForm.username || !registerForm.nickname || !registerForm.password || !registerForm.confirm_password) {
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
.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-actions {
  margin-top: 8px;
}

@media (max-width: 640px) {
  .register-form {
    gap: 20px;
  }
  
  .form-row {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>
