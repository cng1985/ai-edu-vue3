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
          {{ group.title }}
        </button>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup>
const props = defineProps({
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
  border-bottom: 1px solid #e5e7eb;
  background: #fff;
}

.admin-topnav__tabs {
  display: flex;
  gap: 4px;
  padding: 0 16px;
  min-height: 40px;
  align-items: center;
}

.admin-topnav__tab {
  border: none;
  background: transparent;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
}

.admin-topnav__tab:hover {
  background: #f3f4f6;
  color: #374151;
}

.admin-topnav__tab--active {
  background: #eef2ff;
  color: #4f46e5;
  font-weight: 600;
}
</style>
