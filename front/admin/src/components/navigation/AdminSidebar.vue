<template>
  <aside class="admin-sidebar" :class="{ 'admin-sidebar--collapsed': collapsed }">
    <div class="admin-sidebar__glow" aria-hidden="true" />

    <div class="admin-sidebar__brand">
      <span class="admin-sidebar__logo">AI</span>
      <transition name="admin-sidebar-fade">
        <div v-if="!collapsed" class="admin-sidebar__brand-text">
          <strong>学习管理</strong>
          <span>Admin Console</span>
        </div>
      </transition>
    </div>

    <el-scrollbar class="admin-sidebar__scroll">
      <nav class="admin-sidebar__nav">
        <div v-for="group in visibleGroups" :key="group.key" class="admin-sidebar__group">
          <div v-if="!collapsed" class="admin-sidebar__group-title">{{ group.title }}</div>
          <div v-else class="admin-sidebar__group-rule" />

          <el-tooltip
            v-for="item in group.items"
            :key="item.path"
            :content="item.title"
            placement="right"
            :disabled="!collapsed"
            :show-after="200"
          >
            <router-link
              :to="item.path"
              class="admin-sidebar__link"
              :class="{ 'admin-sidebar__link--active': isActive(item) }"
            >
              <span class="admin-sidebar__link-icon">
                <el-icon :size="18"><component :is="item.icon" /></el-icon>
              </span>
              <span v-if="!collapsed" class="admin-sidebar__link-text">{{ item.title }}</span>
              <span v-if="!collapsed && isActive(item)" class="admin-sidebar__link-dot" />
            </router-link>
          </el-tooltip>
        </div>
      </nav>
    </el-scrollbar>

    <div class="admin-sidebar__foot">
      <button type="button" class="admin-sidebar__collapse" :title="collapsed ? '展开侧栏' : '收起侧栏'" @click="$emit('toggle')">
        <el-icon :size="16"><component :is="collapsed ? Expand : Fold" /></el-icon>
        <span v-if="!collapsed">收起侧栏</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { Expand, Fold } from '@element-plus/icons-vue'

const props = defineProps({
  collapsed: { type: Boolean, default: true },
  visibleGroups: { type: Array, required: true },
  activeMenu: { type: String, required: true }
})

defineEmits(['toggle'])

const route = useRoute()

function isActive(item) {
  if (item.activePrefix) return route.path.startsWith(item.activePrefix)
  return props.activeMenu === item.path
}
</script>

<style scoped>
.admin-sidebar {
  position: sticky;
  top: 0;
  z-index: 30;
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100vh;
  flex-shrink: 0;
  background: linear-gradient(180deg, var(--ink-2) 0%, var(--ink) 100%);
  color: var(--ink-text);
  overflow: hidden;
  transition: width var(--t);
}

.admin-sidebar--collapsed {
  width: var(--sidebar-collapsed);
}

.admin-sidebar__glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(360px 260px at 0% 0%, rgba(107, 92, 255, 0.35), transparent 70%),
    radial-gradient(300px 260px at 100% 100%, rgba(15, 185, 129, 0.16), transparent 70%);
}

.admin-sidebar__brand {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  height: var(--header-height);
  padding: 0 18px;
  flex-shrink: 0;
}

.admin-sidebar--collapsed .admin-sidebar__brand {
  justify-content: center;
  padding: 0;
}

.admin-sidebar__logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  border-radius: 13px;
  background: linear-gradient(135deg, var(--primary) 0%, #a396ff 100%);
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  font-size: 15px;
  box-shadow: 0 10px 24px rgba(107, 92, 255, 0.45), inset 0 1px 0 rgba(255, 255, 255, 0.35);
}

.admin-sidebar__brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
  white-space: nowrap;
}

.admin-sidebar__brand-text strong {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-size: 15px;
  color: #fff;
}

.admin-sidebar__brand-text span {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ink-text-2);
}

.admin-sidebar__scroll {
  position: relative;
  flex: 1;
}

.admin-sidebar__nav {
  padding: 8px 14px 16px;
}

.admin-sidebar--collapsed .admin-sidebar__nav {
  padding: 8px 14px 16px;
}

.admin-sidebar__group + .admin-sidebar__group {
  margin-top: 14px;
}

.admin-sidebar__group-title {
  padding: 0 10px 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-text-2);
  white-space: nowrap;
}

.admin-sidebar__group-rule {
  height: 1px;
  margin: 0 10px 10px;
  background: rgba(255, 255, 255, 0.08);
}

.admin-sidebar__link {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  height: 42px;
  padding: 0 10px;
  margin-bottom: 4px;
  border-radius: 12px;
  color: var(--ink-text);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  transition: background var(--t-fast), color var(--t-fast), transform var(--t-fast);
}

.admin-sidebar--collapsed .admin-sidebar__link {
  justify-content: center;
  padding: 0;
  width: 48px;
  height: 44px;
}

.admin-sidebar__link:hover {
  background: rgba(255, 255, 255, 0.07);
  color: #fff;
}

.admin-sidebar__link--active {
  background: linear-gradient(135deg, rgba(107, 92, 255, 0.95) 0%, rgba(134, 118, 255, 0.85) 100%);
  color: #fff;
  box-shadow: 0 10px 24px rgba(107, 92, 255, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.admin-sidebar__link-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 9px;
  background: rgba(255, 255, 255, 0.05);
  flex-shrink: 0;
  transition: background var(--t-fast);
}

.admin-sidebar__link--active .admin-sidebar__link-icon {
  background: rgba(255, 255, 255, 0.18);
}

.admin-sidebar__link-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}

.admin-sidebar__link-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.18);
}

.admin-sidebar__foot {
  position: relative;
  padding: 12px 14px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.07);
}

.admin-sidebar__collapse {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 38px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 11px;
  background: rgba(255, 255, 255, 0.04);
  color: var(--ink-text-2);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-sidebar__collapse:hover {
  background: rgba(255, 255, 255, 0.09);
  color: #fff;
}

.admin-sidebar-fade-enter-active,
.admin-sidebar-fade-leave-active {
  transition: opacity 0.15s ease;
}

.admin-sidebar-fade-enter-from,
.admin-sidebar-fade-leave-to {
  opacity: 0;
}
</style>
