<template>
  <el-header class="navigation">
    <div class="nav-container">
      <div class="nav-brand">
        <router-link to="/" class="brand-link">
          <el-icon class="brand-icon" size="24">
            <Lock />
          </el-icon>
          <span class="brand-text">YAccount</span>
        </router-link>
      </div>

      <el-menu :default-active="$route.path" mode="horizontal" class="nav-menu" router background-color="transparent"
        text-color="var(--el-text-color-regular)" active-text-color="var(--el-color-primary)">
        <el-menu-item index="/">
          <el-icon>
            <House />
          </el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item :index="`/user/${user?.id}`">
          <el-icon>
            <User />
          </el-icon>
          <span>个人资料</span>
        </el-menu-item>
        <el-menu-item v-if="isAdmin" index="/admin/users">
          <el-icon>
            <UserFilled />
          </el-icon>
          <span>用户列表</span>
        </el-menu-item>
      </el-menu>

      <div class="nav-user">
        <el-dropdown trigger="click" @command="handleUserCommand">
          <div class="user-info">
            <span class="username">{{ user?.nickname || user?.username }}</span>
            <UserAvatar :size="32" :avatar="user?.avatar" :username="user?.username" :nickname="user?.nickname"
              class="user-avatar" />
            <el-icon class="dropdown-icon">
              <ArrowDown />
            </el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon>
                  <User />
                </el-icon>
                个人资料
              </el-dropdown-item>
              <el-dropdown-item command="logout" divided>
                <el-icon>
                  <SwitchButton />
                </el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </el-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import {
  Lock,
  House,
  User,
  Setting,
  ArrowDown,
  SwitchButton,
  UserFilled
} from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const isAdmin = computed(() => authStore.isAdmin())

const handleUserCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push(`/user/${user.value?.id}`)
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        authStore.logout()
        router.push('/login')
      } catch {
        // 用户取消操作
      }
      break
  }
}
</script>

<style scoped>
.navigation {
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 0;
  height: auto;
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
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 20px;
}

.brand-icon {
  color: var(--el-color-primary);
}

.brand-text {
  font-size: 20px;
}

.nav-menu {
  border-bottom: none;
  flex: 1;
  justify-content: center;
}

.nav-menu :deep(.el-menu-item) {
  height: 64px;
  line-height: 64px;
  border-bottom: none;
  font-weight: 500;
}

.nav-menu :deep(.el-menu-item.is-active) {
  background-color: var(--el-color-primary-light-9);
  border-bottom: 2px solid var(--el-color-primary);
}

.nav-menu :deep(.el-menu-item:hover) {
  background-color: var(--el-fill-color-light);
}

.nav-user {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.user-info:hover {
  background-color: var(--el-fill-color-light);
}

.username {
  color: var(--el-text-color-primary);
  font-weight: 500;
}

.user-avatar {
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.dropdown-icon {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  transition: transform 0.2s ease;
}

.user-info:hover .dropdown-icon {
  transform: rotate(180deg);
}

@media (max-width: 768px) {
  .nav-container {
    padding: 0 15px;
  }

  .nav-menu {
    gap: 16px;
  }

  .nav-menu :deep(.el-menu-item) {
    padding: 0 12px;
    font-size: 14px;
  }

  .username {
    display: none;
  }

  .brand-text {
    font-size: 18px;
  }
}

@media (max-width: 640px) {
  .nav-menu :deep(.el-menu-item span) {
    display: none;
  }

  .nav-menu :deep(.el-menu-item) {
    padding: 0 8px;
  }
}
</style>
