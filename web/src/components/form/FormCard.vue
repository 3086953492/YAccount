<template>
  <div class="form-container">
    <el-card class="form-card" shadow="hover">
      <template #header>
        <div class="form-header">
          <div class="logo-container">
            <el-icon class="logo-icon" size="32">
              <Lock />
            </el-icon>
          </div>
          <h1 class="title">{{ title }}</h1>
          <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
        </div>
      </template>

      <div class="form-content">
        <slot />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Lock } from '@element-plus/icons-vue'

interface Props {
  title: string
  subtitle?: string
}

defineProps<Props>()
</script>

<style scoped>
.form-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-bg-color-page);
  padding: 20px;
  position: relative;
}

.form-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%23ffffff" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%23ffffff" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="%23ffffff" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="%23ffffff" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="%23ffffff" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  pointer-events: none;
  opacity: 0.3;
}

.form-card {
  width: 100%;
  max-width: 600px;
  border-radius: 24px;
  position: relative;
  z-index: 1;
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.form-card :deep(.el-card__header) {
  padding: 36px 36px 24px 36px;
  border-bottom: none;
  background: transparent;
}

.form-card :deep(.el-card__body) {
  padding: 0 36px 36px 36px;
}

.logo-container {
  text-align: center;
  margin-bottom: 24px;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: white;
  margin-bottom: 16px;
}

.form-header {
  text-align: center;
}

.title {
  color: var(--el-text-color-primary);
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.025em;
}

.subtitle {
  color: var(--el-text-color-regular);
  margin: 0;
  font-size: 16px;
  line-height: 1.5;
}

.form-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (max-width: 640px) {
  .form-container {
    padding: 15px;
  }

  .form-card {
    max-width: 100%;
    margin: 0;
  }

  .form-card :deep(.el-card__header) {
    padding: 32px 24px 20px 24px;
  }

  .form-card :deep(.el-card__body) {
    padding: 0 24px 32px 24px;
  }

  .title {
    font-size: 28px;
  }

  .logo-icon {
    width: 56px;
    height: 56px;
  }
}

@media (max-width: 480px) {
  .form-card :deep(.el-card__header) {
    padding: 24px 20px 16px 20px;
  }

  .form-card :deep(.el-card__body) {
    padding: 0 20px 24px 20px;
  }

  .title {
    font-size: 24px;
  }

  .subtitle {
    font-size: 14px;
  }
}
</style>
