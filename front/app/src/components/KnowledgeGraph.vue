<script setup>
import { computed, ref } from 'vue'
import { masteryColor, pct } from '../utils/eco'

/**
 * 分层知识图谱：按前置关系深度分列，前置/高级关系画实线箭头，关联/相似关系画虚线。
 */
const props = defineProps({
  nodes: { type: Array, default: () => [] },
  edges: { type: Array, default: () => [] },
  highlight: { type: Array, default: () => [] }
})
const emit = defineEmits(['select'])

const COL_W = 210
const ROW_H = 74
const NODE_W = 168
const NODE_H = 50
const PAD = 24

const hovered = ref('')

const layout = computed(() => {
  const cols = new Map()
  const sorted = [...props.nodes].sort((a, b) => a.domain.localeCompare(b.domain) || a.difficulty - b.difficulty)
  for (const n of sorted) {
    if (!cols.has(n.depth)) cols.set(n.depth, [])
    cols.get(n.depth).push(n)
  }
  const pos = {}
  let maxRows = 0
  for (const [depth, list] of cols) {
    maxRows = Math.max(maxRows, list.length)
    list.forEach((n, i) => {
      pos[n.id] = { x: PAD + depth * COL_W, y: PAD + i * ROW_H, node: n }
    })
  }
  const depthCount = cols.size ? Math.max(...cols.keys()) + 1 : 1
  return { pos, width: PAD * 2 + (depthCount - 1) * COL_W + NODE_W, height: PAD * 2 + Math.max(1, maxRows - 1) * ROW_H + NODE_H }
})

const focus = computed(() => {
  const id = hovered.value
  if (!id) return null
  const set = new Set([id])
  for (const e of props.edges) {
    if (e.fromId === id) set.add(e.toId)
    if (e.toId === id) set.add(e.fromId)
  }
  return set
})

const lines = computed(() =>
  props.edges
    .filter((e) => layout.value.pos[e.fromId] && layout.value.pos[e.toId])
    .map((e) => {
      const a = layout.value.pos[e.fromId]
      const b = layout.value.pos[e.toId]
      const directed = e.type === 'prerequisite' || e.type === 'advanced'
      const forward = b.x >= a.x
      const x1 = a.x + (forward ? NODE_W : 0)
      const x2 = b.x + (forward ? 0 : NODE_W)
      const y1 = a.y + NODE_H / 2
      const y2 = b.y + NODE_H / 2
      const dx = Math.max(40, Math.abs(x2 - x1) / 2)
      const path = a.x === b.x
        ? `M${a.x + NODE_W} ${y1} C${a.x + NODE_W + 40} ${y1} ${b.x + NODE_W + 40} ${y2} ${b.x + NODE_W} ${y2}`
        : `M${x1} ${y1} C${x1 + (forward ? dx : -dx)} ${y1} ${x2 - (forward ? dx : -dx)} ${y2} ${x2} ${y2}`
      const active = !focus.value || (focus.value.has(e.fromId) && focus.value.has(e.toId))
      return { id: e.id, path, directed, type: e.type, active }
    })
)

const highlightSet = computed(() => new Set(props.highlight))
</script>

<template>
  <div class="kg">
    <svg :width="layout.width" :height="layout.height" :viewBox="`0 0 ${layout.width} ${layout.height}`">
      <defs>
        <marker id="kg-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="#9aa0bf" />
        </marker>
      </defs>
      <path
        v-for="l in lines"
        :key="l.id"
        :d="l.path"
        class="kg__edge"
        :class="[`kg__edge--${l.type}`, { 'kg__edge--dim': !l.active }]"
        :marker-end="l.directed ? 'url(#kg-arrow)' : null"
      />
      <g
        v-for="(p, id) in layout.pos"
        :key="id"
        class="kg__node"
        :class="{ 'kg__node--dim': focus && !focus.has(id), 'kg__node--hl': highlightSet.has(id) }"
        :transform="`translate(${p.x}, ${p.y})`"
        @mouseenter="hovered = id"
        @mouseleave="hovered = ''"
        @click="emit('select', p.node)"
      >
        <rect :width="NODE_W" :height="NODE_H" rx="12" class="kg__box" :style="{ stroke: masteryColor(p.node.mastery) }" />
        <rect :width="Math.max(4, NODE_W * p.node.mastery)" height="4" :y="NODE_H - 4" rx="2" :style="{ fill: masteryColor(p.node.mastery) }" />
        <text x="12" y="21" class="kg__name">{{ p.node.name.length > 11 ? p.node.name.slice(0, 11) + '…' : p.node.name }}</text>
        <text x="12" y="38" class="kg__meta">{{ p.node.domain }} · {{ pct(p.node.mastery) }}%</text>
        <circle v-if="p.node.status === 'mastered'" :cx="NODE_W - 14" cy="14" r="5" fill="#0fb981" />
      </g>
    </svg>
  </div>
</template>

<style scoped>
.kg { overflow: auto; padding-bottom: 6px; }
.kg svg { display: block; }
.kg__edge { fill: none; stroke: #b6bad0; stroke-width: 1.6; transition: opacity var(--t-fast); }
.kg__edge--related, .kg__edge--similar { stroke-dasharray: 5 5; stroke: #cfd3e3; }
.kg__edge--advanced { stroke: #a99cff; }
.kg__edge--dim { opacity: 0.15; }
.kg__node { cursor: pointer; transition: opacity var(--t-fast); }
.kg__node--dim { opacity: 0.25; }
.kg__box { fill: #fff; stroke-width: 1.6; filter: drop-shadow(0 2px 4px rgba(22, 24, 44, 0.06)); }
.kg__node:hover .kg__box { fill: var(--primary-soft); }
.kg__node--hl .kg__box { fill: #fff8e8; stroke-width: 2.4; }
.kg__name { fill: var(--text); font-size: 12.5px; font-weight: 700; }
.kg__meta { fill: var(--text-3); font-size: 10.5px; }
</style>
