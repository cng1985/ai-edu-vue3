<template>
  <div>
    <div class="page-header">
      <h2>数据看板</h2>
      <el-button :icon="Refresh" @click="loadStats">刷新</el-button>
    </div>

    <el-row :gutter="16" v-loading="loading">
      <el-col :span="6" v-for="item in userCards" :key="item.label">
        <el-card shadow="hover" class="stat-card">
          <div class="value">{{ item.value }}</div>
          <div class="label">{{ item.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6" v-for="item in courseCards" :key="item.label">
        <el-card shadow="hover" class="stat-card">
          <div class="value" :style="{ color: item.color }">{{ item.value }}</div>
          <div class="label">{{ item.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><span>内容概览</span></template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="测验套数">{{ stats?.quizStats?.total || 0 }}</el-descriptions-item>
            <el-descriptions-item label="题目总数">{{ stats?.quizStats?.questions || 0 }}</el-descriptions-item>
            <el-descriptions-item label="待审核">{{ stats?.reviewStats?.pending || 0 }}</el-descriptions-item>
            <el-descriptions-item label="已通过">{{ stats?.reviewStats?.approved || 0 }}</el-descriptions-item>
            <el-descriptions-item label="已驳回">{{ stats?.reviewStats?.rejected || 0 }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header><span>快捷操作</span></template>
          <div style="display: flex; flex-wrap: wrap; gap: 12px">
            <el-button v-if="auth.hasPermission(PERM.COURSE_READ)" type="primary" @click="$router.push('/courses')">管理课程</el-button>
            <el-button v-if="auth.hasPermission(PERM.QUIZ_READ)" @click="$router.push('/quizzes')">管理题库</el-button>
            <el-button v-if="auth.hasPermission(PERM.REVIEW_READ)" type="warning" @click="$router.push('/reviews')">审核队列</el-button>
            <el-button v-if="auth.hasPermission(PERM.USER_READ)" @click="$router.push('/users')">用户管理</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { dashboardApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

const auth = useAuthStore()
const loading = ref(false)
const stats = ref(null)

const userCards = computed(() => [
  { label: '用户总数', value: stats.value?.userStats?.total || 0 },
  { label: '学员数', value: stats.value?.userStats?.learners || 0 },
  { label: '管理/运营', value: stats.value?.userStats?.admins || 0 },
  { label: '活跃账号', value: stats.value?.userStats?.active || 0 }
])

const courseCards = computed(() => [
  { label: '课程总数', value: stats.value?.courseStats?.total || 0, color: '#6366f1' },
  { label: '已发布', value: stats.value?.courseStats?.published || 0, color: '#10b981' },
  { label: '草稿', value: stats.value?.courseStats?.draft || 0, color: '#f59e0b' },
  { label: '章节总数', value: stats.value?.courseStats?.chapters || 0, color: '#0ea5e9' }
])

async function loadStats() {
  loading.value = true
  try {
    stats.value = await dashboardApi.stats()
  } finally {
    loading.value = false
  }
}

onMounted(loadStats)
</script>
