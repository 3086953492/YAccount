import api from './auth'

// 获取系统信息列表
export const getSystemInfoList = () => {
  return api.get('/system/infos')
}

// 根据配置键获取系统信息
export const getSystemInfoByKey = (configKey: string) => {
  return api.get(`/system/infos/${configKey}`)
}

// 批量更新系统信息
export const batchUpdateSystemInfo = (data: { system_infos: Array<{ id: number; config_value: string }> }) => {
  return api.post('/system/infos', data)
}
