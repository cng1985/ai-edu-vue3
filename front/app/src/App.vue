<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from './components/AppSidebar.vue'

const route = useRoute()
const sidebarOpen = ref(false)
</script>

<template>
  <!-- 登录/注册等页面使用全屏空白布局 -->
  <router-view v-if="route.meta.blank" />

  <div v-else class="layout">
    <button class="layout__menu-btn" @click="sidebarOpen = !sidebarOpen" aria-label="切换菜单">
      <span></span><span></span><span></span>
    </button>
    <AppSidebar :open="sidebarOpen" @navigate="sidebarOpen = false" />
    <div v-if="sidebarOpen" class="layout__backdrop" @click="sidebarOpen = false"></div>
    <main class="layout__main">
      <div class="layout__glow" aria-hidden="true"></div>
      <router-view v-slot="{ Component }">
        <transition name="route" mode="out-in">
          <component :is="Component" :key="route.path" />
        </transition>
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
  position: relative;
  flex: 1;
  min-width: 0;
  margin-left: var(--sidebar-width);
  isolation: isolate;
}

.layout__glow {
  position: absolute;
  inset: 0 0 auto 0;
  height: 420px;
  z-index: -1;
  pointer-events: none;
  background:
    radial-gradient(60% 80% at 15% 0%, rgba(107, 92, 255, 0.12), transparent 60%),
    radial-gradient(40% 60% at 85% 10%, rgba(15, 185, 129, 0.08), transparent 60%);
}

.route-enter-active,
.route-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.route-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.route-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.layout__menu-btn {
  display: none;
  position: fixed;
  top: 14px;
  left: 14px;
  z-index: 60;
  width: 42px;
  height: 42px;
  padding: 0;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  cursor: pointer;
  box-shadow: var(--shadow);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.layout__menu-btn span {
  display: block;
  width: 16px;
  height: 2px;
  border-radius: 2px;
  background: var(--text);
}

.layout__backdrop {
  display: none;
}

@media (max-width: 860px) {
  .layout__main {
    margin-left: 0;
    padding-top: 56px;
  }

  .layout__menu-btn {
    display: flex;
  }

  .layout__backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(18, 20, 48, 0.5);
    backdrop-filter: blur(2px);
    z-index: 40;
  }
}
</style>
