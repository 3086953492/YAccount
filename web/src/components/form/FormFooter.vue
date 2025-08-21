<template>
  <div class="form-footer">
    <span class="footer-text">{{ text }}</span>
    <router-link :to="linkTo" class="footer-link">{{ linkText }}</router-link>
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
  padding: 20px 0 0 0;
  border-top: 1px solid #f3f4f6;
  margin-top: 8px;
}

.footer-text {
  color: #6b7280;
  font-size: 14px;
  margin-right: 8px;
}

.footer-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
  position: relative;
}

.footer-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.footer-link:hover {
  color: #5a67d8;
}

.footer-link:hover::after {
  width: 100%;
}

@media (max-width: 640px) {
  .form-footer {
    padding: 16px 0 0 0;
  }
}
</style>
