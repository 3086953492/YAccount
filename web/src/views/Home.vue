<template>
  <div class="home-container">
    <Navigation />

    <div class="main-content">
      <div class="user-card">
        <div class="avatar-section">
          <div class="avatar">
            {{ user?.nickname?.charAt(0) || user?.username?.charAt(0) }}
          </div>
        </div>

        <div class="user-details">
          <h2>{{ user?.nickname || user?.username }}</h2>
          <p class="username">@{{ user?.username }}</p>
          <p class="email">{{ user?.email }}</p>
          <div class="role-badge" :class="user?.role">
            {{ user?.role === 'admin' ? '管理员' : '普通用户' }}
          </div>
        </div>
      </div>

      <div class="info-cards">
        <div class="info-card">
          <h3>账户状态</h3>
          <p class="status" :class="{ active: user?.status === 1, inactive: user?.status === 0 }">
            {{ user?.status === 1 ? '正常' : '已禁用' }}
          </p>
        </div>

        <div class="info-card">
          <h3>注册时间</h3>
          <p>{{ formatDate(user?.created_at) }}</p>
        </div>

        <div class="info-card">
          <h3>最后更新</h3>
          <p>{{ formatDate(user?.updated_at) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  try {
    return new Date(dateString).toLocaleString('zh-CN')
  } catch {
    return '未知'
  }
}

onMounted(() => {
  // 检查是否已登录
  if (!authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.user-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 30px;
  margin-bottom: 30px;
}

.avatar-section {
  flex-shrink: 0;
}

.avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 32px;
  font-weight: 600;
}

.user-details {
  flex: 1;
}

.user-details h2 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 28px;
  font-weight: 600;
}

.username {
  color: #666;
  margin: 0 0 8px 0;
  font-size: 16px;
}

.email {
  color: #666;
  margin: 0 0 16px 0;
  font-size: 16px;
}

.role-badge {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  text-transform: uppercase;
}

.role-badge.admin {
  background: #ff6b6b;
  color: white;
}

.role-badge.user {
  background: #4ecdc4;
  color: white;
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.info-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.info-card h3 {
  margin: 0 0 16px 0;
  color: #333;
  font-size: 18px;
  font-weight: 600;
}

.info-card p {
  margin: 0;
  color: #666;
  font-size: 16px;
}

.status.active {
  color: #4caf50;
  font-weight: 600;
}

.status.inactive {
  color: #f44336;
  font-weight: 600;
}

@media (max-width: 768px) {
  .main-content {
    padding: 20px 15px;
  }

  .user-card {
    flex-direction: column;
    text-align: center;
    padding: 30px 20px;
  }

  .info-cards {
    grid-template-columns: 1fr;
  }
}
</style>
