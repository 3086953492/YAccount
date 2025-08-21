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

// 请求拦截器 - 添加token
api.interceptors.request.use(
  (axiosConfig) => {
    const token = localStorage.getItem(config.storageKeys.token)
    if (token) {
      axiosConfig.headers.Authorization = token
    }
    return axiosConfig
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // token过期，清除本地存储并跳转到登录页
      localStorage.removeItem(config.storageKeys.token)
      localStorage.removeItem(config.storageKeys.user)
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 登录接口
export const login = (data: { username: string; password: string }) => {
  return api.post('/auth/login', data)
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

// 刷新token接口
export const refreshToken = () => {
  return api.put('/auth/token')
}

export default api
