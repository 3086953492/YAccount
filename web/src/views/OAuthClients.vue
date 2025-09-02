<template>
  <div class="oauth-clients-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon size="20">
              <Key />
            </el-icon>
            <span>OAuth客户端管理</span>
          </div>
          <el-dropdown @command="handleCreateCommand">
            <el-button type="primary">
              <el-icon>
                <Plus />
              </el-icon>
              新建客户端
              <el-icon class="el-icon--right">
                <ArrowDown />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="dialog">对话框模式</el-dropdown-item>
                <el-dropdown-item command="page">页面模式</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>

      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form :model="searchForm" inline>
          <el-form-item label="客户端名称">
            <el-input v-model="searchForm.name" placeholder="请输入客户端名称" clearable @keyup.enter="handleSearch"
              style="width: 240px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch" :loading="loading">
              <el-icon>
                <Search />
              </el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon>
                <Refresh />
              </el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格区域 -->
      <div class="table-section">
        <el-table :data="tableData" v-loading="loading" stripe border style="width: 100%"
          @sort-change="handleSortChange">
          <el-table-column prop="name" label="客户端名称" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="client-name">
                <span class="name">{{ row.name }}</span>
                <el-tag v-if="row.description" size="small" type="info" class="description-tag">
                  {{ row.description }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="client_id" label="Client ID" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="client-id">
                <code>{{ row.client_id }}</code>
                <el-button type="primary" link size="small" @click="copyToClipboard(row.client_id)" class="copy-btn">
                  <el-icon>
                    <CopyDocument />
                  </el-icon>
                </el-button>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="client_type" label="客户端类型" width="120" align="center">
            <template #default="{ row }">
              <el-tag :type="row.client_type === 'confidential' ? 'success' : 'warning'">
                {{ row.client_type === 'confidential' ? '机密型' : '公开型' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
                {{ row.status === 'active' ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="grant_types" label="授权类型" min-width="140" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="grant-types">
                <el-tag v-for="type in parseGrantTypes(row.grant_types)" :key="type" size="small"
                  class="grant-type-tag">
                  {{ formatGrantType(type) }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="scopes" label="权限范围" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="scopes">
                <el-tag v-for="scope in parseScopes(row.scopes)" :key="scope" size="small" type="info"
                  class="scope-tag">
                  {{ scope }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>

          <el-table-column label="操作" width="200" fixed="right" align="center">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button type="primary" link size="small" @click="handleView(row)">
                  <el-icon>
                    <View />
                  </el-icon>
                  查看
                </el-button>
                <el-dropdown @command="(command) => handleEditCommand(command, row)" trigger="click">
                  <el-button type="primary" link size="small">
                    <el-icon>
                      <Edit />
                    </el-icon>
                    编辑
                    <el-icon class="el-icon--right" style="font-size: 10px;">
                      <ArrowDown />
                    </el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="dialog">对话框模式</el-dropdown-item>
                      <el-dropdown-item command="page">页面模式</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <el-popconfirm title="确定要删除这个OAuth客户端吗？" @confirm="handleDelete(row)" confirm-button-text="确定"
                  cancel-button-text="取消">
                  <template #reference>
                    <el-button type="danger" link size="small">
                      <el-icon>
                        <Delete />
                      </el-icon>
                      删除
                    </el-button>
                  </template>
                </el-popconfirm>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-section">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
          :total="pagination.total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>

    <!-- OAuth客户端表单对话框 -->
    <OAuthClientFormDialog
      v-model:visible="dialogVisible"
      :client-data="currentClient"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listOAuthClients, deleteOAuthClient, getOAuthClient, type OAuthClient } from '@/api/oauth'
import { usePageMeta } from '@/composables/usePageMeta'
import OAuthClientFormDialog from '@/components/ui/OAuthClientFormDialog.vue'
import {
  Key,
  Plus,
  Search,
  Refresh,
  CopyDocument,
  View,
  Edit,
  Delete,
  ArrowDown
} from '@element-plus/icons-vue'

const router = useRouter()

// 使用页面元信息
usePageMeta()

// 响应式数据
const loading = ref(false)
const tableData = ref<OAuthClient[]>([])

// 对话框状态
const dialogVisible = ref(false)
const currentClient = ref<OAuthClient | null>(null)

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 排序信息
const sortInfo = reactive({
  prop: '',
  order: ''
})

// 获取客户端列表
const fetchClientList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: searchForm.name || undefined
    }

    const response = await listOAuthClients(params)

    if (response.data.success) {
      console.log(response.data.data)
      tableData.value = response.data.data.items || []
      pagination.total = response.data.data.total || 0
      console.log(tableData.value)
    } else {
      ElMessage.error(response.data.message || '获取客户端列表失败')
    }
  } catch (error: any) {
    console.error('获取客户端列表失败:', error)
    ElMessage.error(error.response?.data?.message || '获取客户端列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchClientList()
}

// 重置搜索
const handleReset = () => {
  searchForm.name = ''
  pagination.page = 1
  fetchClientList()
}

// 分页大小改变
const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchClientList()
}

// 页码改变
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchClientList()
}

// 排序改变
const handleSortChange = ({ prop, order }: { prop: string; order: string }) => {
  sortInfo.prop = prop
  sortInfo.order = order
  fetchClientList()
}

// 处理新建客户端命令
const handleCreateCommand = (command: string) => {
  if (command === 'dialog') {
    currentClient.value = null
    dialogVisible.value = true
  } else if (command === 'page') {
    router.push('/admin/oauth/clients/new')
  }
}

// 新建客户端（保留向后兼容）
const handleCreate = () => {
  router.push('/admin/oauth/clients/new')
}

// 查看客户端详情
const handleView = (row: OAuthClient) => {
  // TODO: 实现查看客户端详情功能
  ElMessage.info('查看客户端详情功能暂未实现')
}

// 处理编辑客户端命令
const handleEditCommand = async (command: string, row: OAuthClient) => {
  if (command === 'dialog') {
    try {
      const response = await getOAuthClient(row.client_id)
      if (response.data.success) {
        currentClient.value = response.data.data
        dialogVisible.value = true
      } else {
        ElMessage.error(response.data.message || '获取客户端详情失败')
      }
    } catch (error: any) {
      console.error('获取客户端详情失败:', error)
      ElMessage.error(error.response?.data?.message || '获取客户端详情失败')
    }
  } else if (command === 'page') {
    router.push(`/admin/oauth/clients/${row.client_id}/edit`)
  }
}

// 编辑客户端（保留向后兼容）
const handleEdit = (row: OAuthClient) => {
  router.push(`/admin/oauth/clients/${row.client_id}/edit`)
}

// 删除客户端
const handleDelete = async (row: OAuthClient) => {
  try {
    const response = await deleteOAuthClient(row.client_id)
    if (response.data.success) {
      ElMessage.success('删除成功')
      fetchClientList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error: any) {
    console.error('删除客户端失败:', error)
    ElMessage.error(error.response?.data?.message || '删除失败')
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 解析授权类型
const parseGrantTypes = (grantTypes: string): string[] => {
  if (!grantTypes) return []
  return grantTypes.split(',').map(type => type.trim()).filter(Boolean)
}

// 格式化授权类型显示
const formatGrantType = (type: string): string => {
  const typeMap: Record<string, string> = {
    'authorization_code': '授权码',
    'client_credentials': '客户端凭证',
    'refresh_token': '刷新令牌',
    'password': '密码',
    'implicit': '隐式'
  }
  return typeMap[type] || type
}

// 解析权限范围
const parseScopes = (scopes: string): string[] => {
  if (!scopes) return []
  return scopes.split(',').map(scope => scope.trim()).filter(Boolean)
}

// 格式化日期时间
const formatDateTime = (dateTime: string): string => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 对话框成功回调
const handleDialogSuccess = () => {
  fetchClientList()
}

// 页面初始化
onMounted(() => {
  fetchClientList()
})
</script>

<style scoped>
.oauth-clients-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.search-section {
  margin-bottom: 20px;
  padding: 16px;
  background: var(--el-fill-color-extra-light);
  border-radius: 8px;
}

.table-section {
  margin-bottom: 20px;
}

.client-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.description-tag {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.client-id {
  display: flex;
  align-items: center;
  gap: 8px;
}

.client-id code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: var(--el-fill-color-lighter);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--el-color-primary);
}

.copy-btn {
  padding: 2px;
  min-height: auto;
}

.grant-types,
.scopes {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.grant-type-tag,
.scope-tag {
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: center;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .oauth-clients-container {
    padding: 10px;
  }

  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .search-section :deep(.el-form--inline) {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .search-section :deep(.el-form-item) {
    margin-right: 0;
    margin-bottom: 0;
  }

  .action-buttons {
    flex-direction: column;
    gap: 4px;
  }
}

/* 表格样式优化 */
:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table__header) {
  background: var(--el-fill-color-extra-light);
}

:deep(.el-table__header th) {
  background: var(--el-fill-color-extra-light);
  color: var(--el-text-color-primary);
  font-weight: 600;
}

:deep(.el-table__row:hover > td) {
  background-color: var(--el-fill-color-light);
}

:deep(.el-pagination) {
  --el-pagination-button-color: var(--el-text-color-primary);
  --el-pagination-hover-color: var(--el-color-primary);
}
</style>
