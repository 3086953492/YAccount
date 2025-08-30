<template>
  <div class="system-info-editor">
    <el-dialog
      v-model="dialogVisible"
      title="编辑系统信息"
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
          :label="item.description || item.config_key"
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
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存
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

// 表单验证规则
const formRules: FormRules = {
  items: {} as Record<string, any>
}

// 初始化表单数据
const initFormData = () => {
  editableItems.value.forEach(item => {
    formData.items[item.id] = {
      id: item.id,
      config_value: item.config_value
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
    
    loading.value = true
    
    // 转换为后端需要的格式
    const updateData = {
      system_infos: Object.values(formData.items)
    }
    
    const success = await systemStore.updateSystemInfo(updateData)
    
    if (success) {
      ElMessage.success('系统信息更新成功')
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
