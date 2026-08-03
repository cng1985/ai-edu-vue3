import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { NAV_GROUPS, NAV_CHILDREN, flattenNavItems } from '../config/nav'

const RECENT_KEY = 'admin-nav-recent'
const MAX_RECENT = 6

function loadRecent() {
  try {
    return JSON.parse(localStorage.getItem(RECENT_KEY) || '[]')
  } catch {
    return []
  }
}

function saveRecent(path) {
  const items = loadRecent().filter((p) => p !== path)
  items.unshift(path)
  localStorage.setItem(RECENT_KEY, JSON.stringify(items.slice(0, MAX_RECENT)))
}

export function useAdminNav() {
  const route = useRoute()
  const router = useRouter()
  const auth = useAuthStore()

  const visibleGroups = computed(() =>
    NAV_GROUPS.map((group) => ({
      ...group,
      items: group.items.filter(
        (item) => !item.permission || auth.hasPermission(item.permission)
      )
    })).filter((group) => group.items.length > 0)
  )

  const visibleItems = computed(() => flattenNavItems(visibleGroups.value))

  const activeMenu = computed(() => {
    const path = route.path
    const match = visibleItems.value.find(
      (item) => item.activePrefix ? path.startsWith(item.activePrefix) : path === item.path
    )
    return match?.path || path
  })

  const activeGroup = computed(() => {
    const current = visibleItems.value.find((item) => item.path === activeMenu.value)
    return current?.group || visibleGroups.value[0]?.key || ''
  })

  const breadcrumbs = computed(() => {
    const crumbs = [{ title: '管理后台', path: '/dashboard' }]
    const childMeta = NAV_CHILDREN[route.name]

    if (childMeta) {
      const parent = visibleItems.value.find((item) => item.name === childMeta.parent)
      if (parent) {
        crumbs.push({ title: parent.title, path: parent.path })
      }
      crumbs.push({ title: route.meta?.title || childMeta.title, path: route.path })
      return crumbs
    }

    const current = visibleItems.value.find((item) => item.path === activeMenu.value)
    if (current) {
      crumbs.push({ title: current.title, path: current.path })
    } else if (route.meta?.title) {
      crumbs.push({ title: route.meta.title, path: route.path })
    }
    return crumbs
  })

  const recentPaths = ref(loadRecent())

  function recordVisit(path) {
    if (!path || path === '/login') return
    saveRecent(path)
    recentPaths.value = loadRecent()
  }

  function getItemByPath(path) {
    return visibleItems.value.find((item) => item.path === path)
  }

  const recentItems = computed(() =>
    recentPaths.value
      .map((path) => getItemByPath(path))
      .filter(Boolean)
  )

  function navigateTo(path) {
    router.push(path)
  }

  function searchNav(query) {
    const q = query.trim().toLowerCase()
    if (!q) return visibleItems.value
    return visibleItems.value.filter((item) => {
      const haystack = [item.title, item.groupTitle, ...(item.keywords || [])].join(' ').toLowerCase()
      return haystack.includes(q)
    })
  }

  return {
    visibleGroups,
    visibleItems,
    activeMenu,
    activeGroup,
    breadcrumbs,
    recentItems,
    recordVisit,
    navigateTo,
    searchNav
  }
}

/** 全局快捷键：Ctrl/Cmd + K 打开命令面板 */
export function useCommandPaletteShortcut(onOpen) {
  function handleKeydown(e) {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      onOpen()
    }
  }

  onMounted(() => window.addEventListener('keydown', handleKeydown))
  onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
}
