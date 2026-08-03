<template>
  <el-container class="admin-layout">
    <el-aside :width="collapsed ? '64px' : '220px'" class="sidebar">
      <div class="sidebar__logo">
        <span v-if="!collapsed">AI 学习管理</span>
        <span v-else>AI</span>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="collapsed"
        router
        background-color="#1e1b4b"
        text-color="#c7d2fe"
        active-text-color="#ffffff"
      >
        <el-menu-item v-if="auth.hasPermission(PERM.DASHBOARD)" index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据看板</span>
        </el-menu-item>
        <el-menu-item v-if="auth.hasPermission(PERM.USER_READ)" index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item v-if="auth.hasPermission(PERM.COURSE_READ)" index="/courses">
          <el-icon><Reading /></el-icon>
          <span>课程管理</span>
        </el-menu-item>
        <el-menu-item v-if="auth.hasPermission(PERM.QUIZ_READ)" index="/quizzes">
          <el-icon><EditPen /></el-icon>
          <span>题库管理</span>
        </el-menu-item>
        <el-menu-item v-if="auth.hasPermission(PERM.CUSTOMER_READ)" index="/customers">
          <el-icon><Service /></el-icon>
          <span>客户咨询</span>
        </el-menu-item>
        <el-menu-item v-if="auth.hasPermission(PERM.ROLE_MANAGE)" index="/roles">
          <el-icon><Lock /></el-icon>
          <span>权限管理</span>
        </el-menu-item>
        <el-menu-item v-if="auth.isAdmin" index="/settings">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
        <el-menu-item v-if="auth.isReviewer" index="/reviews">
          <el-icon><DocumentChecked /></el-icon>
          <span>内容审核</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header__left">
          <el-button :icon="collapsed ? Expand : Fold" text @click="collapsed = !collapsed" />
          <el-breadcrumb separator="/">
            <el-breadcrumb-item>管理后台</el-breadcrumb-item>
            <el-breadcrumb-item>{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header__right">
          <el-tag size="small" :type="roleTagType">{{ roleLabel }}</el-tag>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :style="{ background: auth.user?.avatarColor }">
                {{ auth.user?.avatar }}
              </el-avatar>
              <span class="nickname">{{ auth.user?.nickname }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const collapsed = ref(false)

const roleMap = { admin: '管理员', reviewer: '审核员', operator: '运营' }
const roleLabel = computed(() => auth.user?.roleName || roleMap[auth.user?.role] || auth.user?.role)
const roleTagType = computed(() => {
  if (auth.user?.role === 'admin') return 'danger'
  if (auth.user?.role === 'reviewer') return 'warning'
  return 'info'
})

function handleCommand(cmd) {
  if (cmd === 'logout') {
    auth.logout()
    router.push('/login')
  }
}
</script>
