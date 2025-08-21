<template>
  <div v-if="loading" class="loading-spinner" :class="size">
    <div class="spinner"></div>
    <p v-if="text" class="loading-text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
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

.loading-spinner.small .spinner {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

.loading-spinner.small .loading-text {
  font-size: 12px;
  margin-top: 8px;
}

.loading-spinner.medium .spinner {
  width: 40px;
  height: 40px;
  border-width: 4px;
}

.loading-spinner.large .spinner {
  width: 60px;
  height: 60px;
  border-width: 6px;
}

.spinner {
  border: 4px solid #f3f4f6;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s cubic-bezier(0.68, -0.55, 0.265, 1.55) infinite;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.loading-text {
  margin-top: 16px;
  color: #6b7280;
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
  .loading-spinner.small .spinner {
    width: 14px;
    height: 14px;
  }
  
  .loading-spinner.medium .spinner {
    width: 32px;
    height: 32px;
  }
  
  .loading-spinner.large .spinner {
    width: 48px;
    height: 48px;
  }
}
</style>
