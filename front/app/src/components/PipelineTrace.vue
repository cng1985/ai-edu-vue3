<script setup>
import Icon from './Icon.vue'

defineProps({
  stages: { type: Array, default: () => [] },
  traces: { type: Array, default: () => [] },
  running: { type: Boolean, default: false }
})

function traceOf(traces, name) {
  return traces.find((t) => t.name === name)
}
</script>

<template>
  <ol class="trace">
    <li
      v-for="(stage, i) in stages"
      :key="stage.name"
      :class="[
        'trace__item',
        traceOf(traces, stage.name) ? `trace__item--${traceOf(traces, stage.name).status}` : '',
        { 'trace__item--active': running && !traceOf(traces, stage.name) && traces.length === i }
      ]"
    >
      <span class="trace__dot">
        <Icon v-if="traceOf(traces, stage.name)?.status === 'ok'" name="check" :size="12" :stroke="3" />
        <Icon v-else-if="traceOf(traces, stage.name)?.status === 'skipped'" name="arrowRight" :size="12" :stroke="3" />
        <Icon v-else-if="traceOf(traces, stage.name)?.status === 'failed'" name="x" :size="12" :stroke="3" />
        <b v-else>{{ i + 1 }}</b>
      </span>
      <div class="trace__body">
        <div class="trace__title">
          <strong>{{ stage.title }}</strong>
          <code>{{ stage.name }}</code>
          <small v-if="traceOf(traces, stage.name)" class="num">{{ traceOf(traces, stage.name).durationMs }}ms</small>
        </div>
        <p v-if="traceOf(traces, stage.name)">{{ traceOf(traces, stage.name).summary }}</p>
      </div>
    </li>
  </ol>
</template>

<style scoped>
.trace { list-style: none; margin: 0; padding: 0; }
.trace__item { position: relative; display: flex; gap: 12px; padding-bottom: 14px; }
.trace__item:not(:last-child)::before {
  content: '';
  position: absolute;
  left: 13px;
  top: 28px;
  bottom: 0;
  width: 2px;
  background: var(--border);
}
.trace__dot {
  display: grid;
  place-items: center;
  width: 28px;
  height: 28px;
  flex: 0 0 28px;
  border-radius: 50%;
  border: 2px solid var(--border-strong);
  background: var(--surface);
  color: var(--text-3);
  font-size: 11px;
}
.trace__item--ok .trace__dot { border-color: var(--success); background: var(--success); color: #fff; }
.trace__item--skipped .trace__dot { border-color: var(--warning); background: var(--warning-soft); color: var(--warning-strong); }
.trace__item--failed .trace__dot { border-color: var(--danger); background: var(--danger); color: #fff; }
.trace__item--active .trace__dot { border-color: var(--primary); color: var(--primary); animation: pulse 1s infinite; }
.trace__body { flex: 1; min-width: 0; padding-top: 3px; }
.trace__title { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.trace__title strong { font-size: 13.5px; }
.trace__title code { padding: 1px 6px; border-radius: 6px; background: var(--surface-3); color: var(--text-3); font-size: 11px; }
.trace__title small { margin-left: auto; color: var(--text-3); font-size: 11px; }
.trace__body p { margin: 3px 0 0; color: var(--text-2); font-size: 12.5px; line-height: 1.55; }
@keyframes pulse { 50% { box-shadow: 0 0 0 5px var(--primary-soft-2); } }
</style>
