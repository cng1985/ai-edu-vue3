<template>
  <router-view v-if="route.meta.chatgptFullscreen" />
  <el-container v-else class="admin-layout" :class="layoutClass">
    <!-- 经典风格：可折叠侧边栏 -->
    <AdminSidebar
      v-if="showClassicSidebar"
      :collapsed="collapsed"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
    />

    <!-- Win11 风格：图标导航栏 -->
    <AdminSidebarWin11
      v-if="showWin11Nav"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
      :active-group="activeGroup"
    />

    <AdminMobileNav
      :open="mobileNavOpen"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
      @close="mobileNavOpen = false"
    />

    <el-container class="admin-layout__main">
      <el-header class="admin-header" :class="{ 'admin-header--win11': navMode === 'win11' }">
        <div class="admin-header__left">
          <el-button
            v-if="isMobile"
            :icon="Menu"
            text
            @click="mobileNavOpen = true"
          />
          <el-button
            v-else-if="showCollapseToggle"
            :icon="collapsed ? Expand : Fold"
            text
            @click="toggleCollapsed"
          />
          <AdminBreadcrumb :breadcrumbs="breadcrumbs" />
        </div>
        <div class="admin-header__right">
          <AdminNavModeSwitcher
            v-if="!isMobile"
            :nav-mode="navMode"
            :compact="isMobile"
            @change="setNavMode"
          />
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

      <!-- 经典风格：顶部分组标签 -->
      <AdminTopNav
        v-if="showClassicTopNav"
        :visible-groups="visibleGroups"
        :active-group="activeGroup"
        @navigate="onNavigate"
      />

      <!-- 现代风格：水平下拉导航 -->
      <AdminNavModern
        v-if="showModernNav"
        :visible-groups="visibleGroups"
        :active-menu="activeMenu"
        :active-group="activeGroup"
      />

      <el-main class="admin-main" :class="{ 'admin-main--win11': navMode === 'win11' }">
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
import { useAdminNavLayout } from '../composables/useAdminNavLayout'
import AdminSidebar from '../components/navigation/AdminSidebar.vue'
import AdminSidebarWin11 from '../components/navigation/AdminSidebarWin11.vue'
import AdminTopNav from '../components/navigation/AdminTopNav.vue'
import AdminNavModern from '../components/navigation/AdminNavModern.vue'
import AdminNavModeSwitcher from '../components/navigation/AdminNavModeSwitcher.vue'
import AdminBreadcrumb from '../components/navigation/AdminBreadcrumb.vue'
import AdminCommandPalette from '../components/navigation/AdminCommandPalette.vue'
import AdminMobileNav from '../components/navigation/AdminMobileNav.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const { navMode, collapsed, setNavMode, toggleCollapsed } = useAdminNavLayout()

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

const showClassicSidebar = computed(() => !isMobile.value && navMode.value === 'classic')
const showWin11Nav = computed(() => !isMobile.value && navMode.value === 'win11')
const showClassicTopNav = computed(() => !isMobile.value && navMode.value === 'classic')
const showModernNav = computed(() => !isMobile.value && navMode.value === 'modern')
const showCollapseToggle = computed(() => navMode.value === 'classic')

const layoutClass = computed(() => ({
  'admin-layout--classic': navMode.value === 'classic',
  'admin-layout--win11': navMode.value === 'win11',
  'admin-layout--modern': navMode.value === 'modern'
}))

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

.admin-header--win11 {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  border-bottom-color: rgba(0, 0, 0, 0.06);
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

.admin-layout--modern .admin-main {
  min-height: calc(100vh - 100px);
}

.admin-layout--win11 .admin-main {
  background: #f5f5f5;
  min-height: calc(100vh - 56px);
}

.admin-main--win11 {
  min-height: calc(100vh - 56px);
}

@media (max-width: 900px) {
  .admin-main {
    min-height: calc(100vh - 56px);
  }
}
</style>
