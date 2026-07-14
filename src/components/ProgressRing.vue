<script setup>
import { computed } from 'vue'

const props = defineProps({
  percent: { type: Number, default: 0 },
  size: { type: Number, default: 64 },
  stroke: { type: Number, default: 6 },
  color: { type: String, default: 'var(--primary)' }
})

const radius = computed(() => (props.size - props.stroke) / 2)
const circumference = computed(() => 2 * Math.PI * radius.value)
const offset = computed(
  () => circumference.value * (1 - Math.min(100, Math.max(0, props.percent)) / 100)
)
</script>

<template>
  <svg :width="size" :height="size" class="progress-ring">
    <circle
      :cx="size / 2"
      :cy="size / 2"
      :r="radius"
      fill="none"
      stroke="var(--border)"
      :stroke-width="stroke"
    />
    <circle
      :cx="size / 2"
      :cy="size / 2"
      :r="radius"
      fill="none"
      :stroke="color"
      :stroke-width="stroke"
      stroke-linecap="round"
      :stroke-dasharray="circumference"
      :stroke-dashoffset="offset"
      :transform="`rotate(-90 ${size / 2} ${size / 2})`"
      class="progress-ring__value"
    />
    <text
      :x="size / 2"
      :y="size / 2"
      text-anchor="middle"
      dominant-baseline="central"
      class="progress-ring__text"
    >
      {{ Math.round(percent) }}%
    </text>
  </svg>
</template>

<style scoped>
.progress-ring__value {
  transition: stroke-dashoffset 0.5s ease;
}

.progress-ring__text {
  font-size: 13px;
  font-weight: 700;
  fill: var(--text);
}
</style>
