import { ref, watch } from 'vue'
import { DEFAULT_COLLAPSED, DEFAULT_NAV_MODE } from '../config/navLayout'

const STORAGE_KEY = 'admin-nav-layout'

function loadLayoutPrefs() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return { navMode: DEFAULT_NAV_MODE, collapsed: DEFAULT_COLLAPSED }
    }
    const parsed = JSON.parse(raw)
    return {
      navMode: parsed.navMode || DEFAULT_NAV_MODE,
      collapsed: parsed.collapsed ?? DEFAULT_COLLAPSED
    }
  } catch {
    return { navMode: DEFAULT_NAV_MODE, collapsed: DEFAULT_COLLAPSED }
  }
}

function saveLayoutPrefs(navMode, collapsed) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({ navMode, collapsed }))
}

const prefs = loadLayoutPrefs()
const navMode = ref(prefs.navMode)
const collapsed = ref(prefs.collapsed)

watch([navMode, collapsed], () => {
  saveLayoutPrefs(navMode.value, collapsed.value)
})

export function useAdminNavLayout() {
  function setNavMode(mode) {
    navMode.value = mode
  }

  function toggleCollapsed() {
    collapsed.value = !collapsed.value
  }

  return {
    navMode,
    collapsed,
    setNavMode,
    toggleCollapsed
  }
}
