// 环境配置
export const config = {
  // API基础URL - 使用相对路径通过代理
  apiBaseURL: import.meta.env.VITE_API_BASE_URL || '/api/account/v1',
  
  // 请求超时时间（毫秒）
  requestTimeout: 10000,
  
  // 本地存储键名
  storageKeys: {
    token: 'token',
    user: 'user'
  }
}

export default config
