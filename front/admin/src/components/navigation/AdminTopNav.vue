<template>
  <div class="admin-topnav">
    <el-scrollbar>
      <div class="admin-topnav__tabs">
        <button
          v-for="group in visibleGroups"
          :key="group.key"
          type="button"
          class="admin-topnav__tab"
          :class="{ 'admin-topnav__tab--active': activeGroup === group.key }"
          @click="onGroupClick(group)"
        >
          <el-icon :size="15"><component :is="getGroupIcon(group.key)" /></el-icon>
          <span>{{ group.title }}</span>
        </button>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup>
import { getGroupIcon } from '../../utils/navGroupIcons'

defineProps({
  visibleGroups: { type: Array, required: true },
  activeGroup: { type: String, required: true }
})

const emit = defineEmits(['navigate'])

function onGroupClick(group) {
  const first = group.items[0]
  if (first) emit('navigate', first.path)
}
</script>

<style scoped>
.admin-topnav {
  background: rgba(255, 255, 255, 0.6);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(12px);
}

.admin-topnav__tabs {
  display: flex;
  gap: 4px;
  padding: 8px 24px;
  align-items: center;
}

.admin-topnav__tab {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  border: none;
  background: transparent;
  padding: 7px 13px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-3);
  cursor: pointer;
  white-space: nowrap;
  transition: all var(--t-fast);
}

.admin-topnav__tab:hover {
  background: var(--surface-3);
  color: var(--text);
}

.admin-topnav__tab--active {
  background: var(--ink);
  color: #fff;
  box-shadow: 0 6px 14px rgba(18, 20, 48, 0.2);
}
</style>
