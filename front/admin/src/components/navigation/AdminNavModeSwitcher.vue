<template>
  <el-dropdown trigger="click" @command="onSelect">
    <el-button :icon="Brush" text class="admin-nav-mode-switcher">
      <span v-if="!compact" class="admin-nav-mode-switcher__label">{{ currentLabel }}</span>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="mode in NAV_MODES"
          :key="mode.key"
          :command="mode.key"
          :class="{ 'is-active': navMode === mode.key }"
        >
          <div class="admin-nav-mode-switcher__option">
            <span class="admin-nav-mode-switcher__option-title">{{ mode.label }}</span>
            <span class="admin-nav-mode-switcher__option-desc">{{ mode.description }}</span>
          </div>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { Brush } from '@element-plus/icons-vue'
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
  color: #6b7280;
}

.admin-nav-mode-switcher__label {
  margin-left: 4px;
  font-size: 13px;
}

.admin-nav-mode-switcher__option {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 2px 0;
}

.admin-nav-mode-switcher__option-title {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.admin-nav-mode-switcher__option-desc {
  font-size: 11px;
  color: #9ca3af;
}

:deep(.el-dropdown-menu__item.is-active .admin-nav-mode-switcher__option-title) {
  color: #4f46e5;
  font-weight: 600;
}
</style>
