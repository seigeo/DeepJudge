<template>
  <div class="app-container">
    <router-view />
    <Footer />
  </div>
</template>

<script setup>
import Footer from './components/Footer.vue'
</script>

<style>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 让路由视图占据除页脚外的所有空间 */
.app-container > :first-child {
  flex: 1;
}

/* 重置默认样式 */
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

/* 可加一些全局样式 */
</style>