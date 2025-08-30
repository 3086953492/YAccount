<template>
  <div class="system-management">
    <el-card class="management-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon>
              <Setting />
            </el-icon>
            <span>系统信息管理</span>
          </div>
          <div class="header-right">
            <el-button
              type="primary"
              @click="handleRefresh"
              :loading="loading"
              :icon="Refresh"
            >
              刷新
            </el-button>
            <el-button
              type="success"
              @click="handleEdit"
              :icon="Edit"
              :disabled="!hasEditableItems"
            >
              编辑配置{{ hasEditableItems ? ` (${systemStore.systemInfoList.filter(item => item.status === 1).length})` : '' }}
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 系统信息列表 -->
      <div class="system-info-section">
        <el-table
          :data="systemStore.systemInfoList"
          v-loading="loading"
          stripe
          border
          class="system-info-table"
        >
          <el-table-column prop="config_key" label="配置键" width="180" />
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column prop="config_value" label="配置值" min-width="250">
            <template #default="{ row }">
              <div class="config-value-cell">
                <span v-if="row.config_type === 'url'" class="url-value">
                  <el-link :href="row.config_value" target="_blank" type="primary">
                    {{ row.config_value }}
                  </el-link>
                </span>
                <span v-else-if="row.config_type === 'image'" class="image-value">
                  <el-image
                    :src="row.config_value"
                    :preview-src-list="[row.config_value]"
                    fit="contain"
                    style="width: 60px; height: 30px;"
                    class="config-image"
                  />
                </span>
                <span v-else class="text-value">
                  {{ row.config_value || '-' }}
                </span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="config_type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.config_type)" size="small">
                {{ getTypeDisplayName(row.config_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="row.status === 1"
                type="primary"
                size="small"
                @click="handleEditSingle(row)"
                :icon="Edit"
              >
                编辑
              </el-button>
              <el-button
                v-else
                type="info"
                size="small"
                disabled
              >
                已禁用
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    
    <!-- 系统信息编辑组件 -->
    <SystemInfoEditor
      v-model:visible="editorVisible"
      :system-info-list="editingItems"
      @success="handleEditSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Setting, Edit, Refresh } from '@element-plus/icons-vue'
import { useSystemStore } from '@/stores/system'
import SystemInfoEditor from '@/components/ui/SystemInfoEditor.vue'
import type { SystemInfoList } from '@/types/system'

const systemStore = useSystemStore()
const loading = ref(false)
const editorVisible = ref(false)
const editingItems = ref<SystemInfoList[]>([])

// 是否有可编辑的项目
const hasEditableItems = computed(() => {
  return systemStore.systemInfoList.some(item => item.status === 1)
})

onMounted(async () => {
  await fetchSystemInfo()
})

// 获取系统信息
const fetchSystemInfo = async () => {
  loading.value = true
  try {
    await systemStore.fetchSystemInfo(true)
  } catch (error) {
    ElMessage.error('获取系统信息失败')
  } finally {
    loading.value = false
  }
}

// 刷新数据
const handleRefresh = () => {
  fetchSystemInfo()
}

// 编辑所有配置
const handleEdit = () => {
  editingItems.value = systemStore.systemInfoList.filter(item => item.status === 1)
  editorVisible.value = true
}

// 编辑单个配置
const handleEditSingle = (item: SystemInfoList) => {
  editingItems.value = [item]
  editorVisible.value = true
}

// 编辑成功回调
const handleEditSuccess = () => {
  ElMessage.success('系统信息更新成功')
  // 重新获取数据
  fetchSystemInfo()
}

// 获取类型标签样式
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    'text': 'info',
    'string': 'info',
    'textarea': 'warning',
    'url': 'success',
    'image': 'primary',
    'number': 'danger'
  }
  return typeMap[type] || 'info'
}

// 获取类型显示名称
const getTypeDisplayName = (type: string) => {
  const typeMap: Record<string, string> = {
    'text': '文本',
    'string': '字符串',
    'textarea': '多行文本',
    'url': '链接',
    'image': '图片',
    'number': '数字'
  }
  return typeMap[type] || type
}
</script>

<style scoped>
.system-management {
  padding: 20px;
}

.management-card {
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
  gap: 12px;
}

.system-info-section {
  margin-top: 20px;
}

.system-info-table {
  border-radius: 8px;
  overflow: hidden;
}

.config-value-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.url-value {
  word-break: break-all;
}

.image-value {
  display: flex;
  align-items: center;
}

.config-image {
  border-radius: 4px;
  border: 1px solid var(--el-border-color-light);
}

.text-value {
  word-break: break-all;
  max-width: 200px;
}

@media (max-width: 768px) {
  .system-management {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .header-right {
    align-self: stretch;
    justify-content: flex-end;
  }
  
  .system-info-table {
    font-size: 12px;
  }
}
</style>
