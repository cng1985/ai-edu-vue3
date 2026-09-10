<template>
  <el-dropdown trigger="click" @command="onSelect">
    <button type="button" class="admin-nav-mode-switcher" title="切换导航风格">
      <el-icon :size="16"><Brush /></el-icon>
      <span v-if="!compact" class="admin-nav-mode-switcher__label">{{ currentLabel }}</span>
    </button>
    <template #dropdown>
      <el-dropdown-menu class="admin-nav-mode-switcher__menu">
        <el-dropdown-item
          v-for="mode in NAV_MODES"
          :key="mode.key"
          :command="mode.key"
          :class="{ 'is-active': navMode === mode.key }"
        >
          <div class="admin-nav-mode-switcher__option">
            <span class="admin-nav-mode-switcher__option-icon">
              <el-icon :size="16"><component :is="mode.icon" /></el-icon>
            </span>
            <span class="admin-nav-mode-switcher__option-text">
              <span class="admin-nav-mode-switcher__option-title">{{ mode.label }}</span>
              <span class="admin-nav-mode-switcher__option-desc">{{ mode.description }}</span>
            </span>
            <el-icon v-if="navMode === mode.key" :size="15" class="admin-nav-mode-switcher__check"><Check /></el-icon>
          </div>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { Brush, Check } from '@element-plus/icons-vue'
import { NAV_MODES } from '../../config/navLayout'

const props = defineProps({
  navMode: { type: String, required: true },
  compact: { type: Boolean, default: false }
})

const emit = defineEmits(['change'])

const currentLabel = computed(() =>
  NAV_MODES.find((m) => m.key === props.navMode)?.label || '导航风格'
)

function onSelect(mode) {
  emit('change', mode)
}
</script>

<style scoped>
.admin-nav-mode-switcher {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  height: 38px;
  padding: 0 12px;
  border: 1px solid transparent;
  border-radius: 11px;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-nav-mode-switcher:hover {
  background: var(--surface-2);
  border-color: var(--border);
  color: var(--text);
}

.admin-nav-mode-switcher__option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 2px;
  min-width: 240px;
}

.admin-nav-mode-switcher__option-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: var(--surface-3);
  color: var(--text-3);
  flex-shrink: 0;
}

.admin-nav-mode-switcher__option-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  flex: 1;
  line-height: 1.3;
}

.admin-nav-mode-switcher__option-title {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text);
}

.admin-nav-mode-switcher__option-desc {
  font-size: 11.5px;
  color: var(--text-3);
  font-weight: 400;
}

.admin-nav-mode-switcher__check {
  color: var(--primary);
}

:deep(.el-dropdown-menu__item.is-active .admin-nav-mode-switcher__option-icon) {
  background: var(--primary);
  color: #fff;
}

:deep(.el-dropdown-menu__item.is-active .admin-nav-mode-switcher__option-title) {
  color: var(--primary-deep);
}
</style>
