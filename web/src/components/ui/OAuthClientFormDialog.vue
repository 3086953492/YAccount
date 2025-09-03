<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑OAuth客户端' : '新建OAuth客户端'"
    width="700px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      class="oauth-dialog-form"
    >
      <!-- 基本信息 -->
      <div class="form-section">
        <div class="section-title">基本信息</div>
        
        <el-form-item label="客户端名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入客户端名称"
            maxlength="255"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="客户端描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            placeholder="请输入客户端描述（可选）"
            :rows="2"
            maxlength="500"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="客户端类型" prop="client_type">
          <el-radio-group v-model="formData.client_type">
            <el-radio value="confidential">机密型</el-radio>
            <el-radio value="public">公开型</el-radio>
          </el-radio-group>
          <div class="form-tip">
            机密型：能安全保存密钥的应用；公开型：无法安全保存密钥的应用
          </div>
        </el-form-item>
      </div>

      <!-- 授权配置 -->
      <div class="form-section">
        <div class="section-title">授权配置</div>
        
        <el-form-item label="授权类型" prop="grant_types">
          <el-checkbox-group v-model="formData.grant_types">
            <el-checkbox value="authorization_code">授权码模式</el-checkbox>
            <el-checkbox value="client_credentials">客户端凭证模式</el-checkbox>
            <el-checkbox value="refresh_token">刷新令牌</el-checkbox>
            <el-checkbox value="password" v-if="formData.client_type === 'confidential'">
              密码模式
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="权限范围" prop="scopes">
          <el-checkbox-group v-model="formData.scopes">
            <el-checkbox
              v-for="scope in availableScopes"
              :key="scope.value"
              :value="scope.value"
            >
              {{ scope.label }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="重定向URI" prop="redirect_uris">
          <div class="redirect-uris-container">
            <div
              v-for="(uri, index) in formData.redirect_uris"
              :key="index"
              class="uri-item"
            >
              <el-input
                v-model="formData.redirect_uris[index]"
                placeholder="https://example.com/callback"
                clearable
              />
              <el-button
                type="danger"
                text
                @click="removeRedirectUri(index)"
                :disabled="formData.redirect_uris.length <= 1"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button type="primary" text @click="addRedirectUri">
              <el-icon><Plus /></el-icon>
              添加重定向URI
            </el-button>
          </div>
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 客户端密钥显示弹窗 -->
  <ClientSecretDisplay
    v-model:visible="secretDialogVisible"
    :client-data="createdClient"
    @confirmed="handleSecretConfirmed"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { 
  createOAuthClient, 
  updateOAuthClient,
  type CreateOAuthClientRequest,
  type UpdateOAuthClientRequest,
  type OAuthClient
} from '@/api/oauth'
import { Plus, Delete } from '@element-plus/icons-vue'
import ClientSecretDisplay from './ClientSecretDisplay.vue'

interface Props {
  visible: boolean
  clientData?: OAuthClient | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  clientData: null
})

const emit = defineEmits<Emits>()

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 是否为编辑模式
const isEdit = computed(() => !!props.clientData)

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
  { value: 'read', label: '读取' },
  { value: 'write', label: '写入' },
  { value: 'profile', label: '用户资料' },
  { value: 'email', label: '邮箱' },
  { value: 'admin', label: '管理员' }
])

// 提交状态
const submitLoading = ref(false)

// 密钥显示状态
const secretDialogVisible = ref(false)
const createdClient = ref<OAuthClient | null>(null)

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

// 监听客户端数据变化，填充表单
watch(() => props.clientData, (clientData) => {
  if (clientData) {
    formData.name = clientData.name
    formData.description = clientData.description
    formData.client_type = clientData.client_type as 'public' | 'confidential'
    formData.grant_types = clientData.grant_types.split(',').map(s => s.trim()).filter(Boolean)
    formData.scopes = clientData.scopes.split(',').map(s => s.trim()).filter(Boolean)
    formData.redirect_uris = JSON.parse(clientData.redirect_uris || '[""]')
  }
}, { immediate: true })

// 监听对话框显示状态，重置表单
watch(dialogVisible, (visible) => {
  if (visible && !isEdit.value) {
    resetForm()
  }
})

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.description = ''
  formData.redirect_uris = ['']
  formData.grant_types = ['authorization_code']
  formData.scopes = ['read']
  formData.client_type = 'confidential'
  
  if (formRef.value) {
    formRef.value.resetFields()
  }
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

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}

// 密钥确认处理
const handleSecretConfirmed = () => {
  secretDialogVisible.value = false
  createdClient.value = null
  emit('success')
  resetForm()
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
    if (isEdit.value && props.clientData) {
      response = await updateOAuthClient(props.clientData.client_id, submitData as UpdateOAuthClientRequest)
    } else {
      response = await createOAuthClient(submitData as CreateOAuthClientRequest)
    }

    if (response.data.success) {
      if (isEdit.value) {
        ElMessage.success('客户端更新成功')
        emit('success')
        handleClose()
      } else {
        // 创建成功，显示密钥
        createdClient.value = response.data.data
        secretDialogVisible.value = true
        dialogVisible.value = false // 关闭表单对话框
      }
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
</script>

<style scoped>
.oauth-dialog-form {
  max-height: 60vh;
  overflow-y: auto;
}

.form-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 16px;
  padding-bottom: 4px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
  line-height: 1.4;
}

.redirect-uris-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.uri-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Element Plus 样式覆盖 */
:deep(.el-checkbox) {
  margin-bottom: 8px;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

:deep(.el-dialog__header) {
  padding: 20px 20px 10px;
}

:deep(.el-dialog__body) {
  padding: 10px 20px 20px;
}

:deep(.el-dialog__footer) {
  padding: 10px 20px 20px;
}
</style>
