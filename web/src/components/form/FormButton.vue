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
  padding: 16px 24px;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  position: relative;
  overflow: hidden;
  letter-spacing: 0.025em;
}

.form-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.form-btn:hover::before {
  left: 100%;
}

.form-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 
    0 4px 14px 0 rgba(102, 126, 234, 0.39),
    0 0 0 1px rgba(102, 126, 234, 0.1);
}

.form-btn.secondary {
  background: #ffffff;
  color: #374151;
  border: 2px solid #e5e7eb;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.form-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.form-btn.primary:hover:not(:disabled) {
  box-shadow: 
    0 8px 25px 0 rgba(102, 126, 234, 0.5),
    0 0 0 1px rgba(102, 126, 234, 0.2);
}

.form-btn.secondary:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
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

.form-btn.loading {
  cursor: wait;
}

.form-btn span {
  position: relative;
  z-index: 1;
}

@media (max-width: 640px) {
  .form-btn {
    padding: 14px 20px;
    font-size: 16px;
    border-radius: 12px;
  }
}
</style>
