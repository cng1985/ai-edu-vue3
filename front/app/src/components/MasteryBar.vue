<script setup>
import { computed } from 'vue'
import { masteryColor, pct } from '../utils/eco'

const props = defineProps({
  value: { type: Number, default: 0 },
  label: { type: String, default: '' },
  compact: { type: Boolean, default: false }
})

const percent = computed(() => pct(props.value))
</script>

<template>
  <div class="mastery" :class="{ 'mastery--compact': compact }">
    <div v-if="!compact" class="mastery__head">
      <span>{{ label || '掌握度' }}</span>
      <strong class="num">{{ percent }}%</strong>
    </div>
    <div class="progress progress--thin">
      <i :style="{ width: percent + '%', background: masteryColor(value) }"></i>
    </div>
    <strong v-if="compact" class="num mastery__value">{{ percent }}%</strong>
  </div>
</template>

<style scoped>
.mastery__head { display: flex; justify-content: space-between; margin-bottom: 6px; color: var(--text-2); font-size: 12.5px; }
.mastery--compact { display: flex; align-items: center; gap: 8px; min-width: 110px; }
.mastery--compact .progress { flex: 1; }
.mastery__value { width: 36px; color: var(--text-2); font-size: 12px; text-align: right; }
</style>
