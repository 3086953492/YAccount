export interface SystemInfo {
  id: number
  config_key: string
  config_value: string
  config_type: string
  description: string
  status: number
  created_at: string
  updated_at: string
  created_by: number
  updated_by: number
}

export interface SystemInfoList {
  id: number
  config_key: string
  config_value: string
  config_type: string
  description: string
  status: number
}

// 根据后端实际提供的配置项定义
export interface SystemConfig {
  system_name?: string
  system_icon?: string
  system_logo?: string
  system_description?: string
  [key: string]: string | undefined
}
