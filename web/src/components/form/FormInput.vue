<template>
  <el-form-item :label="label" :prop="prop" :error="errorMessage">
    <el-input v-model="inputValue" :type="type" :placeholder="placeholder" :disabled="disabled" :clearable="clearable"
      :show-password="type === 'password'" :prefix-icon="prefixIcon" :suffix-icon="suffixIcon" @input="handleInput"
      @blur="handleBlur" @focus="handleFocus" @clear="handleClear" size="large" />
    <template #label v-if="hint">
      <span>{{ label }}</span>
      <el-tooltip v-if="hint" :content="hint" placement="top">
        <el-icon class="hint-icon">
          <QuestionFilled />
        </el-icon>
      </el-tooltip>
    </template>
  </el-form-item>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { QuestionFilled } from '@element-plus/icons-vue'

interface Props {
  modelValue: string
  label: string
  prop?: string
  type?: string
  placeholder?: string
  disabled?: boolean
  errorMessage?: string
  hint?: string
  clearable?: boolean
  prefixIcon?: string
  suffixIcon?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  placeholder: '',
  disabled: false,
  errorMessage: '',
  hint: '',
  clearable: true,
  prefixIcon: '',
  suffixIcon: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'input': [value: string]
  'blur': [value: string]
  'focus': [value: string]
  'clear': []
}>()

const inputValue = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})

const handleInput = (value: string) => {
  emit('input', value)
}

const handleBlur = (event: FocusEvent) => {
  const target = event.target as HTMLInputElement
  emit('blur', target.value)
}

const handleFocus = (event: FocusEvent) => {
  const target = event.target as HTMLInputElement
  emit('focus', target.value)
}

const handleClear = () => {
  emit('clear')
}
</script>

<style scoped>
.hint-icon {
  margin-left: 4px;
  color: var(--el-text-color-secondary);
  cursor: help;
}

:deep(.el-form-item__label) {
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 14px;
}

:deep(.el-input__wrapper) {
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--el-border-color-hover);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-color-primary);
}

:deep(.el-input__inner) {
  font-size: 16px;
  padding: 14px 18px;
}

:deep(.el-form-item__error) {
  color: var(--el-color-danger);
  font-size: 13px;
  margin-top: 4px;
  font-weight: 500;
}

@media (max-width: 640px) {
  :deep(.el-input__inner) {
    padding: 14px 16px;
    font-size: 16px;
  }
}
</style>
