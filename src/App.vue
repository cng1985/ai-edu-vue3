<script setup>
import { ref } from 'vue'
import AppSidebar from './components/AppSidebar.vue'

const sidebarOpen = ref(false)
</script>

<template>
  <div class="layout">
    <button class="layout__menu-btn" @click="sidebarOpen = !sidebarOpen" aria-label="切换菜单">
      ☰
    </button>
    <AppSidebar :open="sidebarOpen" @navigate="sidebarOpen = false" />
    <div v-if="sidebarOpen" class="layout__backdrop" @click="sidebarOpen = false"></div>
    <main class="layout__main">
      <router-view v-slot="{ Component }">
        <component :is="Component" />
      </router-view>
    </main>
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  min-height: 100vh;
}

.layout__main {
  flex: 1;
  min-width: 0;
  margin-left: var(--sidebar-width);
}

.layout__menu-btn {
  display: none;
  position: fixed;
  top: 14px;
  left: 14px;
  z-index: 60;
  width: 40px;
  height: 40px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--surface);
  font-size: 17px;
  cursor: pointer;
  box-shadow: var(--shadow);
}

.layout__backdrop {
  display: none;
}

@media (max-width: 860px) {
  .layout__main {
    margin-left: 0;
    padding-top: 52px;
  }

  .layout__menu-btn {
    display: block;
  }

  .layout__backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(15, 23, 42, 0.4);
    z-index: 40;
  }
}
</style>
