<script setup>
defineProps({
  user: { type: Object, default: null },
  size: { type: String, default: '' }
})
</script>

<template>
  <span class="avatar" :class="size && `avatar--${size}`" :style="{ background: user?.avatarColor || '#6b5cff' }">
    {{ user?.avatar || user?.nickname?.slice(0, 1) || '?' }}
  </span>
</template>
