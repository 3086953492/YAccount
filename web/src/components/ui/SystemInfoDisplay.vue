<template>
  <div class="system-info-display">
    <el-card class="system-info-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon>
            <InfoFilled />
          </el-icon>
          <span>系统信息</span>
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
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSystemStore } from '@/stores/system'
import { InfoFilled } from '@element-plus/icons-vue'

const systemStore = useSystemStore()

onMounted(async () => {
  // 获取系统信息
  await systemStore.fetchSystemInfo()
})
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
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
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
}
</style>
