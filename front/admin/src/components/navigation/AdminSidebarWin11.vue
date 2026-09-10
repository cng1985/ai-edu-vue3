<template>
  <aside class="admin-win11-nav">
    <div class="admin-win11-nav__rail">
      <div class="admin-win11-nav__logo" title="AI 学习管理">
        <span>AI</span>
      </div>

      <nav class="admin-win11-nav__icons">
        <el-tooltip
          v-for="group in visibleGroups"
          :key="group.key"
          :content="group.title"
          placement="right"
          :show-after="300"
          :disabled="openGroup === group.key"
        >
          <button
            type="button"
            class="admin-win11-nav__icon-btn"
            :class="{
              'admin-win11-nav__icon-btn--active': activeGroup === group.key,
              'admin-win11-nav__icon-btn--open': openGroup === group.key
            }"
            @click="toggleGroup(group)"
          >
            <el-icon :size="20"><component :is="getGroupIcon(group.key)" /></el-icon>
            <span v-if="activeGroup === group.key" class="admin-win11-nav__indicator" />
          </button>
        </el-tooltip>
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
            <div class="admin-win11-nav__panel-title">
              <span class="admin-win11-nav__panel-icon">
                <el-icon :size="18"><component :is="getGroupIcon(openGroupData.key)" /></el-icon>
              </span>
              <h3>{{ openGroupData.title }}</h3>
            </div>
            <button type="button" class="admin-win11-nav__close" @click="closeFlyout">
              <el-icon :size="16"><Close /></el-icon>
            </button>
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
              <span class="admin-win11-nav__item-icon">
                <el-icon :size="17"><component :is="item.icon" /></el-icon>
              </span>
              <span class="admin-win11-nav__item-text">{{ item.title }}</span>
              <el-icon :size="14" class="admin-win11-nav__item-arrow"><ArrowRight /></el-icon>
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
import { Close, ArrowRight } from '@element-plus/icons-vue'
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
  position: sticky;
  top: 0;
  z-index: 30;
  height: 100vh;
  flex-shrink: 0;
}

.admin-win11-nav__rail {
  width: 68px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(20px) saturate(1.4);
  -webkit-backdrop-filter: blur(20px) saturate(1.4);
  border-right: 1px solid var(--border);
}

.admin-win11-nav__logo {
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: center;
}

.admin-win11-nav__logo span {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary) 0%, #a396ff 100%);
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  font-size: 14px;
  box-shadow: 0 8px 18px var(--primary-glow);
}

.admin-win11-nav__icons {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 10px 0;
}

.admin-win11-nav__icon-btn {
  position: relative;
  width: 46px;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 14px;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-win11-nav__icon-btn:hover {
  background: var(--surface-3);
  color: var(--text);
}

.admin-win11-nav__icon-btn--active {
  color: var(--primary);
  background: var(--primary-soft);
}

.admin-win11-nav__icon-btn--open {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  color: #fff;
  box-shadow: 0 8px 18px var(--primary-glow);
}

.admin-win11-nav__indicator {
  position: absolute;
  left: -11px;
  top: 50%;
  width: 4px;
  height: 18px;
  border-radius: 999px;
  background: var(--primary);
  transform: translateY(-50%);
}

.admin-win11-nav__icon-btn--open .admin-win11-nav__indicator {
  opacity: 0;
}

.admin-win11-nav__flyout {
  position: fixed;
  top: 0;
  left: 68px;
  bottom: 0;
  right: 0;
  z-index: 100;
}

.admin-win11-nav__panel {
  width: 300px;
  height: 100%;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(24px) saturate(1.4);
  -webkit-backdrop-filter: blur(24px) saturate(1.4);
  border-right: 1px solid var(--border);
  box-shadow: 12px 0 40px rgba(22, 24, 44, 0.08);
  display: flex;
  flex-direction: column;
}

.admin-win11-nav__panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--header-height);
  padding: 0 14px 0 20px;
  border-bottom: 1px solid var(--border);
}

.admin-win11-nav__panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-win11-nav__panel-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: var(--primary-soft);
  color: var(--primary);
}

.admin-win11-nav__panel-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
}

.admin-win11-nav__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
}

.admin-win11-nav__close:hover {
  background: var(--surface-2);
  color: var(--text);
}

.admin-win11-nav__panel-items {
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.admin-win11-nav__item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 12px;
  color: var(--text-2);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all var(--t-fast);
}

.admin-win11-nav__item-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: var(--surface-2);
  color: var(--text-3);
  transition: all var(--t-fast);
}

.admin-win11-nav__item-text {
  flex: 1;
}

.admin-win11-nav__item-arrow {
  color: var(--border-strong);
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--t-fast);
}

.admin-win11-nav__item:hover {
  background: var(--surface-2);
  color: var(--text);
}

.admin-win11-nav__item:hover .admin-win11-nav__item-arrow {
  opacity: 1;
  transform: translateX(0);
}

.admin-win11-nav__item--active {
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-weight: 700;
}

.admin-win11-nav__item--active .admin-win11-nav__item-icon {
  background: var(--primary);
  color: #fff;
  box-shadow: 0 6px 14px var(--primary-glow);
}

.admin-win11-flyout-enter-active,
.admin-win11-flyout-leave-active {
  transition: opacity 0.2s ease;
}

.admin-win11-flyout-enter-active .admin-win11-nav__panel,
.admin-win11-flyout-leave-active .admin-win11-nav__panel {
  transition: transform 0.22s var(--ease);
}

.admin-win11-flyout-enter-from,
.admin-win11-flyout-leave-to {
  opacity: 0;
}

.admin-win11-flyout-enter-from .admin-win11-nav__panel,
.admin-win11-flyout-leave-to .admin-win11-nav__panel {
  transform: translateX(-16px);
}
</style>
