<template>
  <nav class="admin-breadcrumb" aria-label="面包屑">
    <template v-for="(crumb, index) in breadcrumbs" :key="crumb.path">
      <el-icon v-if="index > 0" :size="12" class="admin-breadcrumb__sep"><ArrowRight /></el-icon>
      <router-link
        v-if="index < breadcrumbs.length - 1"
        :to="crumb.path"
        class="admin-breadcrumb__link"
      >
        {{ crumb.title }}
      </router-link>
      <span v-else class="admin-breadcrumb__current">{{ crumb.title }}</span>
    </template>
  </nav>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'

defineProps({
  breadcrumbs: { type: Array, required: true }
})
</script>

<style scoped>
.admin-breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13.5px;
  white-space: nowrap;
  overflow: hidden;
}

.admin-breadcrumb__sep {
  color: var(--border-strong);
  flex-shrink: 0;
}

.admin-breadcrumb__link {
  color: var(--text-3);
  text-decoration: none;
  font-weight: 500;
  padding: 3px 6px;
  border-radius: 7px;
  transition: all var(--t-fast);
}

.admin-breadcrumb__link:hover {
  color: var(--primary);
  background: var(--primary-soft);
}

.admin-breadcrumb__current {
  color: var(--text);
  font-weight: 700;
  padding: 3px 6px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
