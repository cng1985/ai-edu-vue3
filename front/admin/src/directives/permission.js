import { useAuthStore } from '../stores/auth'

function allowed(value) {
  const auth = useAuthStore()
  if (!value) return true
  if (Array.isArray(value)) return auth.hasAnyPermission(value)
  return auth.hasPermission(value)
}

/** 无权限时移除 DOM 元素，用于按钮级控制 */
export const permission = {
  mounted(el, binding) {
    if (!allowed(binding.value)) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el, binding) {
    if (!allowed(binding.value)) {
      el.style.display = 'none'
    }
  }
}
