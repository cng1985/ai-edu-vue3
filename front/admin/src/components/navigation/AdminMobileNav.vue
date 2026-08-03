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
        <div class="admin-mobile-nav__header">
          <span>导航菜单</span>
          <el-button :icon="Close" text @click="$emit('close')" />
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
              <el-icon><component :is="item.icon" /></el-icon>
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
  background: rgba(15, 23, 42, 0.45);
  z-index: 2000;
}

.admin-mobile-nav {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: min(280px, 85vw);
  background: #1e1b4b;
  z-index: 2001;
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.2);
}

.admin-mobile-nav__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 12px 12px 16px;
  color: #fff;
  font-weight: 600;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-mobile-nav__body {
  flex: 1;
  padding: 8px 0 16px;
}

.admin-mobile-nav__group {
  padding: 8px 12px;
}

.admin-mobile-nav__group-title {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: #a5b4fc;
  padding: 4px 8px 8px;
}

.admin-mobile-nav__link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #c7d2fe;
  text-decoration: none;
  font-size: 14px;
  transition: background 0.12s ease;
}

.admin-mobile-nav__link:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

.admin-mobile-nav__link--active {
  background: rgba(99, 102, 241, 0.35);
  color: #fff;
  font-weight: 600;
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
  transition: transform 0.25s ease;
}

.admin-drawer-slide-enter-from,
.admin-drawer-slide-leave-to {
  transform: translateX(-100%);
}
</style>
