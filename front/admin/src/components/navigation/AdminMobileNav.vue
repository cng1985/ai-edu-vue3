<template>
  <Teleport to="body">
    <Transition name="admin-drawer-fade">
      <div
        v-if="open"
        class="admin-mobile-nav__backdrop"
        @click="$emit('close')"
      />
    </Transition>

    <Transition name="admin-drawer-slide">
      <aside v-if="open" class="admin-mobile-nav">
        <div class="admin-mobile-nav__glow" aria-hidden="true" />
        <div class="admin-mobile-nav__header">
          <div class="admin-mobile-nav__brand">
            <span class="admin-mobile-nav__logo">AI</span>
            <strong>学习管理</strong>
          </div>
          <button type="button" class="admin-mobile-nav__close" @click="$emit('close')">
            <el-icon :size="18"><Close /></el-icon>
          </button>
        </div>
        <el-scrollbar class="admin-mobile-nav__body">
          <div
            v-for="group in visibleGroups"
            :key="group.key"
            class="admin-mobile-nav__group"
          >
            <div class="admin-mobile-nav__group-title">{{ group.title }}</div>
            <router-link
              v-for="item in group.items"
              :key="item.path"
              :to="item.path"
              class="admin-mobile-nav__link"
              :class="{ 'admin-mobile-nav__link--active': activeMenu === item.path || (item.activePrefix && $route.path.startsWith(item.activePrefix)) }"
              @click="$emit('close')"
            >
              <span class="admin-mobile-nav__link-icon">
                <el-icon :size="17"><component :is="item.icon" /></el-icon>
              </span>
              <span>{{ item.title }}</span>
            </router-link>
          </div>
        </el-scrollbar>
      </aside>
    </Transition>
  </Teleport>
</template>

<script setup>
import { Close } from '@element-plus/icons-vue'

defineProps({
  open: { type: Boolean, default: false },
  visibleGroups: { type: Array, required: true },
  activeMenu: { type: String, required: true }
})

defineEmits(['close'])
</script>

<style scoped>
.admin-mobile-nav__backdrop {
  position: fixed;
  inset: 0;
  background: rgba(18, 20, 48, 0.5);
  backdrop-filter: blur(4px);
  z-index: 2000;
}

.admin-mobile-nav {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: min(300px, 86vw);
  background: linear-gradient(180deg, var(--ink-2) 0%, var(--ink) 100%);
  color: var(--ink-text);
  z-index: 2001;
  display: flex;
  flex-direction: column;
  box-shadow: 12px 0 40px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.admin-mobile-nav__glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(360px 260px at 0% 0%, rgba(107, 92, 255, 0.35), transparent 70%);
}

.admin-mobile-nav__header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 12px 14px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-mobile-nav__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
}

.admin-mobile-nav__logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 11px;
  background: linear-gradient(135deg, var(--primary) 0%, #a396ff 100%);
  font-weight: 800;
  font-size: 13px;
}

.admin-mobile-nav__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.06);
  color: var(--ink-text);
  cursor: pointer;
}

.admin-mobile-nav__body {
  position: relative;
  flex: 1;
  padding: 8px 0 16px;
}

.admin-mobile-nav__group {
  padding: 8px 14px;
}

.admin-mobile-nav__group-title {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--ink-text-2);
  padding: 4px 10px 8px;
}

.admin-mobile-nav__link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  margin-bottom: 4px;
  border-radius: 12px;
  color: var(--ink-text);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: background var(--t-fast);
}

.admin-mobile-nav__link-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 9px;
  background: rgba(255, 255, 255, 0.05);
}

.admin-mobile-nav__link:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}

.admin-mobile-nav__link--active {
  background: linear-gradient(135deg, rgba(107, 92, 255, 0.95) 0%, rgba(134, 118, 255, 0.85) 100%);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 10px 24px rgba(107, 92, 255, 0.35);
}

.admin-mobile-nav__link--active .admin-mobile-nav__link-icon {
  background: rgba(255, 255, 255, 0.18);
}

.admin-drawer-fade-enter-active,
.admin-drawer-fade-leave-active {
  transition: opacity 0.2s ease;
}

.admin-drawer-fade-enter-from,
.admin-drawer-fade-leave-to {
  opacity: 0;
}

.admin-drawer-slide-enter-active,
.admin-drawer-slide-leave-active {
  transition: transform 0.28s var(--ease);
}

.admin-drawer-slide-enter-from,
.admin-drawer-slide-leave-to {
  transform: translateX(-100%);
}
</style>
