<template>
  <aside class="admin-win11-nav">
    <div class="admin-win11-nav__rail">
      <div class="admin-win11-nav__logo" title="AI 学习管理">
        <span>AI</span>
      </div>

      <nav class="admin-win11-nav__icons">
        <button
          v-for="group in visibleGroups"
          :key="group.key"
          type="button"
          class="admin-win11-nav__icon-btn"
          :class="{ 'admin-win11-nav__icon-btn--active': activeGroup === group.key || openGroup === group.key }"
          :title="group.title"
          @click="toggleGroup(group)"
        >
          <el-icon :size="20"><component :is="getGroupIcon(group.key)" /></el-icon>
        </button>
      </nav>
    </div>

    <Transition name="admin-win11-flyout">
      <div
        v-if="openGroupData"
        class="admin-win11-nav__flyout"
        @click.self="closeFlyout"
      >
        <div class="admin-win11-nav__panel">
          <div class="admin-win11-nav__panel-header">
            <h3>{{ openGroupData.title }}</h3>
            <el-button :icon="Close" text size="small" @click="closeFlyout" />
          </div>
          <div class="admin-win11-nav__panel-items">
            <router-link
              v-for="item in openGroupData.items"
              :key="item.path"
              :to="item.path"
              class="admin-win11-nav__item"
              :class="{ 'admin-win11-nav__item--active': isActive(item) }"
              @click="closeFlyout"
            >
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </router-link>
          </div>
        </div>
      </div>
    </Transition>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Close } from '@element-plus/icons-vue'
import { getGroupIcon } from '../../utils/navGroupIcons'

const props = defineProps({
  visibleGroups: { type: Array, required: true },
  activeMenu: { type: String, required: true },
  activeGroup: { type: String, required: true }
})

const route = useRoute()
const openGroup = ref(null)

const openGroupData = computed(() =>
  props.visibleGroups.find((g) => g.key === openGroup.value) || null
)

function toggleGroup(group) {
  openGroup.value = openGroup.value === group.key ? null : group.key
}

function closeFlyout() {
  openGroup.value = null
}

function isActive(item) {
  if (item.activePrefix) return route.path.startsWith(item.activePrefix)
  return props.activeMenu === item.path
}
</script>

<style scoped>
.admin-win11-nav {
  position: relative;
  flex-shrink: 0;
}

.admin-win11-nav__rail {
  width: 52px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: rgba(243, 243, 243, 0.85);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(0, 0, 0, 0.06);
}

.admin-win11-nav__logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  color: #4f46e5;
  letter-spacing: 0.02em;
}

.admin-win11-nav__icons {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 6px;
}

.admin-win11-nav__icon-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #4b5563;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.admin-win11-nav__icon-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.admin-win11-nav__icon-btn--active {
  background: rgba(79, 70, 229, 0.12);
  color: #4f46e5;
}

.admin-win11-nav__flyout {
  position: fixed;
  top: 0;
  left: 52px;
  bottom: 0;
  right: 0;
  z-index: 100;
}

.admin-win11-nav__panel {
  width: 280px;
  height: 100%;
  background: rgba(252, 252, 252, 0.92);
  backdrop-filter: blur(24px);
  border-right: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
}

.admin-win11-nav__panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 12px 12px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.admin-win11-nav__panel-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
}

.admin-win11-nav__panel-items {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.admin-win11-nav__item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 8px;
  color: #374151;
  text-decoration: none;
  font-size: 14px;
  transition: background 0.12s ease;
}

.admin-win11-nav__item:hover {
  background: rgba(0, 0, 0, 0.04);
}

.admin-win11-nav__item--active {
  background: rgba(79, 70, 229, 0.1);
  color: #4f46e5;
  font-weight: 600;
}

.admin-win11-flyout-enter-active,
.admin-win11-flyout-leave-active {
  transition: opacity 0.2s ease;
}

.admin-win11-flyout-enter-active .admin-win11-nav__panel,
.admin-win11-flyout-leave-active .admin-win11-nav__panel {
  transition: transform 0.2s ease;
}

.admin-win11-flyout-enter-from,
.admin-win11-flyout-leave-to {
  opacity: 0;
}

.admin-win11-flyout-enter-from .admin-win11-nav__panel,
.admin-win11-flyout-leave-to .admin-win11-nav__panel {
  transform: translateX(-12px);
}
</style>
