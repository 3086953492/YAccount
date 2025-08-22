<template>
  <el-alert v-if="message" :title="message" :type="alertType" :closable="dismissible" :show-icon="true" :center="false"
    @close="$emit('dismiss')" class="error-message" :class="type" />
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  message?: string
  type?: 'error' | 'warning' | 'info' | 'success'
  dismissible?: boolean
}

interface Emits {
  (e: 'dismiss'): void
}

const props = withDefaults(defineProps<Props>(), {
  type: 'error',
  dismissible: true
})

defineEmits<Emits>()

const alertType = computed(() => {
  return props.type
})
</script>

<style scoped>
.error-message {
  margin: 16px 0;
  border-radius: 12px;
}

.error-message :deep(.el-alert__title) {
  font-weight: 500;
  font-size: 14px;
  line-height: 1.5;
}

.error-message :deep(.el-alert__content) {
  padding: 0;
}

.error-message :deep(.el-alert__icon) {
  font-size: 20px;
}

.error-message :deep(.el-alert__closebtn) {
  font-size: 16px;
  color: inherit;
  opacity: 0.7;
}

.error-message :deep(.el-alert__closebtn:hover) {
  opacity: 1;
}

/* 自定义类型样式 */
.error-message.error :deep(.el-alert) {
  background-color: var(--el-color-danger-light-9);
  border-color: var(--el-color-danger-light-7);
  color: var(--el-color-danger);
}

.error-message.warning :deep(.el-alert) {
  background-color: var(--el-color-warning-light-9);
  border-color: var(--el-color-warning-light-7);
  color: var(--el-color-warning);
}

.error-message.info :deep(.el-alert) {
  background-color: var(--el-color-info-light-9);
  border-color: var(--el-color-info-light-7);
  color: var(--el-color-info);
}

.error-message.success :deep(.el-alert) {
  background-color: var(--el-color-success-light-9);
  border-color: var(--el-color-success-light-7);
  color: var(--el-color-success);
}

@media (max-width: 640px) {
  .error-message {
    margin: 14px 0;
  }

  .error-message :deep(.el-alert__title) {
    font-size: 13px;
  }
}
</style>
