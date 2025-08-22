<template>
  <el-button :type="buttonType" :disabled="disabled || loading" :loading="loading" :size="size"
    :native-type="nativeType" :class="['form-btn', variant]" @click="$emit('click')">
    <template #default>
      <span>{{ loading ? loadingText : text }}</span>
    </template>
  </el-button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  text: string
  loadingText?: string
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary'
  loading?: boolean
  disabled?: boolean
  size?: 'large' | 'default' | 'small'
}

const props = withDefaults(defineProps<Props>(), {
  loadingText: '',
  type: 'button',
  variant: 'primary',
  loading: false,
  disabled: false,
  size: 'large'
})

const emit = defineEmits<{
  click: []
}>()

const buttonType = computed(() => {
  return props.variant === 'primary' ? 'primary' : 'default'
})

const nativeType = computed(() => {
  return props.type
})
</script>

<style scoped>
.form-btn {
  width: 100%;
  height: 48px;
  border-radius: 14px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.025em;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.form-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow:
    0 4px 14px 0 rgba(102, 126, 234, 0.39),
    0 0 0 1px rgba(102, 126, 234, 0.1);
}

.form-btn.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow:
    0 8px 25px 0 rgba(102, 126, 234, 0.5),
    0 0 0 1px rgba(102, 126, 234, 0.2);
}

.form-btn.secondary {
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
  border: 2px solid var(--el-border-color);
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.form-btn.secondary:hover:not(:disabled) {
  transform: translateY(-2px);
  background: var(--el-fill-color-light);
  border-color: var(--el-border-color-hover);
  box-shadow:
    0 4px 12px 0 rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(0, 0, 0, 0.05);
}

.form-btn:active:not(:disabled) {
  transform: translateY(0);
}

.form-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 640px) {
  .form-btn {
    height: 44px;
    font-size: 16px;
    border-radius: 12px;
  }
}
</style>
