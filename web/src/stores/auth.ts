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

export interface TokenData {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
  scope: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const tokenData = ref<TokenData | null>(null)
  const isAuthenticated = ref(false)

  // 初始化状态
  const initAuth = () => {
    const storedTokenData = localStorage.getItem(config.storageKeys.token)
    const storedUser = localStorage.getItem(config.storageKeys.user)

    if (storedTokenData && storedUser) {
      try {
        tokenData.value = JSON.parse(storedTokenData)
        user.value = JSON.parse(storedUser)
        isAuthenticated.value = true
      } catch (error) {
        console.error('解析认证信息失败:', error)
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

  // 更新用户信息（用于同步更新）
  const updateUser = (userData: Partial<User>) => {
    if (user.value) {
      user.value = { ...user.value, ...userData }
      localStorage.setItem(config.storageKeys.user, JSON.stringify(user.value))
    }
  }

  // 设置token数据
  const setTokenData = (token: TokenData) => {
    tokenData.value = token
    localStorage.setItem(config.storageKeys.token, JSON.stringify(token))
  }

  // 获取访问令牌
  const getAccessToken = () => {
    return tokenData.value?.access_token || null
  }

  // 获取刷新令牌
  const getRefreshToken = () => {
    return tokenData.value?.refresh_token || null
  }

  // 清除认证信息
  const clearAuth = () => {
    user.value = null
    tokenData.value = null
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
    tokenData,
    isAuthenticated,
    initAuth,
    setUser,
    setTokenData,
    getAccessToken,
    getRefreshToken,
    updateUser,
    clearAuth,
    logout,
    hasRole,
    isAdmin
  }
})
