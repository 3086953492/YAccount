<template>
  <div class="main-layout">
    <Navigation />
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import Navigation from './Navigation.vue'
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding-top: 0;
}
</style>
