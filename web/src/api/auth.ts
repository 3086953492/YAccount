import axios from 'axios'
import { config } from '@/config'

// 创建axios实例
const api = axios.create({
  baseURL: config.apiBaseURL,
  timeout: config.requestTimeout,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加OAuth访问令牌
api.interceptors.request.use(
  (axiosConfig) => {
    const tokenDataStr = localStorage.getItem(config.storageKeys.token)
    if (tokenDataStr) {
      try {
        const tokenData = JSON.parse(tokenDataStr)
        if (tokenData.access_token) {
          axiosConfig.headers.Authorization = `${tokenData.token_type || 'Bearer'} ${tokenData.access_token}`
        }
      } catch (error) {
        console.error('解析令牌数据失败:', error)
        // 清除无效的令牌数据
        localStorage.removeItem(config.storageKeys.token)
        localStorage.removeItem(config.storageKeys.user)
      }
    }
    return axiosConfig
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误和令牌刷新
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      // 尝试使用刷新令牌获取新的访问令牌
      const tokenDataStr = localStorage.getItem(config.storageKeys.token)
      if (tokenDataStr) {
        try {
          const tokenData = JSON.parse(tokenDataStr)
          if (tokenData.refresh_token) {
            const refreshResponse = await refreshToken(tokenData.refresh_token)
            if (refreshResponse.data.success) {
              const newTokenData = refreshResponse.data.data
              localStorage.setItem(config.storageKeys.token, JSON.stringify(newTokenData))
              
              // 重新设置请求头并重试原请求
              originalRequest.headers.Authorization = `${newTokenData.token_type || 'Bearer'} ${newTokenData.access_token}`
              return api(originalRequest)
            }
          }
        } catch (refreshError) {
          console.error('刷新令牌失败:', refreshError)
        }
      }
      
      // 刷新失败，清除本地存储并跳转到登录页
      localStorage.removeItem(config.storageKeys.token)
      localStorage.removeItem(config.storageKeys.user)
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 登录接口
export const login = (data: { username: string; password: string }) => {
  return api.post('/auth/login', data, {
    headers: {
      'X-Client-ID': config.clientId
    }
  })
}

// 注册接口
export const register = (data: {
  username: string
  password: string
  confirm_password: string
  nickname: string
}) => {
  return api.post('/auth/register', data)
}

// 刷新OAuth令牌接口
export const refreshToken = (refreshToken: string, clientId?: string) => {
  return api.put('/auth/token', {
    refresh_token: refreshToken,
    client_id: clientId || config.clientId
  })
}

export default api
