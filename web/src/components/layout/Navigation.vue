<template>
  <nav class="navigation">
    <div class="nav-container">
      <div class="nav-brand">
        <router-link to="/" class="brand-link">
          <span class="brand-icon">🔐</span>
          <span class="brand-text">YAccount</span>
        </router-link>
      </div>

      <div class="nav-menu">
        <router-link to="/" class="nav-link" active-class="active">
          首页
        </router-link>
        <router-link to="/profile" class="nav-link" active-class="active">
          个人资料
        </router-link>
        <router-link v-if="isAdmin" to="/admin" class="nav-link" active-class="active">
          管理后台
        </router-link>
      </div>

      <div class="nav-user">
        <div class="user-info">
          <span class="username">{{ user?.nickname || user?.username }}</span>
          <div class="user-avatar">
            {{ user?.nickname?.charAt(0) || user?.username?.charAt(0) }}
          </div>
        </div>
        <button @click="handleLogout" class="logout-btn">
          登出
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const isAdmin = computed(() => authStore.isAdmin())

const handleLogout = () => {
  if (confirm('确定要登出吗？')) {
    authStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.navigation {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.nav-brand {
  display: flex;
  align-items: center;
}

.brand-link {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: #333;
  font-weight: 600;
  font-size: 20px;
}

.brand-icon {
  font-size: 24px;
}

.brand-text {
  font-size: 20px;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 32px;
}

.nav-link {
  text-decoration: none;
  color: #666;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: #333;
  background-color: #f3f4f6;
}

.nav-link.active {
  color: #667eea;
  background-color: #eff6ff;
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  color: #333;
  font-weight: 500;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.logout-btn {
  background: #ef4444;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.logout-btn:hover {
  background: #dc2626;
}

@media (max-width: 768px) {
  .nav-container {
    padding: 0 15px;
  }

  .nav-menu {
    gap: 16px;
  }

  .nav-link {
    padding: 6px 12px;
    font-size: 14px;
  }

  .username {
    display: none;
  }
}
</style>
