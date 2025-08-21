<template>
  <div class="form-footer">
    <p>{{ text }} <router-link :to="linkTo">{{ linkText }}</router-link></p>
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
  margin-top: 20px;
}

.form-footer p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.form-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.form-footer a:hover {
  text-decoration: underline;
}
</style>
