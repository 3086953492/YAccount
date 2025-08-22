<template>
  <div class="home-container">
    <Navigation />

    <div class="main-content">
      <el-card class="user-card" shadow="hover">
        <div class="user-info">
          <UserAvatar :size="80" :avatar="user?.avatar" :username="user?.username" :nickname="user?.nickname"
            class="user-avatar" />

          <div class="user-details">
            <h2>{{ user?.nickname || user?.username }}</h2>
            <p class="username">@{{ user?.username }}</p>
            <el-tag :type="user?.role === 'admin' ? 'danger' : 'success'" size="large" effect="dark">
              {{ user?.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </div>
        </div>
      </el-card>

      <div class="info-cards">
        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon>
                <User />
              </el-icon>
              <span>账户状态</span>
            </div>
          </template>
          <el-tag :type="user?.status === 1 ? 'success' : 'danger'" size="large" effect="light">
            {{ user?.status === 1 ? '正常' : '已禁用' }}
          </el-tag>
        </el-card>

        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon>
                <Calendar />
              </el-icon>
              <span>注册时间</span>
            </div>
          </template>
          <p class="info-text">{{ formatDate(user?.created_at) }}</p>
        </el-card>

        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon>
                <Clock />
              </el-icon>
              <span>最后更新</span>
            </div>
          </template>
          <p class="info-text">{{ formatDate(user?.updated_at) }}</p>
        </el-card>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/layout/Navigation.vue'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import { User, Calendar, Clock } from '@element-plus/icons-vue'

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

// 查看自己的资料
const viewMyProfile = () => {
  if (user.value?.id) {
    router.push(`/user/${user.value.id}`)
  } else {
    console.error('用户ID不存在，无法查看资料')
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
  background: var(--el-bg-color-page);
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.user-card {
  margin-bottom: 30px;
  border-radius: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 30px;
}

.user-avatar {
  flex-shrink: 0;
  font-size: 32px;
  font-weight: 600;
  color: white;
}

.user-details {
  flex: 1;
}

.user-details h2 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
  font-size: 28px;
  font-weight: 600;
}

.username {
  color: var(--el-text-color-regular);
  margin: 0 0 16px 0;
  font-size: 16px;
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.info-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.info-text {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 16px;
}

@media (max-width: 768px) {
  .main-content {
    padding: 20px 15px;
  }

  .user-info {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }

  .info-cards {
    grid-template-columns: 1fr;
  }
}
</style>
