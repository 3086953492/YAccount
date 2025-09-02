import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/Home.vue'),
      meta: { requiresAuth: true, title: '首页' }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false, title: '登录' }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresAuth: false, title: '注册' }
    },
    {
      path: '/user/:id',
      name: 'user-profile',
      component: () => import('@/views/Profile.vue'),
      meta: { requiresAuth: true, title: '用户资料' }
    },
    {
      path: '/admin/users',
      name: 'user-list',
      component: () => import('@/views/UserList.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, title: '用户管理' }
    },
    {
      path: '/admin/system',
      name: 'system-management',
      component: () => import('@/views/SystemManagement.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, title: '系统管理' }
    },
    {
      path: '/admin/oauth/clients',
      name: 'oauth-clients',
      component: () => import('@/views/OAuthClients.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, title: 'OAuth客户端管理' }
    },
    {
      path: '/admin/oauth/clients/new',
      name: 'oauth-client-create',
      component: () => import('@/views/OAuthClientForm.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, title: '新建OAuth客户端' }
    },
    {
      path: '/admin/oauth/clients/:clientId/edit',
      name: 'oauth-client-edit',
      component: () => import('@/views/OAuthClientForm.vue'),
      meta: { requiresAuth: true, requiresAdmin: true, title: '编辑OAuth客户端' }
    }
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // 初始化认证状态
  authStore.initAuth()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // 需要认证但未登录，跳转到登录页
    next('/login')
  } else if (to.meta.requiresAdmin && !authStore.isAdmin()) {
    // 需要管理员权限但用户不是管理员，跳转到首页
    next('/')
  } else if (to.name === 'login' && authStore.isAuthenticated) {
    // 已登录用户访问登录页，跳转到首页
    next('/')
  } else if (to.name === 'register' && authStore.isAuthenticated) {
    // 已登录用户访问注册页，跳转到首页
    next('/')
  } else {
    next()
  }
})

export default router
