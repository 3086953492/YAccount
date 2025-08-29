import api from './auth'

// 获取系统信息列表
export const getSystemInfoList = () => {
  return api.get('/system/infos')
}

// 根据配置键获取系统信息
export const getSystemInfoByKey = (configKey: string) => {
  return api.get(`/system/infos/${configKey}`)
}
