<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  size: { type: Number, default: 260 }
})

const center = computed(() => props.size / 2)
const radius = computed(() => props.size * 0.32)
const count = computed(() => Math.max(3, props.items.length))

function point(index, ratio = 1) {
  const angle = -Math.PI / 2 + index * Math.PI * 2 / count.value
  return [
    center.value + Math.cos(angle) * radius.value * ratio,
    center.value + Math.sin(angle) * radius.value * ratio
  ]
}

function pointsFor(ratio) {
  return Array.from({ length: count.value }, (_, index) => point(index, ratio).join(',')).join(' ')
}

const valuePoints = computed(() => props.items
  .map((item, index) => point(index, Math.max(0, Math.min(100, item.progress)) / 100).join(','))
  .join(' '))

function labelPosition(index) {
  const [x, y] = point(index, 1.32)
  return { x, y }
}
</script>

<template>
  <svg :width="size" :height="size" :viewBox="`0 0 ${size} ${size}`" class="radar" role="img" aria-label="能力雷达图">
    <polygon v-for="level in 4" :key="level" :points="pointsFor(level / 4)" class="radar__grid" />
    <line
      v-for="(_, index) in items"
      :key="index"
      :x1="center"
      :y1="center"
      :x2="point(index)[0]"
      :y2="point(index)[1]"
      class="radar__axis"
    />
    <polygon :points="pointsFor(1)" class="radar__target" />
    <polygon v-if="items.length" :points="valuePoints" class="radar__value" />
    <g v-for="(item, index) in items" :key="item.id">
      <circle :cx="point(index, item.progress / 100)[0]" :cy="point(index, item.progress / 100)[1]" r="4" class="radar__dot" />
      <text
        :x="labelPosition(index).x"
        :y="labelPosition(index).y"
        text-anchor="middle"
        dominant-baseline="central"
        class="radar__label"
      >
        {{ item.name }}
      </text>
      <text
        :x="labelPosition(index).x"
        :y="labelPosition(index).y + 15"
        text-anchor="middle"
        dominant-baseline="central"
        class="radar__score"
      >
        {{ item.progress }}%
      </text>
    </g>
  </svg>
</template>

<style scoped>
.radar { display: block; max-width: 100%; margin: auto; overflow: visible; }
.radar__grid, .radar__axis { fill: #f8fafc; stroke: #cbd5e1; stroke-width: 1; }
.radar__axis { fill: none; }
.radar__target { fill: rgba(99, 102, 241, .04); stroke: #a5b4fc; stroke-width: 1.5; stroke-dasharray: 4 4; }
.radar__value { fill: rgba(99, 102, 241, .22); stroke: var(--primary); stroke-width: 2.5; }
.radar__dot { fill: var(--primary); stroke: white; stroke-width: 2; }
.radar__label { fill: var(--text-2); font-size: 11px; font-weight: 600; }
.radar__score { fill: var(--primary); font-size: 10px; font-weight: 700; }
</style>
