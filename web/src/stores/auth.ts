import { defineStore } from 'pinia'
import { ref } from 'vue'
import { config } from '@/config'

export interface User {
  id: number
  username: string
  role: string
  nickname: string
  avatar: string
  status: number
  created_at?: string
  updated_at?: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isAuthenticated = ref(false)

  // 初始化状态
  const initAuth = () => {
    const storedToken = localStorage.getItem(config.storageKeys.token)
    const storedUser = localStorage.getItem(config.storageKeys.user)
    
    if (storedToken && storedUser) {
      try {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
        isAuthenticated.value = true
      } catch (error) {
        console.error('解析用户信息失败:', error)
        clearAuth()
      }
    }
  }

  // 设置用户信息
  const setUser = (userData: User) => {
    user.value = userData
    localStorage.setItem(config.storageKeys.user, JSON.stringify(userData))
    isAuthenticated.value = true
  }

  // 设置token
  const setToken = (tokenData: string) => {
    token.value = tokenData
    localStorage.setItem(config.storageKeys.token, tokenData)
  }

  // 清除认证信息
  const clearAuth = () => {
    user.value = null
    token.value = null
    isAuthenticated.value = false
    localStorage.removeItem(config.storageKeys.token)
    localStorage.removeItem(config.storageKeys.user)
  }

  // 登出
  const logout = () => {
    clearAuth()
  }

  // 检查是否有权限
  const hasRole = (role: string) => {
    return user.value?.role === role
  }

  // 检查是否是管理员
  const isAdmin = () => {
    return hasRole('admin')
  }

  return {
    user,
    token,
    isAuthenticated,
    initAuth,
    setUser,
    setToken,
    clearAuth,
    logout,
    hasRole,
    isAdmin
  }
})
