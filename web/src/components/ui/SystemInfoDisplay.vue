<template>
  <div class="system-info-display">
    <el-card class="system-info-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon>
              <InfoFilled />
            </el-icon>
            <span>系统信息</span>
          </div>
          <div class="header-right">
            <el-button
              v-if="canEdit"
              type="primary"
              size="small"
              @click="handleEdit"
              :icon="Edit"
            >
              编辑
            </el-button>
          </div>
        </div>
      </template>
      
      <div class="system-info-grid">
        <div v-if="systemStore.getConfig('system_name')" class="info-item">
          <span class="info-label">系统名称：</span>
          <span class="info-value">{{ systemStore.getConfig('system_name') }}</span>
        </div>
        
        <div v-if="systemStore.getConfig('system_description')" class="info-item full-width">
          <span class="info-label">系统描述：</span>
          <span class="info-value">{{ systemStore.getConfig('system_description') }}</span>
        </div>
        
        <div v-if="systemStore.getConfig('system_logo')" class="info-item">
          <span class="info-label">系统Logo：</span>
          <span class="info-value">
            <el-image 
              :src="systemStore.getConfig('system_logo')" 
              :preview-src-list="[systemStore.getConfig('system_logo')]"
              fit="contain"
              style="width: 100px; height: 40px;"
              class="system-logo"
            />
          </span>
        </div>
        
        <div v-if="systemStore.getConfig('system_icon')" class="info-item">
          <span class="info-label">系统图标：</span>
          <span class="info-value">
            <el-image 
              :src="systemStore.getConfig('system_icon')" 
              :preview-src-list="[systemStore.getConfig('system_icon')]"
              fit="contain"
              style="width: 32px; height: 32px;"
              class="system-icon"
            />
          </span>
        </div>
      </div>
    </el-card>
    
    <!-- 系统信息编辑组件 -->
    <SystemInfoEditor
      v-model:visible="editorVisible"
      :system-info-list="systemStore.systemInfoList"
      @success="handleEditSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useAuthStore } from '@/stores/auth'
import { InfoFilled, Edit } from '@element-plus/icons-vue'
import SystemInfoEditor from './SystemInfoEditor.vue'

const systemStore = useSystemStore()
const authStore = useAuthStore()
const editorVisible = ref(false)

// 判断是否可以编辑（管理员权限）
const canEdit = computed(() => {
  return authStore.isAdmin
})

onMounted(async () => {
  // 获取系统信息
  await systemStore.fetchSystemInfo()
})

// 处理编辑按钮点击
const handleEdit = () => {
  editorVisible.value = true
}

// 处理编辑成功
const handleEditSuccess = () => {
  // 编辑成功后，系统信息会自动刷新
  console.log('系统信息更新成功')
}
</script>

<style scoped>
.system-info-display {
  margin: 20px 0;
}

.system-info-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.system-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-weight: 600;
  color: var(--el-text-color-regular);
  min-width: 80px;
}

.info-value {
  color: var(--el-text-color-primary);
  word-break: break-word;
}

.system-logo {
  border-radius: 4px;
  border: 1px solid var(--el-border-color-light);
}

.system-icon {
  border-radius: 4px;
  border: 1px solid var(--el-border-color-light);
}

@media (max-width: 768px) {
  .system-info-grid {
    grid-template-columns: 1fr;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .info-label {
    min-width: auto;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .header-right {
    align-self: stretch;
  }
}
</style>
