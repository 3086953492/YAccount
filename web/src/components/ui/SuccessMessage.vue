<template>
  <div v-if="message" class="success-message">
    <div class="success-content">
      <div class="success-icon">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <polyline points="20,6 9,17 4,12" stroke="currentColor" stroke-width="2"/>
        </svg>
      </div>
      <span class="success-text">{{ message }}</span>
    </div>
    <button v-if="dismissible" @click="$emit('dismiss')" class="dismiss-btn" aria-label="关闭">
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2"/>
        <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2"/>
      </svg>
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
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px 20px;
  border-radius: 12px;
  margin: 16px 0;
  font-size: 14px;
  line-height: 1.5;
  background-color: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #16a34a;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.success-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  flex: 1;
}

.success-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-top: 1px;
}

.success-icon svg {
  width: 20px;
  height: 20px;
}

.success-text {
  flex: 1;
  font-weight: 500;
}

.dismiss-btn {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 4px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
  flex-shrink: 0;
  margin-left: 12px;
}

.dismiss-btn:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.dismiss-btn svg {
  width: 16px;
  height: 16px;
}

@media (max-width: 640px) {
  .success-message {
    padding: 14px 16px;
    border-radius: 10px;
  }
  
  .success-content {
    gap: 10px;
  }
}
</style>
