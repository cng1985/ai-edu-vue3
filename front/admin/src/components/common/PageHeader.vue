<template>
  <header class="page-header">
    <div class="page-header__main">
      <button v-if="back" type="button" class="page-header__back" title="返回" @click="$emit('back')">
        <el-icon :size="18"><ArrowLeft /></el-icon>
      </button>
      <div class="page-header__text">
        <div v-if="eyebrow" class="page-header__eyebrow">{{ eyebrow }}</div>
        <div class="page-header__title-row">
          <h2>{{ title }}</h2>
          <slot name="title-extra" />
        </div>
        <p v-if="subtitle" class="page-header__sub">{{ subtitle }}</p>
      </div>
    </div>
    <div v-if="$slots.default" class="page-header__actions">
      <slot />
    </div>
  </header>
</template>

<script setup>
import { ArrowLeft } from '@element-plus/icons-vue'

defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  eyebrow: { type: String, default: '' },
  back: { type: Boolean, default: false }
})

defineEmits(['back'])
</script>

<style scoped>
.page-header__text {
  min-width: 0;
}

.page-header__title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
