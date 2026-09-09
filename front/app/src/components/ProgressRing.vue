<script setup>
import { computed } from 'vue'

const props = defineProps({
  percent: { type: Number, default: 0 },
  size: { type: Number, default: 64 },
  stroke: { type: Number, default: 6 },
  color: { type: String, default: 'var(--primary)' },
  track: { type: String, default: 'var(--surface-3)' },
  showText: { type: Boolean, default: true }
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
      :stroke="track"
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
      v-if="showText"
      :x="size / 2"
      :y="size / 2"
      text-anchor="middle"
      dominant-baseline="central"
      class="progress-ring__text"
      :style="{ fontSize: Math.max(12, size * 0.2) + 'px' }"
    >
      {{ Math.round(percent) }}%
    </text>
  </svg>
</template>

<style scoped>
.progress-ring {
  display: block;
  flex-shrink: 0;
}

.progress-ring__value {
  transition: stroke-dashoffset 0.6s var(--ease);
}

.progress-ring__text {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-weight: 700;
  letter-spacing: -0.02em;
  fill: var(--text);
}
</style>
