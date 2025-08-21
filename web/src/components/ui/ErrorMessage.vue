<template>
  <div v-if="message" class="error-message" :class="type">
    <div class="error-content">
      <div class="error-icon">
        <svg v-if="type === 'error'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2"/>
          <line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2"/>
        </svg>
        <svg v-else-if="type === 'warning'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" stroke="currentColor" stroke-width="2"/>
          <line x1="12" y1="9" x2="12" y2="13" stroke="currentColor" stroke-width="2"/>
          <circle cx="12" cy="17" r="1" fill="currentColor"/>
        </svg>
        <svg v-else-if="type === 'info'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <line x1="12" y1="16" x2="12" y2="12" stroke="currentColor" stroke-width="2"/>
          <line x1="12" y1="8" x2="12.01" y2="8" stroke="currentColor" stroke-width="2"/>
        </svg>
        <svg v-else-if="type === 'success'" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <polyline points="20,6 9,17 4,12" stroke="currentColor" stroke-width="2"/>
        </svg>
      </div>
      <span class="error-text">{{ message }}</span>
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
</script>

<style scoped>
.error-message {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px 20px;
  border-radius: 12px;
  margin: 16px 0;
  font-size: 14px;
  line-height: 1.5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  border: 1px solid transparent;
}

.error-message.error {
  background-color: #fef2f2;
  border-color: #fecaca;
  color: #dc2626;
}

.error-message.warning {
  background-color: #fffbeb;
  border-color: #fed7aa;
  color: #d97706;
}

.error-message.info {
  background-color: #eff6ff;
  border-color: #bfdbfe;
  color: #2563eb;
}

.error-message.success {
  background-color: #f0fdf4;
  border-color: #bbf7d0;
  color: #16a34a;
}

.error-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  flex: 1;
}

.error-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-top: 1px;
}

.error-icon svg {
  width: 20px;
  height: 20px;
}

.error-text {
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
  .error-message {
    padding: 14px 16px;
    border-radius: 10px;
  }
  
  .error-content {
    gap: 10px;
  }
}
</style>
