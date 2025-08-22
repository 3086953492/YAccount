<template>
  <div v-if="loading" class="loading-spinner" :class="size">
    <el-icon class="spinner-icon" :size="iconSize">
      <Loading />
    </el-icon>
    <p v-if="text" class="loading-text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loading } from '@element-plus/icons-vue'

interface Props {
  loading?: boolean
  text?: string
  size?: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  loading: true,
  text: '',
  size: 'medium'
})

const iconSize = computed(() => {
  switch (props.size) {
    case 'small': return 16
    case 'large': return 60
    default: return 40
  }
})
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.loading-spinner.small {
  padding: 8px;
}

.loading-spinner.small .loading-text {
  font-size: 12px;
  margin-top: 8px;
}

.loading-spinner.large .loading-text {
  font-size: 16px;
  margin-top: 20px;
}

.spinner-icon {
  color: var(--el-color-primary);
  animation: spin 1s linear infinite;
}

.loading-text {
  margin-top: 16px;
  color: var(--el-text-color-regular);
  font-size: 14px;
  text-align: center;
  font-weight: 500;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

@media (max-width: 640px) {
  .loading-spinner.small .spinner-icon {
    font-size: 14px;
  }

  .loading-spinner.medium .spinner-icon {
    font-size: 32px;
  }

  .loading-spinner.large .spinner-icon {
    font-size: 48px;
  }
}
</style>
