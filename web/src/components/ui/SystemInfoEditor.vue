<template>
  <div class="system-info-editor">
    <el-dialog
      v-model="dialogVisible"
      :title="`编辑系统信息${hasChanges ? ` (${getChangedFields().length} 项已修改)` : ''}`"
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        class="editor-form"
      >
        <el-form-item
          v-for="item in editableItems"
          :key="item.id"
          :label="`${item.description || item.config_key}${isFieldChanged(item.id) ? ' *' : ''}`"
          :prop="`items.${item.id}.config_value`"
        >
          <el-input
            v-if="item.config_type === 'text' || item.config_type === 'string'"
            v-model="formData.items[item.id].config_value"
            :placeholder="`请输入${item.description || item.config_key}`"
            type="text"
            maxlength="255"
            show-word-limit
          />
          <el-input
            v-else-if="item.config_type === 'textarea'"
            v-model="formData.items[item.id].config_value"
            :placeholder="`请输入${item.description || item.config_key}`"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
          />
          <el-input
            v-else-if="item.config_type === 'url'"
            v-model="formData.items[item.id].config_value"
            :placeholder="`请输入${item.description || item.config_key}的URL`"
            type="url"
          />
          <el-input
            v-else
            v-model="formData.items[item.id].config_value"
            :placeholder="`请输入${item.description || item.config_key}`"
            type="text"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleCancel">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit" 
            :loading="loading"
            :disabled="!hasChanges"
          >
            保存{{ hasChanges ? ` (${getChangedFields().length})` : '' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useSystemStore } from '@/stores/system'
import type { SystemInfoList, UpdateSystemInfoRequest } from '@/types/system'

interface Props {
  visible: boolean
  systemInfoList: SystemInfoList[]
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const systemStore = useSystemStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 可编辑的项目（只显示启用的配置）
const editableItems = computed(() => {
  return props.systemInfoList.filter(item => item.status === 1)
})

// 表单数据
const formData = reactive<{
  items: Record<number, UpdateSystemInfoRequest>
}>({
  items: {}
})

// 原始值记录，用于检测变更
const originalValues = ref<Record<number, string>>({})

// 表单验证规则
const formRules: FormRules = {
  items: {} as Record<string, any>
}

// 初始化表单数据
const initFormData = () => {
  editableItems.value.forEach(item => {
    const originalValue = item.config_value
    originalValues.value[item.id] = originalValue
    
    formData.items[item.id] = {
      id: item.id,
      config_value: originalValue
    }
    
    // 为每个字段添加验证规则
    if (item.config_type === 'url') {
      formRules.items[item.id] = [
        { required: true, message: '请输入有效的URL', trigger: 'blur' },
        { type: 'url', message: '请输入有效的URL格式', trigger: 'blur' }
      ]
    } else if (item.config_type === 'text' || item.config_type === 'string') {
      formRules.items[item.id] = [
        { required: true, message: '此项不能为空', trigger: 'blur' },
        { min: 1, max: 255, message: '长度在 1 到 255 个字符', trigger: 'blur' }
      ]
    } else if (item.config_type === 'textarea') {
      formRules.items[item.id] = [
        { required: true, message: '此项不能为空', trigger: 'blur' },
        { min: 1, max: 500, message: '长度在 1 到 500 个字符', trigger: 'blur' }
      ]
    } else {
      formRules.items[item.id] = [
        { required: true, message: '此项不能为空', trigger: 'blur' }
      ]
    }
  })
}

// 检测字段是否有变更
const hasChanges = computed(() => {
  return Object.keys(formData.items).some(id => {
    const itemId = parseInt(id)
    return formData.items[itemId].config_value !== originalValues.value[itemId]
  })
})

// 获取变更的字段
const getChangedFields = () => {
  const changedFields: UpdateSystemInfoRequest[] = []
  
  Object.keys(formData.items).forEach(id => {
    const itemId = parseInt(id)
    const currentValue = formData.items[itemId].config_value
    const originalValue = originalValues.value[itemId]
    
    if (currentValue !== originalValue) {
      changedFields.push({
        id: itemId,
        config_value: currentValue
      })
    }
  })
  
  return changedFields
}

// 检测单个字段是否变更
const isFieldChanged = (itemId: number) => {
  return formData.items[itemId]?.config_value !== originalValues.value[itemId]
}

// 监听系统信息列表变化，重新初始化表单
watch(() => props.systemInfoList, () => {
  if (props.visible) {
    initFormData()
  }
}, { immediate: true })

// 监听对话框显示状态
watch(() => props.visible, (visible) => {
  if (visible) {
    initFormData()
  }
})

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 检查是否有变更
    if (!hasChanges.value) {
      ElMessage.warning('没有检测到任何变更，无需保存')
      return
    }
    
    loading.value = true
    
    // 只获取变更的字段
    const changedFields = getChangedFields()
    const updateData = {
      system_infos: changedFields
    }
    

    
    const success = await systemStore.updateSystemInfo(updateData)
    
    if (success) {

      ElMessage.success(`系统信息更新成功，共更新了 ${changedFields.length} 个配置项`)
      emit('success')
      dialogVisible.value = false
    } else {
      ElMessage.error('系统信息更新失败')
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}

// 取消编辑
const handleCancel = () => {
  dialogVisible.value = false
}
</script>

<style scoped>
.system-info-editor {
  .editor-form {
    max-height: 400px;
    overflow-y: auto;
  }
  
  .dialog-footer {
    text-align: right;
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: var(--el-text-color-primary);
  }
  
  :deep(.el-input__wrapper) {
    border-radius: 8px;
  }
  
  :deep(.el-textarea__inner) {
    border-radius: 8px;
  }
}
</style>
