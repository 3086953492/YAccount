<template>
  <el-form ref="formRef" :model="formModel" :rules="formRules" :label-position="labelPosition" :label-width="labelWidth"
    :size="size" class="form-validator">
    <slot />
  </el-form>
</template>

<script setup lang="ts">
import { ref, provide, reactive, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

interface ValidationRule {
  required?: boolean
  min?: number
  max?: number
  pattern?: RegExp
  validator?: (rule: any, value: any, callback: any) => void
  message?: string
  trigger?: string | string[]
}

interface ValidationField {
  value: any
  rules: ValidationRule[]
  errors: string[]
}

interface ValidationState {
  fields: Record<string, ValidationField>
  isValid: boolean
  validateField: (fieldName: string) => Promise<boolean>
  validateAll: () => Promise<boolean>
  clearErrors: () => void
  resetFields: () => void
}

interface Props {
  modelValue?: Record<string, any>
  rules?: FormRules
  labelPosition?: 'left' | 'right' | 'top'
  labelWidth?: string | number
  size?: 'large' | 'default' | 'small'
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({}),
  rules: () => ({}),
  labelPosition: 'top',
  labelWidth: 'auto',
  size: 'large'
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
  'validate': [isValid: boolean]
}>()

const formRef = ref<FormInstance>()
const formModel = ref(props.modelValue)
const formRules = ref(props.rules)

const validationState = reactive<ValidationState>({
  fields: {},
  isValid: true,
  validateField,
  validateAll,
  clearErrors,
  resetFields
})

provide('formValidator', validationState)

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  formModel.value = newValue
}, { deep: true })

// 监听rules变化
watch(() => props.rules, (newRules) => {
  formRules.value = newRules
}, { deep: true })

async function validateField(fieldName: string): Promise<boolean> {
  if (!formRef.value) return true

  try {
    await formRef.value.validateField(fieldName)
    return true
  } catch {
    return false
  }
}

async function validateAll(): Promise<boolean> {
  if (!formRef.value) return true

  try {
    await formRef.value.validate()
    validationState.isValid = true
    emit('validate', true)
    return true
  } catch {
    validationState.isValid = false
    emit('validate', false)
    return false
  }
}

function clearErrors(): void {
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  validationState.isValid = true
}

function resetFields(): void {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  validationState.isValid = true
}

// 暴露给父组件使用的方法
defineExpose({
  validateField,
  validateAll,
  clearErrors,
  resetFields,
  formRef
})
</script>

<style scoped>
.form-validator {
  width: 100%;
}

:deep(.el-form-item__label) {
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 14px;
}

:deep(.el-form-item__error) {
  color: var(--el-color-danger);
  font-size: 13px;
  margin-top: 4px;
  font-weight: 500;
}
</style>
