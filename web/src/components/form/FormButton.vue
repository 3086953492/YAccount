<template>
  <button :type="type" :disabled="disabled || loading" :class="['form-btn', variant, { 'loading': loading }]"
    @click="$emit('click')">
    <LoadingSpinner v-if="loading" size="small" />
    <span v-else>{{ loading ? loadingText : text }}</span>
  </button>
</template>

<script setup lang="ts">
import LoadingSpinner from '../ui/LoadingSpinner.vue'

interface Props {
  text: string
  loadingText?: string
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary'
  loading?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loadingText: '',
  type: 'button',
  variant: 'primary',
  loading: false,
  disabled: false
})

defineEmits<{
  click: []
}>()
</script>

<style scoped>
.form-btn {
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.form-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.form-btn.secondary {
  background: #f8f9fa;
  color: #333;
  border: 2px solid #e1e5e9;
}

.form-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.form-btn.primary:hover:not(:disabled) {
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
}

.form-btn.secondary:hover:not(:disabled) {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
}

.form-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.form-btn.loading {
  cursor: wait;
}
</style>
