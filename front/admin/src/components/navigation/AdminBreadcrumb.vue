<template>
  <el-breadcrumb separator="/" class="admin-breadcrumb">
    <el-breadcrumb-item
      v-for="(crumb, index) in breadcrumbs"
      :key="crumb.path"
    >
      <router-link
        v-if="index < breadcrumbs.length - 1"
        :to="crumb.path"
        class="admin-breadcrumb__link"
      >
        {{ crumb.title }}
      </router-link>
      <span v-else>{{ crumb.title }}</span>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
defineProps({
  breadcrumbs: { type: Array, required: true }
})
</script>

<style scoped>
.admin-breadcrumb__link {
  color: #6b7280;
  text-decoration: none;
}

.admin-breadcrumb__link:hover {
  color: #4f46e5;
}
</style>
