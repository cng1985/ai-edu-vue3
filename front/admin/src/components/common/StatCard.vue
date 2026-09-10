<template>
  <div class="stat-card" :class="[`stat-card--${tone}`, { 'stat-card--clickable': clickable }]" @click="$emit('click')">
    <div class="stat-card__icon">
      <el-icon :size="20"><component :is="icon" /></el-icon>
    </div>
    <div class="stat-card__body">
      <div class="stat-card__label">{{ label }}</div>
      <div class="stat-card__value num" :class="{ 'stat-card__value--text': isText }">
        <slot>{{ value }}</slot>
      </div>
      <div v-if="hint" class="stat-card__hint">{{ hint }}</div>
    </div>
    <div class="stat-card__glow" />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [Number, String], default: 0 },
  hint: { type: String, default: '' },
  icon: { type: [String, Object], default: 'DataAnalysis' },
  tone: { type: String, default: 'primary' },
  clickable: { type: Boolean, default: false }
})

defineEmits(['click'])

const isText = computed(() => typeof props.value === 'string' && Number.isNaN(Number(props.value)))
</script>

<style scoped>
.stat-card {
  --tone: var(--primary);
  --tone-soft: var(--primary-soft);
  position: relative;
  display: flex;
  gap: 14px;
  padding: 18px 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: transform var(--t), box-shadow var(--t), border-color var(--t);
}

.stat-card--clickable {
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow);
  border-color: color-mix(in srgb, var(--tone) 35%, var(--border));
}

.stat-card--success { --tone: var(--success); --tone-soft: var(--success-soft); }
.stat-card--warning { --tone: var(--warning); --tone-soft: var(--warning-soft); }
.stat-card--danger { --tone: var(--danger); --tone-soft: var(--danger-soft); }
.stat-card--info { --tone: var(--info); --tone-soft: var(--info-soft); }
.stat-card--sky { --tone: var(--sky); --tone-soft: var(--sky-soft); }
.stat-card--rose { --tone: var(--rose); --tone-soft: var(--rose-soft); }
.stat-card--ink { --tone: var(--ink-3); --tone-soft: var(--surface-3); }

.stat-card__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  flex-shrink: 0;
  border-radius: 13px;
  background: var(--tone-soft);
  color: var(--tone);
}

.stat-card__body {
  min-width: 0;
  flex: 1;
}

.stat-card__label {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-3);
  letter-spacing: 0.02em;
}

.stat-card__value {
  margin-top: 2px;
  font-size: 28px;
  font-weight: 800;
  line-height: 1.15;
  color: var(--text);
  letter-spacing: -0.02em;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stat-card__value--text {
  font-size: 16px;
  font-weight: 700;
  margin-top: 6px;
}

.stat-card__hint {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-3);
}

.stat-card__glow {
  position: absolute;
  right: -30px;
  top: -30px;
  width: 110px;
  height: 110px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--tone-soft) 0%, transparent 70%);
  opacity: 0.9;
  pointer-events: none;
}
</style>
