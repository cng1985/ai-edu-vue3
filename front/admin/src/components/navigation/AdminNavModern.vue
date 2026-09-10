<template>
  <nav class="admin-modern-nav">
    <el-scrollbar>
      <div class="admin-modern-nav__inner">
        <el-dropdown
          v-for="group in visibleGroups"
          :key="group.key"
          trigger="hover"
          :show-timeout="80"
          :hide-timeout="120"
          placement="bottom-start"
          @command="onNavigate"
        >
          <button
            type="button"
            class="admin-modern-nav__group"
            :class="{ 'admin-modern-nav__group--active': activeGroup === group.key }"
          >
            <el-icon :size="16"><component :is="getGroupIcon(group.key)" /></el-icon>
            <span>{{ group.title }}</span>
            <el-icon class="admin-modern-nav__arrow" :size="12"><ArrowDown /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                v-for="item in group.items"
                :key="item.path"
                :command="item.path"
                :class="{ 'is-active': isActive(item) }"
              >
                <el-icon><component :is="item.icon" /></el-icon>
                {{ item.title }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-scrollbar>
  </nav>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { getGroupIcon } from '../../utils/navGroupIcons'

defineProps({
  visibleGroups: { type: Array, required: true },
  activeMenu: { type: String, required: true },
  activeGroup: { type: String, required: true }
})

const route = useRoute()
const router = useRouter()

function onNavigate(path) {
  router.push(path)
}

function isActive(item) {
  if (item.activePrefix) return route.path.startsWith(item.activePrefix)
  return route.path === item.path
}
</script>

<style scoped>
.admin-modern-nav {
  background: rgba(255, 255, 255, 0.6);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(12px);
}

.admin-modern-nav__inner {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 24px;
}

.admin-modern-nav__group {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 8px 14px;
  border: none;
  border-radius: 11px;
  background: transparent;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-3);
  cursor: pointer;
  white-space: nowrap;
  transition: all var(--t-fast);
}

.admin-modern-nav__group:hover {
  background: var(--surface-3);
  color: var(--text);
}

.admin-modern-nav__group--active {
  background: var(--primary-soft);
  color: var(--primary-deep);
}

.admin-modern-nav__arrow {
  opacity: 0.5;
  margin-left: 2px;
}

:deep(.el-dropdown-menu__item.is-active) {
  color: var(--primary-deep);
  font-weight: 700;
  background: var(--primary-soft);
}

:deep(.el-dropdown-menu__item .el-icon) {
  margin-right: 8px;
}
</style>
