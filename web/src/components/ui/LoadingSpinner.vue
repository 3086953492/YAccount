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
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  margin-top: 16px;
  color: #666;
  font-size: 14px;
  text-align: center;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}
</style>
