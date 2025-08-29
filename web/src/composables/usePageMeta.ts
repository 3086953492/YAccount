import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSystemStore } from '@/stores/system'

export function usePageMeta() {
  const route = useRoute()
  const systemStore = useSystemStore()
  const pageTitle = ref('')
  const pageIcon = ref('')
  
  // 更新页面标题
  const updatePageTitle = (title?: string) => {
    const systemName = systemStore.getConfig('system_name', 'YAccount')
    const finalTitle = title ? `${title} - ${systemName}` : systemName
    pageTitle.value = finalTitle
    document.title = finalTitle
  }
  
  // 更新页面图标
  const updatePageIcon = () => {
    const systemIcon = systemStore.getConfig('system_icon')
    if (systemIcon) {
      const link = document.querySelector("link[rel*='icon']") as HTMLLinkElement
      if (link) {
        link.href = systemIcon
      } else {
        const newLink = document.createElement('link')
        newLink.rel = 'icon'
        newLink.href = systemIcon
        document.head.appendChild(newLink)
      }
      pageIcon.value = systemIcon
    }
  }
  
  // 更新页面meta信息
  const updatePageMeta = () => {
    updatePageTitle()
    updatePageIcon()
    
    // 更新其他meta信息
    const description = systemStore.getConfig('system_description')
    if (description) {
      let metaDesc = document.querySelector('meta[name="description"]')
      if (!metaDesc) {
        metaDesc = document.createElement('meta')
        metaDesc.setAttribute('name', 'description')
        document.head.appendChild(metaDesc)
      }
      metaDesc.setAttribute('content', description)
    }
  }
  
  // 监听路由变化，更新页面标题
  watch(
    () => route.meta.title,
    (newTitle) => {
      if (newTitle) {
        updatePageTitle(newTitle as string)
      } else {
        updatePageTitle()
      }
    },
    { immediate: true }
  )
  
  // 监听系统配置变化
  watch(
    () => systemStore.config,
    () => {
      updatePageMeta()
    },
    { deep: true }
  )
  
  // 初始化
  onMounted(async () => {
    // 尝试获取系统信息
    await systemStore.fetchSystemInfo()
    updatePageMeta()
  })
  
  return {
    pageTitle,
    pageIcon,
    updatePageTitle,
    updatePageIcon,
    updatePageMeta
  }
}
