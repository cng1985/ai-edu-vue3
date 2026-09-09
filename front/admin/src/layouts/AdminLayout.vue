<template>
  <router-view v-if="route.meta.chatgptFullscreen" />
  <div v-else class="admin-layout" :class="layoutClass">
    <div class="admin-layout__glow" aria-hidden="true" />

    <!-- 经典风格：可折叠侧边栏 -->
    <AdminSidebar
      v-if="showClassicSidebar"
      :collapsed="collapsed"
      :visible-groups="visibleGroups"
      :active-menu="activeMenu"
      @toggle="toggleCollapsed"
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

    <div class="admin-layout__main">
      <header class="admin-header">
        <div class="admin-header__left">
          <button
            v-if="isMobile"
            type="button"
            class="admin-header__icon-btn"
            title="打开菜单"
            @click="mobileNavOpen = true"
          >
            <el-icon :size="18"><Menu /></el-icon>
          </button>
          <button
            v-else-if="showCollapseToggle"
            type="button"
            class="admin-header__icon-btn"
            :title="collapsed ? '展开侧栏' : '收起侧栏'"
            @click="toggleCollapsed"
          >
            <el-icon :size="18"><component :is="collapsed ? Expand : Fold" /></el-icon>
          </button>
          <div v-if="showModernNav" class="admin-header__brand">
            <span class="admin-header__logo">AI</span>
            <span class="admin-header__brand-text">学习管理</span>
          </div>
          <AdminBreadcrumb v-else :breadcrumbs="breadcrumbs" />
        </div>

        <div class="admin-header__right">
          <button type="button" class="admin-header__search" @click="commandOpen = true">
            <el-icon :size="16"><Search /></el-icon>
            <span v-if="!isMobile" class="admin-header__search-text">搜索页面…</span>
            <kbd v-if="!isMobile" class="admin-header__kbd">⌘K</kbd>
          </button>

          <AdminNavModeSwitcher
            v-if="!isMobile"
            :nav-mode="navMode"
            :compact="isMobile"
            @change="setNavMode"
          />

          <el-dropdown trigger="click" @command="handleCommand">
            <button type="button" class="admin-header__user">
              <span class="admin-header__avatar" :style="{ background: auth.user?.avatarColor || 'var(--primary)' }">
                {{ auth.user?.avatar || auth.user?.nickname?.slice(0, 1) || '管' }}
              </span>
              <span v-if="!isMobile" class="admin-header__user-meta">
                <span class="admin-header__nickname">{{ auth.user?.nickname }}</span>
                <span class="admin-header__role" :class="`admin-header__role--${auth.user?.role}`">{{ roleLabel }}</span>
              </span>
              <el-icon v-if="!isMobile" :size="14" class="admin-header__chevron"><ArrowDown /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="dashboard">
                  <el-icon><DataAnalysis /></el-icon>返回看板
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

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

      <main class="admin-main">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <AdminCommandPalette
      v-model="commandOpen"
      :search-nav="searchNav"
      :visible-items="visibleItems"
      :recent-items="recentItems"
      @select="onNavigate"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Expand, Fold, Menu, Search, ArrowDown, DataAnalysis, SwitchButton } from '@element-plus/icons-vue'
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
  position: relative;
  display: flex;
  min-height: 100vh;
  background: var(--bg);
}

.admin-layout__glow {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background:
    radial-gradient(800px 400px at 85% -10%, rgba(107, 92, 255, 0.10), transparent 60%),
    radial-gradient(600px 300px at 10% 110%, rgba(15, 185, 129, 0.08), transparent 60%);
}

.admin-layout__main {
  position: relative;
  z-index: 1;
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.admin-header {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  height: var(--header-height);
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.78);
  backdrop-filter: blur(18px) saturate(1.4);
  -webkit-backdrop-filter: blur(18px) saturate(1.4);
  border-bottom: 1px solid var(--border);
}

.admin-header__left,
.admin-header__right {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.admin-header__right {
  flex-shrink: 0;
}

.admin-header__icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: 1px solid transparent;
  border-radius: 11px;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-header__icon-btn:hover {
  background: var(--surface-2);
  border-color: var(--border);
  color: var(--text);
}

.admin-header__brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-header__logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--primary) 0%, #9a8cff 100%);
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  font-size: 13px;
  box-shadow: 0 6px 14px var(--primary-glow);
}

.admin-header__brand-text {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-weight: 700;
  font-size: 15px;
  color: var(--text);
}

.admin-header__search {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 38px;
  padding: 0 12px;
  border: 1px solid var(--border);
  border-radius: 11px;
  background: var(--surface-2);
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-header__search:hover {
  border-color: var(--border-strong);
  background: var(--surface);
  color: var(--text-2);
}

.admin-header__search-text {
  font-size: 13px;
  min-width: 90px;
  text-align: left;
}

.admin-header__kbd {
  padding: 1px 6px;
  border: 1px solid var(--border-strong);
  border-radius: 6px;
  background: var(--surface);
  font-size: 11px;
  color: var(--text-3);
  font-family: inherit;
}

.admin-header__user {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 44px;
  padding: 0 6px 0 4px;
  border: 1px solid transparent;
  border-radius: 13px;
  background: transparent;
  cursor: pointer;
  transition: all var(--t-fast);
}

.admin-header__user:hover {
  background: var(--surface-2);
  border-color: var(--border);
}

.admin-header__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 11px;
  color: #fff;
  font-weight: 700;
  font-size: 14px;
}

.admin-header__user-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  line-height: 1.2;
}

.admin-header__nickname {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text);
}

.admin-header__role {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-3);
}

.admin-header__role--admin { color: var(--danger); }
.admin-header__role--reviewer { color: var(--warning-strong); }
.admin-header__role--operator { color: var(--info); }

.admin-header__chevron {
  color: var(--text-3);
}

.admin-main {
  flex: 1;
  padding: 26px 28px 40px;
}

@media (max-width: 900px) {
  .admin-header {
    padding: 0 14px;
    height: 58px;
  }

  .admin-main {
    padding: 18px 14px 32px;
  }
}
</style>
