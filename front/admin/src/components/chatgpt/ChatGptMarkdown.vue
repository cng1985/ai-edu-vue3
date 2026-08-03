<template>
  <div class="cgpt-md" ref="rootRef">
    <div v-html="html"></div>
    <span v-if="live" class="cgpt-md__cursor" />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { marked } from 'marked'

const props = defineProps({
  source: { type: String, default: '' },
  live: { type: Boolean, default: false }
})

const rootRef = ref(null)

const html = computed(() => {
  if (!props.source) return ''
  return marked.parse(props.source, { breaks: true, gfm: true })
})

function bindCopyButtons() {
  const root = rootRef.value
  if (!root) return
  root.querySelectorAll('pre').forEach((pre) => {
    if (pre.querySelector('.cgpt-code-copy')) return
    const btn = document.createElement('button')
    btn.type = 'button'
    btn.className = 'cgpt-code-copy'
    btn.textContent = '复制代码'
    btn.addEventListener('click', async () => {
      const code = pre.querySelector('code')
      const text = code?.textContent || pre.textContent || ''
      try {
        await navigator.clipboard.writeText(text)
        btn.textContent = '已复制'
        setTimeout(() => { btn.textContent = '复制代码' }, 1500)
      } catch {
        btn.textContent = '复制失败'
      }
    })
    pre.style.position = 'relative'
    pre.appendChild(btn)
  })
}

watch(() => props.source, () => {
  requestAnimationFrame(bindCopyButtons)
}, { flush: 'post' })

onMounted(bindCopyButtons)
</script>

<style scoped>
.cgpt-md {
  font-size: 16px;
  line-height: 1.65;
  color: #0d0d0d;
  word-break: break-word;
}

.cgpt-md :deep(p) {
  margin: 0 0 12px;
}

.cgpt-md :deep(p:last-child) {
  margin-bottom: 0;
}

.cgpt-md :deep(ul),
.cgpt-md :deep(ol) {
  margin: 8px 0 12px;
  padding-left: 24px;
}

.cgpt-md :deep(li) {
  margin: 4px 0;
}

.cgpt-md :deep(h1),
.cgpt-md :deep(h2),
.cgpt-md :deep(h3) {
  margin: 16px 0 8px;
  font-weight: 600;
}

.cgpt-md :deep(blockquote) {
  margin: 12px 0;
  padding-left: 12px;
  border-left: 3px solid #d1d5db;
  color: #4b5563;
}

.cgpt-md :deep(a) {
  color: #10a37f;
}

.cgpt-md :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.9em;
  background: #f4f4f4;
  padding: 2px 6px;
  border-radius: 6px;
}

.cgpt-md :deep(pre) {
  margin: 12px 0;
  padding: 16px;
  background: #0d0d0d;
  color: #f5f5f5;
  border-radius: 12px;
  overflow-x: auto;
  font-size: 14px;
  line-height: 1.5;
}

.cgpt-md :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
}

.cgpt-md :deep(.cgpt-code-copy) {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 10px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  font-size: 12px;
  cursor: pointer;
}

.cgpt-md :deep(.cgpt-code-copy:hover) {
  background: rgba(255, 255, 255, 0.15);
}

.cgpt-md__cursor {
  display: inline-block;
  width: 8px;
  height: 18px;
  margin-left: 2px;
  vertical-align: text-bottom;
  background: #0d0d0d;
  animation: cgpt-cursor 1s step-end infinite;
}

@keyframes cgpt-cursor {
  50% { opacity: 0; }
}
</style>
