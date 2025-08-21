<template>
  <div class="form-validator">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { provide, reactive } from 'vue'

interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  custom?: (value: any) => boolean | string
}

interface ValidationField {
  value: any
  rules: ValidationRule[]
  errors: string[]
}

interface ValidationState {
  fields: Record<string, ValidationField>
  isValid: boolean
  validateField: (fieldName: string) => boolean
  validateAll: () => boolean
  clearErrors: () => void
}

const validationState = reactive<ValidationState>({
  fields: {},
  isValid: true,
  validateField,
  validateAll,
  clearErrors
})

provide('formValidator', validationState)

function validateField(fieldName: string): boolean {
  const field = validationState.fields[fieldName]
  if (!field) return true

  field.errors = []
  const { value, rules } = field

  for (const rule of rules) {
    // 必填验证
    if (rule.required && (!value || value.trim() === '')) {
      field.errors.push('此字段为必填项')
      continue
    }

    // 长度验证
    if (rule.minLength && value && value.length < rule.minLength) {
      field.errors.push(`最少需要 ${rule.minLength} 个字符`)
    }

    if (rule.maxLength && value && value.length > rule.maxLength) {
      field.errors.push(`最多允许 ${rule.maxLength} 个字符`)
    }

    // 正则验证
    if (rule.pattern && value && !rule.pattern.test(value)) {
      field.errors.push('格式不正确')
    }

    // 自定义验证
    if (rule.custom && value) {
      const result = rule.custom(value)
      if (typeof result === 'string') {
        field.errors.push(result)
      } else if (!result) {
        field.errors.push('验证失败')
      }
    }
  }

  // 更新整体验证状态
  validationState.isValid = Object.values(validationState.fields).every(
    field => field.errors.length === 0
  )

  return field.errors.length === 0
}

function validateAll(): boolean {
  let isValid = true
  for (const fieldName in validationState.fields) {
    if (!validateField(fieldName)) {
      isValid = false
    }
  }
  return isValid
}

function clearErrors(): void {
  for (const fieldName in validationState.fields) {
    validationState.fields[fieldName].errors = []
  }
  validationState.isValid = true
}

// 暴露给子组件使用的方法
defineExpose({
  validateField,
  validateAll,
  clearErrors
})
</script>

<style scoped>
.form-validator {
  width: 100%;
}
</style>
