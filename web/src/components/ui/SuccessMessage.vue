<template>
  <el-alert v-if="message" :title="message" type="success" :closable="dismissible" :show-icon="true" :center="false"
    @close="$emit('dismiss')" class="success-message" />
</template>

<script setup lang="ts">
interface Props {
  message?: string
  dismissible?: boolean
}

interface Emits {
  (e: 'dismiss'): void
}

const props = withDefaults(defineProps<Props>(), {
  dismissible: true
})

defineEmits<Emits>()
</script>

<style scoped>
.success-message {
  margin: 16px 0;
  border-radius: 12px;
}

.success-message :deep(.el-alert__title) {
  font-weight: 500;
  font-size: 14px;
  line-height: 1.5;
}

.success-message :deep(.el-alert__content) {
  padding: 0;
}

.success-message :deep(.el-alert__icon) {
  font-size: 20px;
}

.success-message :deep(.el-alert__closebtn) {
  font-size: 16px;
  color: inherit;
  opacity: 0.7;
}

.success-message :deep(.el-alert__closebtn:hover) {
  opacity: 1;
}

/* 自定义成功样式 */
.success-message :deep(.el-alert) {
  background-color: var(--el-color-success-light-9);
  border-color: var(--el-color-success-light-7);
  color: var(--el-color-success);
}

@media (max-width: 640px) {
  .success-message {
    margin: 14px 0;
  }

  .success-message :deep(.el-alert__title) {
    font-size: 13px;
  }
}
</style>
