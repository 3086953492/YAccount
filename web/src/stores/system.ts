import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getSystemInfoList, batchUpdateSystemInfo } from '@/api/system'
import type { SystemInfo, SystemConfig, BatchUpdateSystemInfoRequest } from '@/types/system'

export const useSystemStore = defineStore('system', () => {
  // 系统信息列表
  const systemInfoList = ref<SystemInfo[]>([])
  
  // 系统配置缓存
  const systemConfig = ref<SystemConfig>({})
  
  // 缓存时间戳
  const lastFetchTime = ref<number>(0)
  
  // 缓存过期时间（5分钟）
  const CACHE_EXPIRY_TIME = 5 * 60 * 1000
  
  // 计算属性：获取系统配置
  const config = computed(() => systemConfig.value)
  
  // 检查缓存是否过期
  const isCacheExpired = () => {
    return Date.now() - lastFetchTime.value > CACHE_EXPIRY_TIME
  }
  
  // 从本地存储加载缓存
  const loadFromCache = () => {
    try {
      const cached = localStorage.getItem('system_config_cache')
      if (cached) {
        const { config: cachedConfig, timestamp } = JSON.parse(cached)
        if (Date.now() - timestamp < CACHE_EXPIRY_TIME) {
          systemConfig.value = cachedConfig
          lastFetchTime.value = timestamp
          return true
        }
      }
    } catch (error) {
      console.warn('加载系统配置缓存失败:', error)
    }
    return false
  }
  
  // 保存到本地存储缓存
  const saveToCache = () => {
    try {
      const cacheData = {
        config: systemConfig.value,
        timestamp: Date.now()
      }
      localStorage.setItem('system_config_cache', JSON.stringify(cacheData))
    } catch (error) {
      console.warn('保存系统配置缓存失败:', error)
    }
  }
  
  // 获取系统信息列表
  const fetchSystemInfo = async (force = false) => {
    console.log(`获取系统信息，强制刷新: ${force}`)
    
    // 如果是强制刷新，完全跳过缓存逻辑
    if (force) {
      console.log('强制刷新模式，跳过所有缓存')
      try {
        const response = await getSystemInfoList()
        if (response.data.success) {
          console.log('强制刷新获取到新数据:', response.data.data)
          systemInfoList.value = response.data.data
          
          // 转换为配置对象
          const config: SystemConfig = {}
          systemInfoList.value.forEach(item => {
            if (item.status === 1) { // 只使用启用的配置
              config[item.config_key] = item.config_value
            }
          })
          
          systemConfig.value = config
          lastFetchTime.value = Date.now()
          
          // 保存到缓存
          saveToCache()
          console.log('强制刷新完成，新配置:', config)
          
          return config
        }
      } catch (error) {
        console.error('强制刷新获取系统信息失败:', error)
        return {}
      }
    }
    
    // 非强制刷新时的正常缓存逻辑
    if (!isCacheExpired() && Object.keys(systemConfig.value).length > 0) {
      console.log('使用现有缓存数据')
      return systemConfig.value
    }
    
    // 尝试从缓存加载
    if (loadFromCache()) {
      console.log('从本地存储加载缓存数据')
      return systemConfig.value
    }
    
    try {
      console.log('从API获取系统信息')
      const response = await getSystemInfoList()
      if (response.data.success) {
        systemInfoList.value = response.data.data
        
        // 转换为配置对象
        const config: SystemConfig = {}
        systemInfoList.value.forEach(item => {
          if (item.status === 1) { // 只使用启用的配置
            config[item.config_key] = item.config_value
          }
        })
        
        systemConfig.value = config
        lastFetchTime.value = Date.now()
        
        // 保存到缓存
        saveToCache()
        
        return config
      }
    } catch (error) {
      console.error('获取系统信息失败:', error)
      // 如果获取失败但有缓存，使用缓存
      if (Object.keys(systemConfig.value).length > 0) {
        return systemConfig.value
      }
    }
    
    return {}
  }
  
  // 获取特定配置值
  const getConfig = (key: string, defaultValue = '') => {
    return systemConfig.value[key] || defaultValue
  }
  
  // 批量更新系统信息
  const updateSystemInfo = async (updateData: BatchUpdateSystemInfoRequest) => {
    try {
      console.log('开始更新系统信息，更新字段数量:', updateData.system_infos.length)
      console.log('更新的字段详情:', updateData.system_infos)
      
      const response = await batchUpdateSystemInfo(updateData)
      if (response.data.success) {
        console.log('系统信息更新成功，已更新', updateData.system_infos.length, '个字段')
        
        // 立即清除所有缓存
        clearCache()
        
        // 强制重新获取系统信息，不使用任何缓存
        const freshData = await fetchSystemInfo(true)
        
        // 确保数据已更新
        if (Object.keys(freshData).length > 0) {
          console.log('缓存刷新成功，新数据:', freshData)
        } else {
          console.warn('缓存刷新后数据为空，可能存在异常')
        }
        
        return true
      }
      return false
    } catch (error) {
      console.error('更新系统信息失败:', error)
      return false
    }
  }
  
  // 清除缓存
  const clearCache = () => {
    console.log('清除系统信息缓存')
    systemConfig.value = {}
    systemInfoList.value = []
    lastFetchTime.value = 0
    localStorage.removeItem('system_config_cache')
  }
  
  // 强制刷新系统信息（清除缓存并重新获取）
  const forceRefresh = async () => {
    console.log('强制刷新系统信息')
    clearCache()
    return await fetchSystemInfo(true)
  }
  
  // 调试方法：显示当前缓存状态
  const debugCacheStatus = () => {
    console.log('=== 系统信息缓存状态 ===')
    console.log('systemConfig:', systemConfig.value)
    console.log('systemInfoList:', systemInfoList.value)
    console.log('lastFetchTime:', lastFetchTime.value)
    console.log('缓存是否过期:', isCacheExpired())
    console.log('本地存储缓存:', localStorage.getItem('system_config_cache'))
    console.log('========================')
  }
  
  // 初始化时尝试加载缓存
  const init = () => {
    loadFromCache()
  }
  
  return {
    systemInfoList,
    systemConfig,
    config,
    fetchSystemInfo,
    getConfig,
    updateSystemInfo,
    clearCache,
    init,
    isCacheExpired,
    forceRefresh,
    debugCacheStatus
  }
})
