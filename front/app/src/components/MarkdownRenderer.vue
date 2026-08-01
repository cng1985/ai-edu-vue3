<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { renderMarkdown, renderMermaidIn } from '../utils/markdown'

const props = defineProps({
  source: { type: String, default: '' },
  // 流式输出场景下关闭 mermaid 渲染，避免对半截代码块反复渲染
  live: { type: Boolean, default: false }
})

const root = ref(null)
const html = computed(() => renderMarkdown(props.source))

async function paintMermaid() {
  if (props.live || !root.value) return
  await nextTick()
  renderMermaidIn(root.value)
}

onMounted(paintMermaid)
watch(() => [props.source, props.live], paintMermaid)
</script>

<template>
  <div ref="root" class="markdown-body" v-html="html"></div>
</template>
