<template>
  <el-aside :width="collapsed ? '64px' : '220px'" class="admin-sidebar">
    <div class="admin-sidebar__logo">
      <span v-if="!collapsed">AI 学习管理</span>
      <span v-else>AI</span>
    </div>

    <el-scrollbar class="admin-sidebar__scroll">
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :default-openeds="openGroups"
        router
        class="admin-sidebar__menu"
        background-color="#1e1b4b"
        text-color="#c7d2fe"
        active-text-color="#ffffff"
      >
        <el-sub-menu
          v-for="group in visibleGroups"
          :key="group.key"
          :index="group.key"
        >
          <template #title>
            <el-icon><component :is="getGroupIcon(group.key)" /></el-icon>
            <span>{{ group.title }}</span>
          </template>
          <el-menu-item
            v-for="item in group.items"
            :key="item.path"
            :index="item.path"
          >
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-scrollbar>
  </el-aside>
</template>

<script setup>
import { computed } from 'vue'
import { GROUP_ICONS, getGroupIcon } from '../../utils/navGroupIcons'

defineProps({
  collapsed: { type: Boolean, default: true },
  visibleGroups: { type: Array, required: true },
  activeMenu: { type: String, required: true }
})

const openGroups = computed(() => Object.keys(GROUP_ICONS))
</script>

<style scoped>
.admin-sidebar {
  background: #1e1b4b;
  transition: width 0.2s ease;
  overflow: hidden;
}

.admin-sidebar__logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.02em;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-sidebar__scroll {
  height: calc(100vh - 56px);
}

.admin-sidebar__menu {
  border-right: none;
}

.admin-sidebar__menu:not(.el-menu--collapse) {
  width: 220px;
}
</style>
