<template>
  <el-dialog
    v-model="visible"
    width="560px"
    class="admin-command-palette"
    :show-close="false"
    align-center
    @closed="query = ''"
  >
    <template #header>
      <div class="admin-command-palette__search">
        <el-icon :size="18" class="admin-command-palette__search-icon"><Search /></el-icon>
        <input
          ref="inputRef"
          v-model="query"
          class="admin-command-palette__input"
          placeholder="搜索页面名称或关键词…"
          @keydown.down.prevent="moveHighlight(1)"
          @keydown.up.prevent="moveHighlight(-1)"
          @keydown.enter.prevent="confirmHighlight"
        />
        <kbd class="admin-command-palette__esc">Esc</kbd>
      </div>
    </template>

    <div class="admin-command-palette__section">
      <div v-if="!query && recentItems.length" class="admin-command-palette__label">最近访问</div>
      <div v-else-if="query && !listItems.length" class="admin-command-palette__empty">
        <el-icon :size="28"><Search /></el-icon>
        <span>未找到匹配的页面</span>
      </div>
      <div v-else-if="!query" class="admin-command-palette__label">全部页面</div>

      <button
        v-for="(item, index) in listItems"
        :key="item.path + '-' + index"
        type="button"
        class="admin-command-palette__item"
        :class="{ 'admin-command-palette__item--active': highlightIndex === index }"
        @click="select(item.path)"
        @mouseenter="highlightIndex = index"
      >
        <span class="admin-command-palette__item-icon">
          <el-icon :size="17"><component :is="item.icon" /></el-icon>
        </span>
        <span class="admin-command-palette__title">{{ item.title }}</span>
        <span class="admin-command-palette__group">{{ item.groupTitle }}</span>
        <el-icon :size="14" class="admin-command-palette__enter"><Right /></el-icon>
      </button>
    </div>

    <template #footer>
      <span class="admin-command-palette__hint">
        <kbd>↑</kbd><kbd>↓</kbd> 选择 <span class="dot" /> <kbd>Enter</kbd> 跳转 <span class="dot" /> <kbd>Esc</kbd> 关闭
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Search, Right } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  searchNav: { type: Function, required: true },
  visibleItems: { type: Array, required: true },
  recentItems: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:modelValue', 'select'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const query = ref('')
const highlightIndex = ref(0)
const inputRef = ref(null)

const results = computed(() => props.searchNav(query.value))

const listItems = computed(() => {
  if (query.value) return results.value
  const recentPaths = new Set(props.recentItems.map((item) => item.path))
  const rest = props.visibleItems.filter((item) => !recentPaths.has(item.path))
  return [...props.recentItems, ...rest]
})

watch(visible, async (open) => {
  if (open) {
    highlightIndex.value = 0
    query.value = ''
    await nextTick()
    inputRef.value?.focus()
  }
})

watch(query, () => {
  highlightIndex.value = 0
})

function moveHighlight(delta) {
  const max = listItems.value.length - 1
  if (max < 0) return
  highlightIndex.value = Math.max(0, Math.min(max, highlightIndex.value + delta))
}

function confirmHighlight() {
  const item = listItems.value[highlightIndex.value]
  if (item) select(item.path)
}

function select(path) {
  emit('select', path)
  visible.value = false
}
</script>

<style>
.admin-command-palette .el-dialog__header {
  padding: 0;
  border-bottom: 1px solid var(--border);
}

.admin-command-palette .el-dialog__body {
  padding: 10px;
}

.admin-command-palette .el-dialog__footer {
  padding: 10px 18px;
}
</style>

<style scoped>
.admin-command-palette__search {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
}

.admin-command-palette__search-icon {
  color: var(--text-3);
}

.admin-command-palette__input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 16px;
  font-family: inherit;
  color: var(--text);
}

.admin-command-palette__input::placeholder {
  color: var(--text-3);
}

.admin-command-palette__esc,
.admin-command-palette__hint kbd {
  display: inline-block;
  padding: 2px 7px;
  border: 1px solid var(--border-strong);
  border-radius: 6px;
  background: var(--surface-2);
  font-size: 11px;
  font-family: inherit;
  color: var(--text-3);
}

.admin-command-palette__section {
  max-height: 360px;
  overflow-y: auto;
}

.admin-command-palette__label {
  font-size: 11.5px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-3);
  padding: 8px 12px 6px;
}

.admin-command-palette__item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 9px 12px;
  border: none;
  border-radius: 12px;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: background var(--t-fast);
}

.admin-command-palette__item-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: var(--surface-2);
  color: var(--text-3);
  flex-shrink: 0;
  transition: all var(--t-fast);
}

.admin-command-palette__item--active {
  background: var(--primary-soft);
}

.admin-command-palette__item--active .admin-command-palette__item-icon {
  background: var(--primary);
  color: #fff;
}

.admin-command-palette__title {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.admin-command-palette__group {
  font-size: 12px;
  color: var(--text-3);
}

.admin-command-palette__enter {
  color: var(--primary);
  opacity: 0;
  transition: opacity var(--t-fast);
}

.admin-command-palette__item--active .admin-command-palette__enter {
  opacity: 1;
}

.admin-command-palette__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 36px;
  text-align: center;
  color: var(--text-3);
  font-size: 14px;
}

.admin-command-palette__hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-3);
}

.admin-command-palette__hint .dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--border-strong);
  margin: 0 4px;
}
</style>
