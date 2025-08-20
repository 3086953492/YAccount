<template>
  <div v-if="message" class="success-message">
    <div class="success-content">
      <span class="success-icon">✅</span>
      <span class="success-text">{{ message }}</span>
    </div>
    <button v-if="dismissible" @click="$emit('dismiss')" class="dismiss-btn">
      ×
    </button>
  </div>
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 16px 0;
  font-size: 14px;
  line-height: 1.5;
  background-color: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #16a34a;
}

.success-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.success-icon {
  font-size: 16px;
}

.success-text {
  flex: 1;
}

.dismiss-btn {
  background: none;
  border: none;
  color: inherit;
  font-size: 18px;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}

.dismiss-btn:hover {
  background-color: rgba(0, 0, 0, 0.1);
}
</style>
