<template>
  <el-container class="admin-layout">
    <AdminSidebar
      v-show="!isMobile"
      :collapsed="collapsed"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
    />

    <AdminMobileNav
      :open="mobileNavOpen"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
      @close="mobileNavOpen = false"
    />

    <el-container class="admin-layout__main">
      <el-header class="admin-header">
        <div class="admin-header__left">
          <el-button
            v-if="isMobile"
            :icon="Menu"
            text
            @click="mobileNavOpen = true"
          />
          <el-button
            v-else
            :icon="collapsed ? Expand : Fold"
            text
            @click="collapsed = !collapsed"
          />
          <AdminBreadcrumb :breadcrumbs="breadcrumbs" />
        </div>
        <div class="admin-header__right">
          <el-button
            class="admin-header__search"
            :icon="Search"
            text
            @click="commandOpen = true"
          >
            <span v-if="!isMobile" class="admin-header__search-text">搜索</span>
            <kbd v-if="!isMobile" class="admin-header__kbd">⌘K</kbd>
          </el-button>
          <el-tag size="small" :type="roleTagType">{{ roleLabel }}</el-tag>
          <el-dropdown @command="handleCommand">
            <span class="admin-header__user">
              <el-avatar :size="32" :style="{ background: auth.user?.avatarColor }">
                {{ auth.user?.avatar }}
              </el-avatar>
              <span v-if="!isMobile" class="admin-header__nickname">{{ auth.user?.nickname }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="dashboard">返回看板</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <AdminTopNav
        v-if="!isMobile"
        :visible-groups="visibleGroups"
        :active-group="activeGroup"
        @navigate="onNavigate"
      />

      <el-main class="admin-main" :class="{ 'admin-main--bleed': route.meta.fullBleed }">
        <router-view />
      </el-main>
    </el-container>

    <AdminCommandPalette
      v-model="commandOpen"
      :search-nav="searchNav"
      :visible-items="visibleItems"
      :recent-items="recentItems"
      @select="onNavigate"
    />
  </el-container>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Expand, Fold, Menu, Search } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'
import { useAdminNav, useCommandPaletteShortcut } from '../composables/useAdminNav'
import AdminSidebar from '../components/navigation/AdminSidebar.vue'
import AdminTopNav from '../components/navigation/AdminTopNav.vue'
import AdminBreadcrumb from '../components/navigation/AdminBreadcrumb.vue'
import AdminCommandPalette from '../components/navigation/AdminCommandPalette.vue'
import AdminMobileNav from '../components/navigation/AdminMobileNav.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const collapsed = ref(false)
const mobileNavOpen = ref(false)
const commandOpen = ref(false)
const isMobile = ref(false)

const {
  visibleGroups,
  visibleItems,
  activeMenu,
  activeGroup,
  breadcrumbs,
  recentItems,
  recordVisit,
  navigateTo,
  searchNav
} = useAdminNav()

useCommandPaletteShortcut(() => {
  commandOpen.value = true
})

watch(
  () => route.path,
  (path) => {
    recordVisit(path)
    mobileNavOpen.value = false
  },
  { immediate: true }
)

const roleMap = { admin: '管理员', reviewer: '审核员', operator: '运营' }
const roleLabel = computed(() => auth.user?.roleName || roleMap[auth.user?.role] || auth.user?.role)
const roleTagType = computed(() => {
  if (auth.user?.role === 'admin') return 'danger'
  if (auth.user?.role === 'reviewer') return 'warning'
  return 'info'
})

function onNavigate(path) {
  navigateTo(path)
}

function handleCommand(cmd) {
  if (cmd === 'logout') {
    auth.logout()
    router.push('/login')
  } else if (cmd === 'dashboard') {
    router.push('/dashboard')
  }
}

function updateMobile() {
  isMobile.value = window.innerWidth < 900
  if (!isMobile.value) mobileNavOpen.value = false
}

onMounted(() => {
  updateMobile()
  window.addEventListener('resize', updateMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateMobile)
})
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
}

.admin-layout__main {
  min-width: 0;
}

.admin-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 16px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}

.admin-header__left,
.admin-header__right {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.admin-header__right {
  flex-shrink: 0;
}

.admin-header__search {
  color: #6b7280;
}

.admin-header__search-text {
  margin-left: 4px;
  font-size: 13px;
}

.admin-header__kbd {
  margin-left: 8px;
  padding: 2px 6px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  background: #f9fafb;
  font-size: 11px;
  color: #9ca3af;
  font-family: inherit;
}

.admin-header__user {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.admin-header__nickname {
  font-size: 14px;
  color: #374151;
}

.admin-main {
  padding: 20px;
  background: #f0f2f5;
  min-height: calc(100vh - 96px);
}

.admin-main--bleed {
  padding: 0;
  background: #f8fafc;
}

@media (max-width: 900px) {
  .admin-main {
    min-height: calc(100vh - 56px);
  }
}
</style>
