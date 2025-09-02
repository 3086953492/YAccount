<template>
  <div class="oauth-client-form-container">
    <el-card class="form-card">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon size="20">
              <Key />
            </el-icon>
            <span>{{ isEdit ? '编辑' : '新建' }}OAuth客户端</span>
          </div>
          <el-button @click="handleBack">
            <el-icon>
              <Back />
            </el-icon>
            返回
          </el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" class="oauth-form"
        @submit.prevent>
        <!-- 基本信息 -->
        <div class="form-section">
          <div class="section-title">
            <el-icon>
              <InfoFilled />
            </el-icon>
            <span>基本信息</span>
          </div>

          <el-form-item label="客户端名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入客户端名称" maxlength="255" show-word-limit clearable />
          </el-form-item>

          <el-form-item label="客户端描述" prop="description">
            <el-input v-model="formData.description" type="textarea" placeholder="请输入客户端描述（可选）" :rows="3"
              maxlength="500" show-word-limit clearable />
          </el-form-item>

          <el-form-item label="客户端类型" prop="client_type">
            <el-radio-group v-model="formData.client_type">
              <el-radio value="confidential">
                <div class="radio-option">
                  <div class="radio-title">机密型（Confidential）</div>
                  <div class="radio-desc">能够安全保存客户端密钥的应用，如服务器端应用</div>
                </div>
              </el-radio>
              <el-radio value="public">
                <div class="radio-option">
                  <div class="radio-title">公开型（Public）</div>
                  <div class="radio-desc">无法安全保存客户端密钥的应用，如移动应用或SPA</div>
                </div>
              </el-radio>
            </el-radio-group>
          </el-form-item>
        </div>

        <!-- 授权配置 -->
        <div class="form-section">
          <div class="section-title">
            <el-icon>
              <Lock />
            </el-icon>
            <span>授权配置</span>
          </div>

          <el-form-item label="授权类型" prop="grant_types">
            <el-checkbox-group v-model="formData.grant_types">
              <el-checkbox value="authorization_code">
                <div class="checkbox-option">
                  <div class="checkbox-title">授权码模式（Authorization Code）</div>
                  <div class="checkbox-desc">最安全的授权模式，适用于有后端的应用</div>
                </div>
              </el-checkbox>
              <el-checkbox value="client_credentials">
                <div class="checkbox-option">
                  <div class="checkbox-title">客户端凭证模式（Client Credentials）</div>
                  <div class="checkbox-desc">适用于客户端本身需要访问资源的场景</div>
                </div>
              </el-checkbox>
              <el-checkbox value="refresh_token">
                <div class="checkbox-option">
                  <div class="checkbox-title">刷新令牌（Refresh Token）</div>
                  <div class="checkbox-desc">允许使用刷新令牌获取新的访问令牌</div>
                </div>
              </el-checkbox>
              <el-checkbox value="password" v-if="formData.client_type === 'confidential'">
                <div class="checkbox-option">
                  <div class="checkbox-title">密码模式（Password）</div>
                  <div class="checkbox-desc">直接使用用户名密码获取令牌（不推荐）</div>
                </div>
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="权限范围" prop="scopes">
            <el-checkbox-group v-model="formData.scopes">
              <el-checkbox v-for="scope in availableScopes" :key="scope.value" :value="scope.value">
                <div class="checkbox-option">
                  <div class="checkbox-title">{{ scope.label }}</div>
                  <div class="checkbox-desc">{{ scope.description }}</div>
                </div>
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="重定向URI" prop="redirect_uris">
            <div class="redirect-uris-container">
              <div v-for="(uri, index) in formData.redirect_uris" :key="index" class="uri-item">
                <el-input v-model="formData.redirect_uris[index]" placeholder="https://example.com/callback"
                  clearable />
                <el-button type="danger" text @click="removeRedirectUri(index)"
                  :disabled="formData.redirect_uris.length <= 1">
                  <el-icon>
                    <Delete />
                  </el-icon>
                </el-button>
              </div>
              <el-button type="primary" text @click="addRedirectUri">
                <el-icon>
                  <Plus />
                </el-icon>
                添加重定向URI
              </el-button>
            </div>
          </el-form-item>
        </div>



        <!-- 表单操作 -->
        <div class="form-actions">
          <el-button @click="handleBack">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  createOAuthClient,
  updateOAuthClient,
  getOAuthClient,
  type CreateOAuthClientRequest,
  type UpdateOAuthClientRequest,
  type OAuthClient
} from '@/api/oauth'
import {
  Key,
  Back,
  InfoFilled,
  Lock,
  Plus,
  Delete
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = computed(() => !!route.params.clientId)
const clientId = computed(() => route.params.clientId as string)

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<CreateOAuthClientRequest>({
  name: '',
  description: '',
  redirect_uris: [''],
  grant_types: ['authorization_code'],
  scopes: ['read'],
  client_type: 'confidential'
})



// 可用的权限范围
const availableScopes = ref([
  { value: 'read', label: '读取', description: '读取用户基本信息' },
  { value: 'write', label: '写入', description: '修改用户信息' },
  { value: 'profile', label: '用户资料', description: '访问用户详细资料' },
  { value: 'email', label: '邮箱', description: '访问用户邮箱地址' },
  { value: 'admin', label: '管理员', description: '管理员权限' }
])

// 提交状态
const submitLoading = ref(false)

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入客户端名称', trigger: 'blur' },
    { max: 255, message: '客户端名称长度不能超过255个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述长度不能超过500个字符', trigger: 'blur' }
  ],
  client_type: [
    { required: true, message: '请选择客户端类型', trigger: 'change' }
  ],
  grant_types: [
    {
      type: 'array',
      required: true,
      message: '请至少选择一种授权类型',
      trigger: 'change',
      min: 1
    }
  ],
  scopes: [
    {
      type: 'array',
      required: true,
      message: '请至少选择一个权限范围',
      trigger: 'change',
      min: 1
    }
  ],
  redirect_uris: [
    {
      type: 'array',
      required: true,
      message: '请至少添加一个重定向URI',
      trigger: 'change',
      min: 1
    },
    {
      validator: (rule, value, callback) => {
        const urlPattern = /^https?:\/\/.+/
        const invalidUris = value.filter((uri: string) => uri && !urlPattern.test(uri))
        if (invalidUris.length > 0) {
          callback(new Error('重定向URI必须是有效的HTTP/HTTPS地址'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 添加重定向URI
const addRedirectUri = () => {
  formData.redirect_uris.push('')
}

// 删除重定向URI
const removeRedirectUri = (index: number) => {
  if (formData.redirect_uris.length > 1) {
    formData.redirect_uris.splice(index, 1)
  }
}



// 返回列表页
const handleBack = () => {
  router.push('/admin/oauth/clients')
}

// 加载客户端数据（编辑模式）
const loadClientData = async () => {
  if (!isEdit.value) return

  try {
    const response = await getOAuthClient(clientId.value)
    if (response.data.success) {
      const client: OAuthClient = response.data.data

      // 填充表单数据
      formData.name = client.name
      formData.description = client.description
      formData.client_type = client.client_type as 'public' | 'confidential'
      formData.grant_types = client.grant_types.split(',').map(s => s.trim()).filter(Boolean)
      formData.scopes = client.scopes.split(',').map(s => s.trim()).filter(Boolean)
      formData.redirect_uris = JSON.parse(client.redirect_uris || '[""]')
    } else {
      ElMessage.error(response.data.message || '加载客户端数据失败')
      handleBack()
    }
  } catch (error: any) {
    console.error('加载客户端数据失败:', error)
    ElMessage.error(error.response?.data?.message || '加载客户端数据失败')
    handleBack()
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    submitLoading.value = true

    // 过滤空的重定向URI
    const cleanRedirectUris = formData.redirect_uris.filter(uri => uri.trim())

    const submitData = {
      name: formData.name,
      description: formData.description,
      redirect_uris: cleanRedirectUris,
      grant_types: formData.grant_types,
      scopes: formData.scopes,
      client_type: formData.client_type
    }

    let response
    if (isEdit.value) {
      response = await updateOAuthClient(clientId.value, submitData as UpdateOAuthClientRequest)
    } else {
      response = await createOAuthClient(submitData as CreateOAuthClientRequest)
    }

    if (response.data.success) {
      ElMessage.success(isEdit.value ? '客户端更新成功' : '客户端创建成功')
      handleBack()
    } else {
      ElMessage.error(response.data.message || (isEdit.value ? '更新失败' : '创建失败'))
    }
  } catch (error: any) {
    console.error('提交失败:', error)
    ElMessage.error(error.response?.data?.message || (isEdit.value ? '更新失败' : '创建失败'))
  } finally {
    submitLoading.value = false
  }
}

// 页面初始化
onMounted(() => {
  if (isEdit.value) {
    loadClientData()
  }
})
</script>

<style scoped>
.oauth-client-form-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.form-card {
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

.oauth-form {
  padding: 20px 0;
}

.form-section {
  margin-bottom: 32px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 20px;
  padding-bottom: 8px;
  border-bottom: 2px solid var(--el-color-primary);
}

.radio-option,
.checkbox-option {
  margin-left: 8px;
}

.radio-title,
.checkbox-title {
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 2px;
}

.radio-desc,
.checkbox-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.scopes-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}



.redirect-uris-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.uri-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .oauth-client-form-container {
    padding: 10px;
  }

  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .uri-item {
    flex-direction: column;
    align-items: stretch;
  }

  .form-actions {
    flex-direction: column;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-radio) {
  align-items: flex-start;
  margin-bottom: 16px;
}

:deep(.el-checkbox) {
  align-items: flex-start;
  margin-bottom: 12px;
}

:deep(.el-radio__input) {
  margin-top: 2px;
}

:deep(.el-checkbox__input) {
  margin-top: 2px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-card__header) {
  background: var(--el-fill-color-extra-light);
}
</style>
