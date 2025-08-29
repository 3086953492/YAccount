<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSystemStore } from '@/stores/system'
import { usePageMeta } from '@/composables/usePageMeta'
import MainLayout from '@/components/layout/MainLayout.vue'

const route = useRoute()
const authStore = useAuthStore()
const systemStore = useSystemStore()

// 使用页面meta管理
const { updatePageMeta } = usePageMeta()

// 判断当前页面是否需要导航栏
const needsNavigation = computed(() => {
  return route.meta.requiresAuth !== false
})

onMounted(async () => {
  // 初始化认证状态
  authStore.initAuth()
  
  // 初始化系统信息
  systemStore.init()
  
  // 获取系统信息并更新页面meta
  await systemStore.fetchSystemInfo()
  updatePageMeta()
})
</script>

<template>
  <el-config-provider>
    <div id="app">
      <MainLayout v-if="needsNavigation">
        <router-view />
      </MainLayout>
      <router-view v-else />
    </div>
  </el-config-provider>
</template>

<style>
/* 重置样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  line-height: 1.6;
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

#app {
  min-height: 100vh;
}

/* 全局样式 */
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .container {
    padding: 0 15px;
  }
}

/* Element Plus 主题变量覆盖 */
:root {
  --el-color-primary: #409eff;
  --el-color-success: #67c23a;
  --el-color-warning: #e6a23c;
  --el-color-danger: #f56c6c;
  --el-color-info: #909399;
}

/* 暗色主题支持 */
html.dark {
  --el-bg-color: #141414;
  --el-bg-color-overlay: #1d1e1f;
  --el-text-color-primary: #e5eaf3;
  --el-text-color-regular: #cfd3dc;
  --el-text-color-secondary: #a3a6ad;
}
</style>
