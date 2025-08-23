<template>
  <div class="user-list-container">
    <div class="page-header">
      <h1>用户列表</h1>
      <p class="page-description">管理系统中的所有用户账户</p>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <el-form :model="searchForm" inline>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchForm.nickname" placeholder="请输入昵称" clearable @keyup.enter="handleSearch" />
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

    <!-- 用户列表 -->
    <div class="user-table-section">
      <el-table :data="userList" v-loading="loading" stripe border class="user-table">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <UserAvatar :size="40" :avatar="row.avatar" :username="row.username" :nickname="row.nickname" />
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleViewUser(row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { getUserList } from '@/api/user'
import UserAvatar from '@/components/ui/UserAvatar.vue'

interface User {
  id: number
  username: string
  nickname: string
  avatar: string
}

interface Pagination {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

const router = useRouter()

// 响应式数据
const loading = ref(false)
const userList = ref<User[]>([])
const searchForm = reactive({
  username: '',
  nickname: ''
})
const pagination = reactive<Pagination>({
  page: 1,
  pageSize: 20,
  total: 0,
  totalPages: 0
})

// 获取用户列表
const fetchUserList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username || undefined,
      nickname: searchForm.nickname || undefined
    }

    const response = await getUserList(params)
    const data = response.data.data
    userList.value = data.items
    pagination.total = data.total
    pagination.totalPages = data.totalPages
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchUserList()
}

// 重置搜索
const handleReset = () => {
  searchForm.username = ''
  searchForm.nickname = ''
  pagination.page = 1
  fetchUserList()
}

// 分页大小改变
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchUserList()
}

// 当前页改变
const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchUserList()
}

// 查看用户详情
const handleViewUser = (user: User) => {
  router.push(`/user/${user.id}`)
}

// 组件挂载时获取数据
onMounted(() => {
  fetchUserList()
})
</script>

<style scoped>
.user-list-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  margin-top: 24px;
}

.page-header {
  margin-bottom: 32px;
  text-align: center;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 16px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.search-section {
  background: var(--el-bg-color-page);
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 24px;
  border: 1px solid var(--el-border-color-light);
}

.user-table-section {
  background: var(--el-bg-color-page);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-light);
  overflow: hidden;
}

.user-table {
  width: 100%;
}

.pagination-section {
  padding: 20px;
  display: flex;
  justify-content: center;
  background: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-light);
}

@media (max-width: 768px) {
  .user-list-container {
    padding: 16px;
  }

  .search-section {
    padding: 16px;
  }

  .page-header h1 {
    font-size: 24px;
  }
}
</style>
