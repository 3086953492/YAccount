<template>
  <div class="form-footer">
    <span class="footer-text">{{ text }}</span>
    <el-link :underline="false" type="primary" @click="$router.push(linkTo)" class="footer-link">
      {{ linkText }}
    </el-link>
  </div>
</template>

<script setup lang="ts">
interface Props {
  text: string
  linkText: string
  linkTo: string
}

defineProps<Props>()
</script>

<style scoped>
.form-footer {
  text-align: center;
  padding: 16px 0 0 0;
  border-top: 1px solid var(--el-border-color-lighter);
  margin-top: 4px;
}

.footer-text {
  color: var(--el-text-color-regular);
  font-size: 14px;
  margin-right: 8px;
}

.footer-link {
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
}

.footer-link:hover {
  transform: translateY(-1px);
}

@media (max-width: 640px) {
  .form-footer {
    padding: 16px 0 0 0;
  }
}
</style>
