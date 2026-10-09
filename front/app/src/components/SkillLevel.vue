<script setup>
import { computed } from 'vue'
import { levelInfo } from '../utils/eco'

const props = defineProps({
  level: { type: Number, default: 0 },
  required: { type: Number, default: null },
  showName: { type: Boolean, default: false }
})

const info = computed(() => levelInfo(props.level))
const met = computed(() => props.required == null || props.level >= props.required)
</script>

<template>
  <span class="skill-level" :title="info.name">
    <b :style="{ background: info.color }">L{{ level }}</b>
    <span class="skill-level__bar">
      <i
        v-for="n in 5"
        :key="n"
        :class="{ on: n <= level, req: required != null && n <= required && n > level }"
        :style="n <= level ? { background: info.color } : null"
      ></i>
    </span>
    <em v-if="required != null" :class="met ? 'ok' : 'gap'">{{ met ? '达标' : `需 L${required}` }}</em>
    <small v-if="showName">{{ info.name }}</small>
  </span>
</template>

<style scoped>
.skill-level { display: inline-flex; align-items: center; gap: 7px; white-space: nowrap; }
.skill-level b {
  min-width: 26px;
  padding: 1px 6px;
  border-radius: 7px;
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-size: 11px;
  text-align: center;
}
.skill-level__bar { display: inline-flex; gap: 3px; }
.skill-level__bar i { width: 12px; height: 6px; border-radius: 3px; background: var(--surface-3); }
.skill-level__bar i.req { background: repeating-linear-gradient(135deg, #ffd9a8 0 2px, #fff5e0 2px 4px); }
.skill-level em { font-style: normal; font-size: 11.5px; font-weight: 700; }
.skill-level em.ok { color: var(--success-strong); }
.skill-level em.gap { color: var(--warning-strong); }
.skill-level small { color: var(--text-3); font-size: 12px; }
</style>
