<template>
  <el-dialog
    v-model="visible"
    title="快速导航"
    width="520px"
    class="admin-command-palette"
    :show-close="true"
    @closed="query = ''"
  >
    <el-input
      ref="inputRef"
      v-model="query"
      placeholder="搜索页面名称或关键词..."
      :prefix-icon="Search"
      clearable
      @keydown.down.prevent="moveHighlight(1)"
      @keydown.up.prevent="moveHighlight(-1)"
      @keydown.enter.prevent="confirmHighlight"
    />

    <div class="admin-command-palette__section">
      <div v-if="!query && recentItems.length" class="admin-command-palette__label">最近访问</div>
      <div v-else-if="query && !listItems.length" class="admin-command-palette__empty">未找到匹配的页面</div>
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
        <el-icon><component :is="item.icon" /></el-icon>
        <span class="admin-command-palette__title">{{ item.title }}</span>
        <span class="admin-command-palette__group">{{ item.groupTitle }}</span>
      </button>
    </div>

    <template #footer>
      <span class="admin-command-palette__hint">
        <kbd>↑</kbd><kbd>↓</kbd> 选择 · <kbd>Enter</kbd> 跳转 · <kbd>Esc</kbd> 关闭
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Search } from '@element-plus/icons-vue'

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

<style scoped>
.admin-command-palette__section {
  margin-top: 12px;
  max-height: 320px;
  overflow-y: auto;
}

.admin-command-palette__label {
  font-size: 12px;
  color: #9ca3af;
  margin-bottom: 6px;
  padding: 0 4px;
}

.admin-command-palette__item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border: none;
  border-radius: 8px;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: background 0.12s ease;
}

.admin-command-palette__item:hover,
.admin-command-palette__item--active {
  background: #f3f4f6;
}

.admin-command-palette__title {
  flex: 1;
  font-size: 14px;
  color: #1f2937;
}

.admin-command-palette__group {
  font-size: 12px;
  color: #9ca3af;
}

.admin-command-palette__empty {
  padding: 24px;
  text-align: center;
  color: #9ca3af;
  font-size: 14px;
}

.admin-command-palette__hint {
  font-size: 12px;
  color: #9ca3af;
}

.admin-command-palette__hint kbd {
  display: inline-block;
  padding: 2px 6px;
  margin: 0 2px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  background: #f9fafb;
  font-size: 11px;
  font-family: inherit;
}
</style>
